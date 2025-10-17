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

import lombok.Data;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v1/documents")
@RequiredArgsConstructor
public class DocumentController {
    
    private final DocumentService documentService;
    private final UserRepository userRepository;
    
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
    
    @GetMapping("/space/{spaceId}")
    public Flux<Document> getDocumentsBySpace(@PathVariable Long spaceId) {
        return documentService.getDocumentsBySpaceId(spaceId);
    }
    
    @GetMapping("/{id}")
    public Mono<ApiResponse<Document>> getDocument(@PathVariable Long id) {
        return documentService.getDocumentById(id)
                .map(ApiResponse::success);
    }
    
    @PostMapping("/chat")
    public Mono<ApiResponse<ChatResponse>> chat(@RequestBody ChatRequest request) {
        return documentService.chatWithDocuments(request.getQuestion(), request.getSpaceId())
                .map(answer -> ApiResponse.success(new ChatResponse(answer)));
    }
    
    @PostMapping("/{id}/publish")
    public Mono<ApiResponse<Document>> publishDocument(@PathVariable Long id) {
        return documentService.publishDocument(id)
                .map(doc -> ApiResponse.success("Document published successfully", doc));
    }
    
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

