package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

import org.springframework.http.codec.multipart.FilePart;
import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.client.MinioClientService;
import com.knowledgebase.platformspring.client.OpenAIClientService;
import com.knowledgebase.platformspring.client.PaddleOCRClientService;
import com.knowledgebase.platformspring.client.QdrantClientService;
import com.knowledgebase.platformspring.dto.ChatDocumentRequest;
import com.knowledgebase.platformspring.dto.ChatDocumentResponse;
import com.knowledgebase.platformspring.dto.HomepageResponse;
import com.knowledgebase.platformspring.dto.KnowledgeSearchRequest;
import com.knowledgebase.platformspring.dto.KnowledgeSearchResponse;
import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.dto.TagCloudResponse;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.Document;
import com.knowledgebase.platformspring.model.DocumentChunk;
import com.knowledgebase.platformspring.model.Workflow;
import com.knowledgebase.platformspring.repository.DocumentChunkRepository;
import com.knowledgebase.platformspring.repository.DocumentRepository;
import com.knowledgebase.platformspring.repository.SpaceRepository;
import com.knowledgebase.platformspring.repository.SubSpaceRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Schedulers;

/**
 * 文档服务 - 整合所有文档相关功能
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class DocumentService {
    
    private final DocumentRepository documentRepository;
    private final DocumentChunkRepository documentChunkRepository;
    private final SpaceRepository spaceRepository;
    private final SubSpaceRepository subSpaceRepository;
    private final MinioClientService minioClient;
    private final OpenAIClientService openAIClient;
    @SuppressWarnings("unused") // Reserved for future OCR implementation
    private final PaddleOCRClientService ocrClient;
    private final QdrantClientService qdrantClient;
    private final WorkflowService workflowService;
    
    public Mono<Document> uploadDocument(FilePart filePart, Long spaceId, Long subSpaceId, 
                                        Long classId, Long userId, String nickName,
                                        String fileName, String tags, String summary,
                                        String department, Boolean needApproval, 
                                        String version, String useType) {
        
        String objectName = "documents/" + UUID.randomUUID().toString() + "_" + fileName;
        String fileExt = fileName.contains(".") ? 
                fileName.substring(fileName.lastIndexOf(".")) : "";
        String title = fileName.contains(".") ?
                fileName.substring(0, fileName.lastIndexOf(".")) : fileName;
        
        Long fileSize = filePart.headers().getContentLength();
        
        Document document = Document.builder()
                .title(title)
                .fileName(fileName)
                .filePath(objectName)
                .fileType(fileExt)
                .fileSize(fileSize)
                .spaceId(spaceId)
                .subSpaceId(subSpaceId)
                .classId(classId)
                .createdBy(userId)
                .creatorNickName(nickName)
                .department(department)
                .tags(tags)
                .summary(summary)
                .needApproval(needApproval != null ? needApproval : false)
                .version(version)
                .useType(useType)
                .status(Document.STATUS_UPLOADING)
                .createdAt(LocalDateTime.now())
                .updatedAt(LocalDateTime.now())
                .build();
        
        return documentRepository.save(document)
                .flatMap(savedDoc -> 
                    filePart.content()
                            .reduce((buffer1, buffer2) -> buffer1.write(buffer2))
                            .flatMap(dataBuffer -> {
                                byte[] bytes = new byte[dataBuffer.readableByteCount()];
                                dataBuffer.read(bytes);
                                
                                return minioClient.uploadFile(
                                    objectName, 
                                    new java.io.ByteArrayInputStream(bytes),
                                    bytes.length,
                                    filePart.headers().getContentType().toString()
                                ).thenReturn(savedDoc);
                            })
                )
                .doOnSuccess(doc -> {
                    // 异步处理文档（OCR/向量化）
                    processDocumentAsync(doc.getId())
                            .subscribeOn(Schedulers.boundedElastic())
                            .subscribe(
                                    result -> log.info("Document {} processed successfully", doc.getId()),
                                    error -> log.error("Failed to process document {}: {}", doc.getId(), error.getMessage())
                            );
                });
    }
    
    /**
     * 异步处理文档，不阻塞上传请求
     */
    private Mono<Document> processDocumentAsync(Long documentId) {
        return processDocument(documentId);
    }
    
    public Mono<Document> processDocument(Long documentId) {
        log.info("ProcessDocument started for document ID: {}", documentId);
        
        return documentRepository.findById(documentId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")))
                // 1. 更新状态为处理中 (10%)
                .flatMap(document -> {
                    return updateDocumentStatus(documentId, Document.STATUS_PROCESSING, 10, "开始处理文档...")
                            .thenReturn(document);
                })
                // 2. 从 MinIO 下载并提取文本 (20-60%)
                .flatMap(document -> downloadAndExtractText(document)
                        .flatMap(text -> {
                            // ⭐ 保存 content 到数据库
                            document.setContent(text);
                            return documentRepository.save(document);
                        })
                        .flatMap(savedDoc -> updateDocumentStatus(documentId, Document.STATUS_PROCESSING, 60, "文本提取完成，开始分块...")
                                .thenReturn(savedDoc)))
                // 3. 分块 (70%)
                .flatMap(document -> updateDocumentStatus(documentId, Document.STATUS_VECTORIZING, 70, "开始向量化处理...")
                        .thenReturn(document))
                // 4. 向量化并存储chunks（创建 DocumentChunk 记录和 Qdrant points）
                .flatMap(this::chunkAndVectorize)
                // 5. 更新 processed_at 和清空 parse_error
                .flatMap(document -> {
                    document.setProcessedAt(LocalDateTime.now());
                    document.setParseError(""); // 清空错误信息
                    return documentRepository.save(document);
                })
                // 6. 根据 needApproval 决定最终状态
                .flatMap(document -> {
                    if (Boolean.TRUE.equals(document.getNeedApproval())) {
                        // 需要审批：创建并启动审批流程
                        return createAndStartWorkflow(document)
                                .flatMap(workflow -> {
                                    document.setWorkflowId(workflow.getId());
                                    return updateDocumentStatus(documentId, Document.STATUS_PENDING_APPROVAL, 100, "处理完成，等待审批");
                                })
                                .onErrorResume(e -> {
                                    log.error("Failed to create workflow for document {}: {}", documentId, e.getMessage());
                                    // 审批流程创建失败，直接设为待发布
                                    return updateDocumentStatus(documentId, Document.STATUS_PENDING_PUBLISH, 100, "处理完成，但审批流程创建失败，已设为待发布");
                                });
                    } else {
                        // 不需要审批：直接设为待发布
                        return updateDocumentStatus(documentId, Document.STATUS_PENDING_PUBLISH, 100, "处理完成，待发布");
                    }
                })
                .onErrorResume(e -> {
                    log.error("Failed to process document {}: {}", documentId, e.getMessage(), e);
                    return markDocumentProcessingError(documentId, e);
                });
    }
    
    /**
     * 从 MinIO 下载文件并提取文本
     */
    private Mono<String> downloadAndExtractText(Document document) {
        return updateDocumentStatus(document.getId(), Document.STATUS_PROCESSING, 20, "下载文件完成，开始读取...")
                .flatMap(doc -> minioClient.downloadFile(document.getFilePath()))
                .flatMap(inputStream -> {
                    try {
                        byte[] fileBytes = inputStream.readAllBytes();
                        inputStream.close();
                        
                        return updateDocumentStatus(document.getId(), Document.STATUS_PROCESSING, 30, "开始文本提取...")
                                .flatMap(d -> extractPlainText(document.getFileType(), fileBytes, document));
                    } catch (Exception e) {
                        return Mono.error(new BusinessException("Failed to read file: " + e.getMessage()));
                    }
                });
    }
    
    /**
     * 提取纯文本（支持多种格式）
     */
    private Mono<String> extractPlainText(String fileType, byte[] data, Document document) {
        String text;
        
        switch (fileType.toLowerCase()) {
            case ".txt":
            case ".md":
            case ".csv":
            case ".log":
                text = new String(data);
                break;
            case ".json":
                text = new String(data);
                break;
            default:
                // 不支持的格式，尝试 OCR（如果配置了）
                log.warn("Unsupported file type: {}, OCR not implemented yet", fileType);
                return Mono.error(new BusinessException("Unsupported file type: " + fileType));
        }
        
        if (text == null || text.trim().isEmpty()) {
            return Mono.error(new BusinessException("Empty text extracted from document"));
        }
        
        return Mono.just(text);
    }
    
    /**
     * 更新文档状态和进度
     */
    private Mono<Document> updateDocumentStatus(Long documentId, String status, Integer progress, String message) {
        return documentRepository.findById(documentId)
                .flatMap(doc -> {
                    doc.setStatus(status);
                    doc.setProcessProgress(progress);
                    log.info("Document {} status updated: {} ({}%) - {}", documentId, status, progress, message);
                    return documentRepository.save(doc);
                });
    }
    
    /**
     * 标记文档处理错误
     */
    private Mono<Document> markDocumentProcessingError(Long documentId, Throwable error) {
        return documentRepository.findById(documentId)
                .flatMap(doc -> {
                    doc.setStatus(Document.STATUS_PROCESS_FAILED);
                    doc.setParseError(error.getMessage());
                    doc.setProcessProgress(0);
                    doc.setVectorCount(0);
                    log.error("Document {} marked as process_failed: {}", documentId, error.getMessage());
                    return documentRepository.save(doc);
                });
    }
    
    /**
     * 分块并向量化（先删除旧 chunks，然后批量处理）
     */
    private Mono<Document> chunkAndVectorize(Document document) {
        if (document.getContent() == null || document.getContent().isEmpty()) {
            return Mono.just(document);
        }
        
        // 使用固定大小分块策略 (800字符一块，120字符重叠)
        List<String> chunks = splitIntoChunks(document.getContent(), 800, 120);
        
        if (chunks.isEmpty()) {
            return Mono.error(new BusinessException("No valid chunks to store"));
        }
        
        log.info("Document {}: preparing to generate embeddings for {} chunks", document.getId(), chunks.size());
        
        // 1. 先删除旧的 chunks
        return documentChunkRepository.deleteByDocumentId(document.getId())
                .then(Mono.defer(() -> {
                    // 2. 批量生成 embeddings 并存储
                    return Flux.fromIterable(chunks)
                            .index()
                            .flatMap(tuple -> {
                                String chunkContent = tuple.getT2();
                                int index = tuple.getT1().intValue();
                                
                                return openAIClient.createEmbedding(chunkContent)
                                        .flatMap(embedding -> {
                                            // 创建 Qdrant point
                                            Map<String, Object> payload = new HashMap<>();
                                            payload.put("document_id", document.getId());
                                            payload.put("space_id", document.getSpaceId());
                                            payload.put("sub_space_id", document.getSubSpaceId());
                                            payload.put("class_id", document.getClassId());
                                            payload.put("index", index);
                                            payload.put("content", chunkContent);
                                            
                                            QdrantClientService.QdrantPoint point = 
                                                QdrantClientService.QdrantPoint.create(embedding, payload);
                                            
                                            // 保存到 Qdrant
                                            return qdrantClient.upsertPoints(List.of(point))
                                                    .then(Mono.fromCallable(() -> {
                                                        // 保存到数据库
                                                        DocumentChunk chunk = DocumentChunk.builder()
                                                                .documentId(document.getId())
                                                                .index(index)
                                                                .content(chunkContent)
                                                                .vectorId(point.getId())
                                                                .tokenCount(countTokens(chunkContent))
                                                                .createdAt(LocalDateTime.now())
                                                                .updatedAt(LocalDateTime.now())
                                                                .build();
                                                        return chunk;
                                                    }))
                                                    .flatMap(documentChunkRepository::save);
                                        });
                            })
                            .collectList()
                            .map(savedChunks -> {
                                document.setVectorCount(savedChunks.size());
                                log.info("Document {}: successfully stored {} chunks", document.getId(), savedChunks.size());
                                return document;
                            });
                }));
    }
    
    /**
     * 分块策略：固定大小 + 重叠
     */
    private List<String> splitIntoChunks(String text, int chunkSize, int overlap) {
        List<String> chunks = new ArrayList<>();
        int start = 0;
        
        while (start < text.length()) {
            int end = Math.min(start + chunkSize, text.length());
            String chunk = text.substring(start, end).trim();
            
            if (!chunk.isEmpty()) {
                chunks.add(chunk);
            }
            
            start += chunkSize - overlap;
        }
        
        return chunks;
    }
    
    /**
     * 简单的 token 计数（估算）
     */
    private int countTokens(String text) {
        // 简单估算：平均 4 个字符 = 1 个 token
        return text.length() / 4;
    }
    
    /**
     * 创建并启动审批流程（对齐 Go 版本）
     */
    private Mono<Workflow> createAndStartWorkflow(Document document) {
        // 1. 创建 Workflow（包含 Step）
        return workflowService.createWorkflowWithStep(
                document.getSpaceId(),
                document.getId(),
                "document",
                document.getCreatedBy()
        )
        .flatMap(workflow -> {
            // 2. 更新文档的 workflow_id
            document.setWorkflowId(workflow.getId());
            return documentRepository.save(document)
                    .thenReturn(workflow);
        })
        .flatMap(workflow -> {
            // 3. 启动审批流程（创建 Task）
            return workflowService.startWorkflow(
                    workflow.getId(),
                    document.getSpaceId(),
                    document.getCreatedBy()
            );
        })
        .doOnSuccess(workflow -> 
            log.info("Workflow {} created and started for document {}", workflow.getId(), document.getId())
        )
        .doOnError(e -> 
            log.error("Failed to create/start workflow for document {}: {}", document.getId(), e.getMessage())
        );
    }
    
    public Flux<Document> getDocumentsBySpaceId(Long spaceId) {
        return documentRepository.findBySpaceId(spaceId);
    }
    
    public Mono<PaginationResponse<List<Document>>> getDocumentsBySpaceIdPaginated(
            Long spaceId, Integer page, Integer pageSize) {
        return documentRepository.findBySpaceId(spaceId)
                .collectList()
                .flatMap(allDocs -> {
                    long total = allDocs.size();
                    int offset = (page - 1) * pageSize;
                    List<Document> pagedDocs = allDocs.stream()
                            .skip(offset)
                            .limit(pageSize)
                            .collect(Collectors.toList());
                    
                    return Mono.just(PaginationResponse.of(pagedDocs, total, page, pageSize));
                });
    }
    
    
    public Mono<Document> retryProcessDocument(Long documentId, boolean forceRetry) {
        return documentRepository.findById(documentId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("文档不存在")))
                .flatMap(document -> {
                    if (!forceRetry && 
                        !Document.STATUS_PROCESS_FAILED.equals(document.getStatus()) &&
                        !Document.STATUS_FAILED.equals(document.getStatus())) {
                        return Mono.just(document);
                    }
                    
                    // 更新重试计数
                    document.setRetryCount(document.getRetryCount() + 1);
                    document.setLastRetryAt(LocalDateTime.now());
                    document.setStatus(Document.STATUS_PROCESSING);
                    document.setProcessProgress(0);
                    document.setParseError(null);
                    
                    return documentRepository.save(document)
                            .flatMap(doc -> processDocument(doc.getId()));
                });
    }
    
    public Mono<Document> getDocumentById(Long id) {
        return documentRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")));
    }
    
    // ============ 从 DocumentEnhancedService 整合过来的方法 ============
    
    /**
     * 获取首页展示数据 - 返回5个知识库，每个3个二级库，每个二级库6个文档
     */
    public Mono<HomepageResponse> getHomepageDocuments() {
        return spaceRepository.findAll()
                .filter(space -> space.getStatus() == 1)
                .take(5)
                .flatMap(space -> {
                    HomepageResponse.HomepageSpace homepageSpace = HomepageResponse.HomepageSpace.builder()
                            .id(space.getId())
                            .name(space.getName())
                            .description(space.getDescription())
                            .subSpaces(new ArrayList<>())
                            .build();
                    
                    return subSpaceRepository.findBySpaceId(space.getId())
                            .filter(sub -> sub.getStatus() == 1)
                            .take(3)
                            .flatMap(subSpace -> {
                                HomepageResponse.HomepageSubSpace homepageSubSpace = 
                                    HomepageResponse.HomepageSubSpace.builder()
                                        .id(subSpace.getId())
                                        .name(subSpace.getName())
                                        .description(subSpace.getDescription())
                                        .documents(new ArrayList<>())
                                        .build();
                                
                                return documentRepository.findBySubSpaceId(subSpace.getId())
                                        .filter(doc -> Document.STATUS_PUBLISHED.equals(doc.getStatus()))
                                        .take(6)
                                        .map(doc -> HomepageResponse.HomepageDocument.builder()
                                                .id(doc.getId())
                                                .title(doc.getTitle())
                                                .fileName(doc.getFileName())
                                                .fileSize(doc.getFileSize())
                                                .fileType(doc.getFileType())
                                                .status(doc.getStatus())
                                                .creatorNickName(doc.getCreatorNickName())
                                                .summary(doc.getSummary())
                                                .createdAt(doc.getCreatedAt())
                                                .updatedAt(doc.getUpdatedAt())
                                                .build())
                                        .collectList()
                                        .doOnNext(homepageSubSpace::setDocuments)
                                        .thenReturn(homepageSubSpace);
                            })
                            .collectList()
                            .doOnNext(homepageSpace::setSubSpaces)
                            .thenReturn(homepageSpace);
                })
                .collectList()
                .map(spaces -> HomepageResponse.builder().spaces(spaces).build());
    }
    
    /**
     * 获取标签云
     */
    public Mono<TagCloudResponse> getTagCloud(Long spaceId, Long subSpaceId, Integer limit) {
        Flux<Document> query = documentRepository.findAll()
                .filter(doc -> Document.STATUS_PUBLISHED.equals(doc.getStatus()))
                .filter(doc -> doc.getTags() != null && !doc.getTags().isEmpty());
        
        if (spaceId != null) {
            query = query.filter(doc -> spaceId.equals(doc.getSpaceId()));
        }
        if (subSpaceId != null) {
            query = query.filter(doc -> subSpaceId.equals(doc.getSubSpaceId()));
        }
        
        return query
                .map(Document::getTags)
                .flatMap(tags -> {
                    // 简单按逗号分割tags
                    String[] tagArray = tags.split(",");
                    return Flux.fromArray(tagArray);
                })
                .map(String::trim)
                .filter(tag -> !tag.isEmpty())
                .collect(Collectors.groupingBy(tag -> tag, Collectors.counting()))
                .map(tagCounts -> {
                    List<TagCloudResponse.TagCloudItem> items = tagCounts.entrySet().stream()
                            .map(entry -> TagCloudResponse.TagCloudItem.builder()
                                    .tag(entry.getKey())
                                    .count(entry.getValue().intValue())
                                    .build())
                            .sorted((a, b) -> Integer.compare(b.getCount(), a.getCount()))
                            .limit(limit)
                            .collect(Collectors.toList());
                    
                    return TagCloudResponse.builder().items(items).build();
                });
    }
    
    /**
     * 知识搜索
     */
    public Mono<KnowledgeSearchResponse> searchKnowledge(KnowledgeSearchRequest request) {
        Integer limit = request.getLimit() != null ? request.getLimit() : 5;
        
        return openAIClient.createEmbedding(request.getQuery())
                .flatMapMany(questionEmbedding -> 
                    qdrantClient.searchPoints(questionEmbedding, limit)
                )
                .flatMap(result -> {
                    Long docId = ((Number) result.getPayload().get("document_id")).longValue();
                    Long chunkId = ((Number) result.getPayload().get("chunk_id")).longValue();
                    String content = (String) result.getPayload().get("content");
                    String title = (String) result.getPayload().getOrDefault("title", "");
                    String fileName = (String) result.getPayload().getOrDefault("file_name", "");
                    
                    String snippet = content.length() > 200 ? 
                            content.substring(0, 200) + "..." : content;
                    
                    return Mono.just(KnowledgeSearchResponse.KnowledgeSearchResult.builder()
                            .documentId(docId)
                            .chunkId(chunkId)
                            .title(title)
                            .content(content)
                            .snippet(snippet)
                            .score(result.getScore())
                            .fileName(fileName)
                            .build());
                })
                .collectList()
                .map(items -> KnowledgeSearchResponse.builder().items(items).build())
                .onErrorResume(e -> {
                    log.error("Knowledge search failed", e);
                    return Mono.just(KnowledgeSearchResponse.builder()
                            .items(new ArrayList<>())
                            .build());
                });
    }
    
    /**
     * 与指定文档对话
     */
    public Mono<ChatDocumentResponse> chatWithDocument(Long documentId, ChatDocumentRequest request) {
        return documentRepository.findById(documentId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("文档不存在")))
                .flatMap(document -> {
                    // 获取文档的chunks
                    return documentChunkRepository.findByDocumentIdOrderByIndexAsc(documentId)
                            .map(DocumentChunk::getContent)
                            .collectList()
                            .flatMap(contents -> {
                                if (contents.isEmpty()) {
                                    return Mono.error(new BusinessException("文档内容为空"));
                                }
                                
                                return openAIClient.chat(request.getQuestion(), contents)
                                        .map(answer -> ChatDocumentResponse.builder()
                                                .answer(answer)
                                                .sources(List.of(ChatDocumentResponse.ChatDocumentSource.builder()
                                                        .documentId(document.getId())
                                                        .title(document.getTitle())
                                                        .filePath(document.getFilePath())
                                                        .build()))
                                                .build());
                            });
                });
    }
    
    /**
     * 流式问答 - 返回字符流
     */
    public Flux<String> chatWithDocumentsStream(String question, Long spaceId) {
        return openAIClient.createEmbedding(question)
                .flatMapMany(questionEmbedding -> 
                    qdrantClient.searchPoints(questionEmbedding, 5)
                )
                .collectList()
                .flatMapMany(results -> {
                    if (results.isEmpty()) {
                        return Flux.error(new BusinessException("未找到相关文档"));
                    }
                    
                    // Extract content from search results
                    List<String> contexts = results.stream()
                            .map(result -> (String) result.getPayload().get("content"))
                            .collect(Collectors.toList());
                    
                    // 调用chat并将结果分割为流
                    return openAIClient.chat(question, contexts)
                            .flatMapMany(answer -> 
                                Flux.fromArray(answer.split("(?<=\\S)(?=\\s)|(?<=\\s)(?=\\S)"))
                            );
                });
    }
    
    /**
     * 预览文档
     */
    public Mono<Void> previewDocument(Long id, org.springframework.http.server.reactive.ServerHttpResponse response) {
        return documentRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("文档不存在")))
                .flatMap(document -> {
                    // 设置响应头
                    response.getHeaders().add("Content-Type", getContentType(document.getFileType()));
                    response.getHeaders().add("Content-Disposition", 
                            "inline; filename=\"" + document.getFileName() + "\"");
                    
                    return minioClient.downloadFile(document.getFilePath())
                            .flatMap(inputStream -> {
                                return org.springframework.core.io.buffer.DataBufferUtils
                                        .readInputStream(
                                            () -> inputStream,
                                            response.bufferFactory(),
                                            4096
                                        )
                                        .as(response::writeWith);
                            });
                });
    }
    
    private String getContentType(String fileType) {
        return switch (fileType) {
            case ".pdf" -> "application/pdf";
            case ".jpg", ".jpeg" -> "image/jpeg";
            case ".png" -> "image/png";
            case ".txt" -> "text/plain";
            default -> "application/octet-stream";
        };
    }
    
    public Mono<Void> deleteDocument(Long id) {
        return documentRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")))
                .flatMap(document -> 
                    minioClient.deleteFile(document.getFilePath())
                            .then(documentChunkRepository.deleteByDocumentId(id))
                            .then(documentRepository.delete(document))
                );
    }
    
    public Mono<Document> publishDocument(Long id) {
        return documentRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")))
                .flatMap(document -> {
                    if (!Document.STATUS_PENDING_PUBLISH.equals(document.getStatus())) {
                        return Mono.error(new BusinessException("Document is not ready to publish"));
                    }
                    document.setStatus(Document.STATUS_PUBLISHED);
                    document.setUpdatedAt(LocalDateTime.now());
                    return documentRepository.save(document);
                });
    }
}

