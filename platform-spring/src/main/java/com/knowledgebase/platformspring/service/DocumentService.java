package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;
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
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.Document;
import com.knowledgebase.platformspring.model.DocumentChunk;
import com.knowledgebase.platformspring.repository.DocumentChunkRepository;
import com.knowledgebase.platformspring.repository.DocumentRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Slf4j
@Service
@RequiredArgsConstructor
public class DocumentService {
    
    private final DocumentRepository documentRepository;
    private final DocumentChunkRepository documentChunkRepository;
    private final MinioClientService minioClient;
    private final OpenAIClientService openAIClient;
    @SuppressWarnings("unused") // Reserved for future OCR implementation
    private final PaddleOCRClientService ocrClient;
    private final QdrantClientService qdrantClient;
    
    public Mono<Document> uploadDocument(FilePart filePart, Long spaceId, Long subSpaceId, 
                                        Long classId, Long userId, String nickName) {
        
        String fileName = filePart.filename();
        String objectName = "documents/" + UUID.randomUUID().toString() + "_" + fileName;
        
        Document document = Document.builder()
                .title(fileName)
                .fileName(fileName)
                .filePath(objectName)
                .spaceId(spaceId)
                .subSpaceId(subSpaceId)
                .classId(classId)
                .createdBy(userId)
                .creatorNickName(nickName)
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
                .flatMap(doc -> processDocument(doc.getId()));
    }
    
    public Mono<Document> processDocument(Long documentId) {
        return documentRepository.findById(documentId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")))
                .flatMap(document -> {
                    document.setStatus(Document.STATUS_PROCESSING);
                    return documentRepository.save(document);
                })
                .flatMap(this::extractText)
                .flatMap(this::chunkAndVectorize)
                .flatMap(document -> {
                    document.setStatus(Document.STATUS_PENDING_PUBLISH);
                    document.setProcessedAt(LocalDateTime.now());
                    return documentRepository.save(document);
                })
                .onErrorResume(e -> {
                    log.error("Failed to process document: {}", documentId, e);
                    return documentRepository.findById(documentId)
                            .flatMap(doc -> {
                                doc.setStatus(Document.STATUS_PROCESS_FAILED);
                                doc.setParseError(e.getMessage());
                                doc.setRetryCount(doc.getRetryCount() + 1);
                                return documentRepository.save(doc);
                            });
                });
    }
    
    private Mono<Document> extractText(Document document) {
        // For now, just return the document
        // TODO: Implement OCR and text extraction
        document.setContent("Extracted text content");
        return Mono.just(document);
    }
    
    private Mono<Document> chunkAndVectorize(Document document) {
        if (document.getContent() == null || document.getContent().isEmpty()) {
            return Mono.just(document);
        }
        
        // Simple chunking: split by paragraphs
        String[] chunks = document.getContent().split("\n\n");
        
        return Flux.fromArray(chunks)
                .index()
                .flatMap(tuple -> {
                    String chunkContent = tuple.getT2();
                    int index = tuple.getT1().intValue();
                    
                    return openAIClient.createEmbedding(chunkContent)
                            .flatMap(embedding -> {
                                // Save to Qdrant
                                Map<String, Object> payload = new HashMap<>();
                                payload.put("document_id", document.getId());
                                payload.put("chunk_index", index);
                                payload.put("content", chunkContent);
                                
                                QdrantClientService.QdrantPoint point = 
                                    QdrantClientService.QdrantPoint.create(embedding, payload);
                                
                                return qdrantClient.upsertPoints(List.of(point))
                                        .then(Mono.fromCallable(() -> {
                                            DocumentChunk chunk = DocumentChunk.builder()
                                                    .documentId(document.getId())
                                                    .chunkIndex(index)
                                                    .content(chunkContent)
                                                    .vectorId(point.getId())
                                                    .createdAt(LocalDateTime.now())
                                                    .build();
                                            return chunk;
                                        }));
                            })
                            .flatMap(documentChunkRepository::save);
                })
                .collectList()
                .map(chunks1 -> {
                    document.setVectorCount(chunks1.size());
                    return document;
                });
    }
    
    public Flux<Document> getDocumentsBySpaceId(Long spaceId) {
        return documentRepository.findBySpaceId(spaceId);
    }
    
    public Mono<Document> getDocumentById(Long id) {
        return documentRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Document not found")));
    }
    
    public Mono<String> chatWithDocuments(String question, Long spaceId) {
        return openAIClient.createEmbedding(question)
                .flatMapMany(questionEmbedding -> 
                    qdrantClient.searchPoints(questionEmbedding, 5)
                )
                .collectList()
                .flatMap(results -> {
                    if (results.isEmpty()) {
                        return Mono.error(new BusinessException("No relevant documents found"));
                    }
                    
                    // Extract content from search results
                    List<String> contexts = results.stream()
                            .map(result -> (String) result.getPayload().get("content"))
                            .collect(Collectors.toList());
                    
                    return openAIClient.chat(question, contexts);
                });
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

