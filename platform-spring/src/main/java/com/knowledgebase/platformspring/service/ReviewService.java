package com.knowledgebase.platformspring.service;

import java.io.ByteArrayInputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;
import java.util.concurrent.TimeoutException;
import java.util.stream.Collectors;

import org.springframework.http.codec.multipart.FilePart;
import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.client.MinioClientService;
import com.knowledgebase.platformspring.client.OpenAIClientService;
import com.knowledgebase.platformspring.client.PaddleOCRClientService;
import com.knowledgebase.platformspring.client.QdrantClientService;
import com.knowledgebase.platformspring.dto.AcceptSuggestionRequest;
import com.knowledgebase.platformspring.dto.DocumentSection;
import com.knowledgebase.platformspring.dto.ReviewRequest;
import com.knowledgebase.platformspring.dto.ReviewSuggestion;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.util.DocumentParser;
import com.knowledgebase.platformspring.util.FormatChecker;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Schedulers;

/**
 * 智能审查服务
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class ReviewService {

    private final MinioClientService minioClient;
    private final OpenAIClientService openAIClient;
    private final QdrantClientService qdrantClient;
    private final PaddleOCRClientService ocrClient;

    // 存储会话的建议映射 (sessionId -> List<ReviewSuggestion>)
    private final Map<String, List<ReviewSuggestion>> sessionSuggestions = new HashMap<>();

    /**
     * 上传文档用于审查
     */
    public Mono<String> uploadForReview(FilePart filePart, String fileName) {
        String sessionId = UUID.randomUUID().toString();
        String tempPath = "review-temp/" + sessionId + "/" + fileName;

        log.info("Uploading file for review: sessionId={}, fileName={}", sessionId, fileName);

        return filePart.content()
                .reduce((buffer1, buffer2) -> buffer1.write(buffer2))
                .flatMap(dataBuffer -> {
                    byte[] bytes = new byte[dataBuffer.readableByteCount()];
                    dataBuffer.read(bytes);

                    return minioClient.uploadFile(
                            tempPath,
                            new ByteArrayInputStream(bytes),
                            bytes.length,
                            Objects.requireNonNull(filePart.headers().getContentType()).toString()
                    ).thenReturn(sessionId);
                })
                .doOnSuccess(id -> log.info("File uploaded for review: sessionId={}", id))
                .doOnError(e -> log.error("Failed to upload file for review: {}", e.getMessage(), e));
    }

    /**
     * 智能审查文档（流式返回建议）
     */
    public Flux<ReviewSuggestion> reviewDocument(ReviewRequest request) {
        log.info("Starting document review: file={}", request.getFileName());

        return extractText(request)
                .flatMapMany(text -> {
                    log.debug("Text extracted, length: {}", text.length());

                    // 解析文档为段落
                    List<DocumentSection> sections = DocumentParser.parse(text);
                    log.debug("Document parsed into {} sections", sections.size());

                    // 先发送一条特殊建议，包含文档内容（用于前端显示）
                    ReviewSuggestion contentMarker = ReviewSuggestion.builder()
                            .id(UUID.randomUUID().toString())
                            .type("DOCUMENT_CONTENT")
                            .severity("INFO")
                            .position(-1)
                            .documentContent(text)
                            .originalText("")
                            .suggestedText("")
                            .reason("文档内容")
                            .build();

                    // 对每个段落并行检查
                    return Flux.just(contentMarker)
                            .concatWith(Flux.fromIterable(sections)
                                    .flatMap(section -> {
                                        ReviewRequest.ReviewOptions options = request.getOptions();
                                        if (options == null) {
                                            options = ReviewRequest.ReviewOptions.builder().build();
                                        }

                                        List<Flux<ReviewSuggestion>> checks = new ArrayList<>();

                                        // 格式检查
                                        if (Boolean.TRUE.equals(options.getCheckFormat())) {
                                            checks.add(checkFormat(section));
                                        }

                                        // 引用验证
                                        if (Boolean.TRUE.equals(options.getVerifyReferences())) {
                                            checks.add(verifyReferences(section, request.getSpaceId()));
                                        }

                                        // 内容建议
                                        if (Boolean.TRUE.equals(options.getSuggestContent())) {
                                            checks.add(suggestContent(section, sections, request.getSpaceId()));
                                        }

                                        return Flux.merge(checks);
                                    }, 3)) // 并发度为3
                            .doOnNext(suggestion -> {
                                // 为每个建议生成ID（如果还没有的话）
                                if (suggestion.getId() == null) {
                                    suggestion.setId(UUID.randomUUID().toString());
                                }

                                // 存储建议到会话映射中
                                String sessionId = extractSessionIdFromPath(request.getTempFilePath());
                                sessionSuggestions.computeIfAbsent(sessionId, k -> new ArrayList<>()).add(suggestion);
                            });
                })
                .doOnComplete(() -> log.info("Document review completed"))
                .doOnError(e -> log.error("Document review failed: {}", e.getMessage(), e));
    }

    /**
     * 提取文本内容
     */
    private Mono<String> extractText(ReviewRequest request) {
        String fileType = request.getFileType().toLowerCase();
        String tempPath = request.getTempFilePath();

        return minioClient.downloadFile(tempPath)
                .flatMap(inputStream -> {
                    try {
                        byte[] bytes = inputStream.readAllBytes();
                        inputStream.close();

                        // 根据文件类型提取文本
                        return switch (fileType) {
                            case ".txt", ".md" -> Mono.just(new String(bytes));
                            case ".pdf", ".doc", ".docx" -> {
                                // 使用Pandoc转Markdown
                                log.debug("Using Pandoc for file type: {}", fileType);
                                yield convertToPandoc(request.getFileName(), bytes, fileType);
                            }
                            case ".jpg", ".jpeg", ".png" -> {
                                // 图片才用OCR
                                log.debug("Using OCR for image file: {}", fileType);
                                yield ocrClient.recognize(request.getFileName(), bytes);
                            }
                            default -> Mono.error(new BusinessException("Unsupported file type: " + fileType));
                        };
                    } catch (Exception e) {
                        return Mono.error(new BusinessException("Failed to read file: " + e.getMessage()));
                    }
                });
    }

    /**
     * 使用Pandoc转换为Markdown
     */
    private Mono<String> convertToPandoc(String fileName, byte[] data, String fileType) {
        return Mono.fromCallable(() -> {
            try {
                // 创建临时文件
                java.io.File tempFile = java.io.File.createTempFile("review-", fileType);
                java.nio.file.Files.write(tempFile.toPath(), data);

                // 执行pandoc命令
                ProcessBuilder pb = new ProcessBuilder(
                        "pandoc",
                        tempFile.getAbsolutePath(),
                        "-t", "markdown",
                        "--wrap=none"
                );

                Process process = pb.start();

                // 读取输出
                String markdown = new String(process.getInputStream().readAllBytes());

                // 等待完成
                int exitCode = process.waitFor();

                // 删除临时文件
                tempFile.delete();

                if (exitCode != 0) {
                    String error = new String(process.getErrorStream().readAllBytes());
                    throw new RuntimeException("Pandoc conversion failed: " + error);
                }

                log.debug("Pandoc conversion successful, markdown length: {}", markdown.length());
                return markdown;

            } catch (Exception e) {
                log.error("Pandoc conversion error: {}", e.getMessage());
                throw new RuntimeException("Failed to convert with Pandoc: " + e.getMessage());
            }
        }).subscribeOn(Schedulers.boundedElastic());
    }

    /**
     * 格式规范检查
     */
    private Flux<ReviewSuggestion> checkFormat(DocumentSection section) {
        return Mono.fromCallable(() -> FormatChecker.check(section))
                .flatMapMany(Flux::fromIterable)
                .map(suggestion -> {
                    // 为格式检查建议生成ID
                    if (suggestion.getId() == null) {
                        suggestion.setId(UUID.randomUUID().toString());
                    }
                    return suggestion;
                })
                .subscribeOn(Schedulers.boundedElastic());
    }

    /**
     * 验证知识库引用
     */
    private Flux<ReviewSuggestion> verifyReferences(DocumentSection section, Long spaceId) {
        // 提取法规引用
        List<String> references = DocumentParser.extractReferences(section.getContent());

        if (references.isEmpty()) {
            return Flux.empty();
        }

        log.debug("Found {} references in section {}", references.size(), section.getPosition());

        return Flux.fromIterable(references)
                .flatMap(ref ->
                        // 生成embedding并检索知识库
                        openAIClient.createEmbedding(ref)
                                .timeout(java.time.Duration.ofSeconds(30)) // Embedding超时30秒
                                .flatMapMany(embedding ->
                                        qdrantClient.searchPoints(embedding, 3)
                                )
                                .timeout(java.time.Duration.ofSeconds(10)) // Qdrant检索超时10秒
                                .filter(result -> result.getScore() > 0.85) // 高相关度
                                .take(1) // 只取最相关的
                                .map(result -> {
                                    String knowledgeContent = (String) result.getPayload().get("content");
                                    String knowledgeTitle = (String) result.getPayload().get("title");
                                    Long docId = ((Number) result.getPayload().get("document_id")).longValue();

                                    // 简单的版本对比：如果知识库中的内容与引用不完全一致，可能有更新
                                    if (!knowledgeContent.contains(ref)) {
                                        return ReviewSuggestion.builder()
                                                .id(UUID.randomUUID().toString())
                                                .type(ReviewSuggestion.TYPE_REFERENCE_OUTDATED)
                                                .severity(ReviewSuggestion.SEVERITY_WARNING)
                                                .position(section.getPosition())
                                                .originalText(ref)
                                                .suggestedText("建议核对：" + knowledgeTitle)
                                                .reason("知识库中可能存在更新版本或相关规定")
                                                .knowledgeSource(knowledgeTitle)
                                                .knowledgeDocumentId(docId)
                                                .build();
                                    }
                                    return null;
                                })
                                .filter(Objects::nonNull)
                )
                .onErrorResume(e -> {
                    if (e instanceof TimeoutException) {
                        log.warn("Reference verification timeout for section {}, skipping",
                                section.getPosition());
                    } else {
                        log.warn("Reference verification failed for section {}: {}",
                                section.getPosition(), e.getMessage());
                    }
                    return Flux.empty(); // 失败时返回空，不中断整个流程
                });
    }

    /**
     * 基于知识库的内容建议
     */
    private Flux<ReviewSuggestion> suggestContent(DocumentSection section,
                                                  List<DocumentSection> allSections,
                                                  Long spaceId) {
        // 跳过标题、日期等段落
        if (!DocumentSection.TYPE_PARAGRAPH.equals(section.getType()) &&
                !DocumentSection.TYPE_ARTICLE.equals(section.getType())) {
            return Flux.empty();
        }

        // 合并上下文（前后各2个段落）
        String context = DocumentParser.mergeContext(allSections, section.getPosition(), 2);

        return openAIClient.createEmbedding(section.getContent())
                .timeout(java.time.Duration.ofSeconds(30)) // Embedding超时30秒
                .flatMapMany(embedding ->
                        qdrantClient.searchPoints(embedding, 5)
                )
                .timeout(java.time.Duration.ofSeconds(10)) // Qdrant检索超时10秒
                .collectList()
                .flatMapMany(results -> {
                    if (results.isEmpty()) {
                        log.debug("No knowledge found for section {}, skipping AI suggestion",
                                section.getPosition());
                        return Flux.empty();
                    }

                    // 提取知识库内容
                    List<String> knowledgeContext = results.stream()
                            .map(r -> (String) r.getPayload().get("content"))
                            .collect(Collectors.toList());

                    // 构造Prompt
                    String prompt = buildReviewPrompt(section.getContent(), context);

                    // AI生成建议 - 设置较短的超时时间
                    return openAIClient.chat(prompt, knowledgeContext)
                            .timeout(java.time.Duration.ofSeconds(45)) // AI生成超时45秒
                            .flux()
                            .filter(suggestion -> suggestion != null && !suggestion.trim().isEmpty())
                            .filter(suggestion -> !suggestion.contains("无建议")) // 过滤掉"无建议"
                            .map(suggestion -> ReviewSuggestion.builder()
                                    .id(UUID.randomUUID().toString())
                                    .type(ReviewSuggestion.TYPE_CONTENT_ENHANCEMENT)
                                    .severity(ReviewSuggestion.SEVERITY_INFO)
                                    .position(section.getPosition())
                                    .originalText(section.getContent())
                                    .suggestedText(suggestion)
                                    .reason("基于知识库的内容补充建议")
                                    .build());
                })
                .onErrorResume(e -> {
                    if (e instanceof TimeoutException) {
                        log.warn("Content suggestion timeout for section {}, skipping",
                                section.getPosition());
                    } else {
                        log.warn("Content suggestion failed for section { {}",
                                section.getPosition(), e.getMessage());
                    }
                    return Flux.empty(); // 失败时返回空，不中断整个流程
                });
    }

    /**
     * 构造审查Prompt
     */
    private String buildReviewPrompt(String content, String context) {
        return String.format("""
                你是一个专业的公文审查助手。请审查以下公文段落，并根据知识库内容提出改进建议。
                
                当前段落：
                %s
                
                上下文：
                %s
                
                请提供具体的修改建议，包括：
                1. 是否有遗漏的重要内容
                2. 表述是否准确、规范
                3. 是否需要引用相关法规或政策
                
                如果没有需要改进的地方，请回复"无建议"。
                """, content, context);
    }

    /**
     * 获取会话的临时文件路径
     */
    public String getTempFilePath(String sessionId, String fileName) {
        return "review-temp/" + sessionId + "/" + fileName;
    }

    /**
     * 从路径中提取sessionId
     */
    private String extractSessionIdFromPath(String tempFilePath) {
        // 路径格式: review-temp/{sessionId}/{fileName}
        String[] parts = tempFilePath.split("/");
        if (parts.length >= 3) {
            return parts[1];
        }
        throw new BusinessException("Invalid temp file path format: " + tempFilePath);
    }

    /**
     * 接受建议并更新文档
     */
    public Mono<String> acceptSuggestions(AcceptSuggestionRequest request) {
        log.info("Accepting suggestions for session: {}, suggestions: {}",
                request.getSessionId(), request.getAcceptedSuggestionIds());

        // 获取会话的建议
        List<ReviewSuggestion> suggestions = sessionSuggestions.get(request.getSessionId());
        if (suggestions == null || suggestions.isEmpty()) {
            return Mono.error(new BusinessException("No suggestions found for session: " + request.getSessionId()));
        }

        // 获取要接受的建议
        List<ReviewSuggestion> acceptedSuggestions;
        if (Boolean.TRUE.equals(request.getApplyAll())) {
            // 接受所有建议（除了文档内容标记）
            acceptedSuggestions = suggestions.stream()
                    .filter(s -> !"DOCUMENT_CONTENT".equals(s.getType()))
                    .collect(Collectors.toList());
        } else {
            // 接受指定的建议
            acceptedSuggestions = suggestions.stream()
                    .filter(s -> request.getAcceptedSuggestionIds().contains(s.getId()))
                    .collect(Collectors.toList());
        }

        if (acceptedSuggestions.isEmpty()) {
            return Mono.error(new BusinessException("No valid suggestions to accept"));
        }

        // 获取原始文档内容
        String tempFilePath = getTempFilePath(request.getSessionId(), request.getFileName());

        return minioClient.downloadFile(tempFilePath)
                .flatMap(inputStream -> {
                    try {
                        byte[] originalBytes = inputStream.readAllBytes();
                        inputStream.close();

                        // 根据文件类型处理
                        return switch (request.getFileType().toLowerCase()) {
                            case ".txt", ".md" -> {
                                // 文本文件直接修改内容
                                String originalText = new String(originalBytes);
                                String modifiedText = applyTextSuggestions(originalText, acceptedSuggestions);
                                yield uploadModifiedFile(tempFilePath, modifiedText.getBytes(), "text/plain");
                            }
                            case ".pdf", ".doc", ".docx" -> {
                                // 对于复杂格式，先转换为Markdown，修改后再转换回原格式
                                yield applyComplexFormatSuggestions(tempFilePath, originalBytes, request.getFileType(), acceptedSuggestions);
                            }
                            default -> Mono.error(new BusinessException("Unsupported file type for modification: " + request.getFileType()));
                        };
                    } catch (Exception e) {
                        return Mono.error(new BusinessException("Failed to process file: " + e.getMessage()));
                    }
                })
                .doOnSuccess(result -> {
                    log.info("Successfully applied {} suggestions to session {}",
                            acceptedSuggestions.size(), request.getSessionId());
                })
                .doOnError(e -> log.error("Failed to apply suggestions: {}", e.getMessage(), e));
    }

    /**
     * 应用文本建议
     */
    private String applyTextSuggestions(String originalText, List<ReviewSuggestion> suggestions) {
        String modifiedText = originalText;

        // 按位置倒序排序，避免位置偏移问题
        List<ReviewSuggestion> sortedSuggestions = suggestions.stream()
                .filter(s -> s.getPosition() != null && s.getPosition() >= 0)
                .sorted((a, b) -> Integer.compare(b.getPosition(), a.getPosition()))
                .collect(Collectors.toList());

        for (ReviewSuggestion suggestion : sortedSuggestions) {
            if (suggestion.getOriginalText() != null && suggestion.getSuggestedText() != null) {
                // 替换原始文本为建议文本
                modifiedText = modifiedText.replace(suggestion.getOriginalText(), suggestion.getSuggestedText());
            }
        }

        return modifiedText;
    }

    /**
     * 应用复杂格式文件的建议
     */
    private Mono<String> applyComplexFormatSuggestions(String tempFilePath, byte[] originalBytes,
                                                       String fileType, List<ReviewSuggestion> suggestions) {
        return Mono.fromCallable(() -> {
                    try {
                        // 创建临时文件
                        java.io.File tempFile = java.io.File.createTempFile("review-", fileType);
                        java.nio.file.Files.write(tempFile.toPath(), originalBytes);

                        // 使用Pandoc转换为Markdown
                        ProcessBuilder toMarkdown = new ProcessBuilder(
                                "pandoc",
                                tempFile.getAbsolutePath(),
                                "-t", "markdown",
                                "--wrap=none"
                        );

                        Process process = toMarkdown.start();
                        String markdown = new String(process.getInputStream().readAllBytes());
                        int exitCode = process.waitFor();

                        if (exitCode != 0) {
                            throw new RuntimeException("Pandoc conversion to markdown failed");
                        }

                        // 应用建议
                        String modifiedMarkdown = applyTextSuggestions(markdown, suggestions);

                        // 创建修改后的Markdown临时文件
                        java.io.File modifiedMarkdownFile = java.io.File.createTempFile("modified-", ".md");
                        java.nio.file.Files.write(modifiedMarkdownFile.toPath(), modifiedMarkdown.getBytes());

                        // 转换回原格式
                        ProcessBuilder fromMarkdown = new ProcessBuilder(
                                "pandoc",
                                modifiedMarkdownFile.getAbsolutePath(),
                                "-f", "markdown",
                                "-t", getPandocOutputFormat(fileType),
                                "-o", tempFile.getAbsolutePath()
                        );

                        process = fromMarkdown.start();
                        exitCode = process.waitFor();

                        // 清理临时文件
                        modifiedMarkdownFile.delete();

                        if (exitCode != 0) {
                            throw new RuntimeException("Pandoc conversion from markdown failed");
                        }

                        // 读取修改后的文件
                        byte[] modifiedBytes = java.nio.file.Files.readAllBytes(tempFile.toPath());
                        tempFile.delete();

                        return modifiedBytes;

                    } catch (Exception e) {
                        log.error("Complex format suggestion application failed: {}", e.getMessage());
                        throw new RuntimeException("Failed to apply suggestions to complex format: " + e.getMessage());
                    }
                })
                .subscribeOn(Schedulers.boundedElastic())
                .flatMap(modifiedBytes -> uploadModifiedFile(tempFilePath, modifiedBytes, getContentType(fileType)));
    }

    /**
     * 获取Pandoc输出格式
     */
    private String getPandocOutputFormat(String fileType) {
        return switch (fileType.toLowerCase()) {
            case ".pdf" -> "pdf";
            case ".doc" -> "docx";
            case ".docx" -> "docx";
            default -> "plain";
        };
    }

    /**
     * 获取内容类型
     */
    private String getContentType(String fileType) {
        return switch (fileType.toLowerCase()) {
            case ".pdf" -> "application/pdf";
            case ".doc" -> "application/msword";
            case ".docx" -> "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
            case ".txt" -> "text/plain";
            case ".md" -> "text/markdown";
            default -> "application/octet-stream";
        };
    }

    /**
     * 上传修改后的文件
     */
    private Mono<String> uploadModifiedFile(String tempFilePath, byte[] modifiedBytes, String contentType) {
        return minioClient.uploadFile(
                tempFilePath,
                new ByteArrayInputStream(modifiedBytes),
                modifiedBytes.length,
                contentType
        ).thenReturn(tempFilePath);
    }

    /**
     * 下载文档
     */
    public Mono<java.io.InputStream> downloadDocument(String sessionId, String fileName, String fileType) {
        log.info("Downloading document: sessionId={}, fileName={}", sessionId, fileName);

        String tempFilePath = getTempFilePath(sessionId, fileName);

        return minioClient.downloadFile(tempFilePath)
                .doOnSuccess(inputStream -> log.info("Document downloaded successfully: sessionId={}", sessionId))
                .doOnError(e -> log.error("Failed to download document: sessionId={}, error={}", sessionId, e.getMessage(), e));
    }
}

