package com.knowledgebase.platformspring.controller;

import org.springframework.http.codec.multipart.FilePart;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.model.Document;
import com.knowledgebase.platformspring.repository.UserRepository;
import com.knowledgebase.platformspring.service.DocumentService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.Data;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Tag(name = "文档管理", description = "文档上传、查询、问答相关接口")
@RestController
@RequestMapping("/api/v1/documents")
@RequiredArgsConstructor
@SecurityRequirement(name = "bearerAuth")
public class DocumentController {
    
    private final DocumentService documentService;
    private final UserRepository userRepository;
    
    @Operation(summary = "上传文档", description = "上传文档文件并自动进行OCR和向量化处理")
    @PostMapping("/upload")
    public Mono<ApiResponse<Document>> uploadDocument(
            @RequestPart("file") FilePart filePart,
            @RequestParam Long spaceId,
            @RequestParam Long subSpaceId,
            @RequestParam Long classId,
            Authentication authentication) {
        
        Long userId = (Long) authentication.getPrincipal();
        
        return userRepository.findById(userId)
                .flatMap(user -> documentService.uploadDocument(
                        filePart, spaceId, subSpaceId, classId, userId, user.getNickname()))
                .map(doc -> ApiResponse.success("Document uploaded successfully", doc));
    }
    
    @Operation(summary = "获取空间文档", description = "获取指定空间下的所有文档")
    @GetMapping("/space/{spaceId}")
    public Flux<Document> getDocumentsBySpace(@PathVariable Long spaceId) {
        return documentService.getDocumentsBySpaceId(spaceId);
    }
    
    @Operation(summary = "获取文档详情", description = "根据ID获取文档的详细信息")
    @GetMapping("/{id}")
    public Mono<ApiResponse<Document>> getDocument(@PathVariable Long id) {
        return documentService.getDocumentById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "文档问答", description = "基于知识库文档内容进行智能问答")
    @PostMapping("/chat")
    public Mono<ApiResponse<ChatResponse>> chat(@RequestBody ChatRequest request) {
        return documentService.chatWithDocuments(request.getQuestion(), request.getSpaceId())
                .map(answer -> ApiResponse.success(new ChatResponse(answer)));
    }
    
    @Operation(summary = "发布文档", description = "将文档状态设置为已发布，使其可被搜索和使用")
    @PostMapping("/{id}/publish")
    public Mono<ApiResponse<Document>> publishDocument(@PathVariable Long id) {
        return documentService.publishDocument(id)
                .map(doc -> ApiResponse.success("Document published successfully", doc));
    }
    
    @Operation(summary = "删除文档", description = "删除指定文档及其关联的向量数据")
    @DeleteMapping("/{id}")
    public Mono<ApiResponse<Void>> deleteDocument(@PathVariable Long id) {
        return documentService.deleteDocument(id)
                .then(Mono.just(ApiResponse.<Void>success("Document deleted successfully", null)));
    }
    
    @Data
    public static class ChatRequest {
        private String question;
        private Long spaceId;
    }
    
    @Data
    @RequiredArgsConstructor
    public static class ChatResponse {
        private final String answer;
    }
}

