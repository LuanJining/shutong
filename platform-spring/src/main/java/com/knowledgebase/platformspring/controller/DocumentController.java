package com.knowledgebase.platformspring.controller;

import java.time.Duration;
import java.util.List;

import org.springframework.http.MediaType;
import org.springframework.http.codec.ServerSentEvent;
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
import com.knowledgebase.platformspring.dto.ChatDocumentRequest;
import com.knowledgebase.platformspring.dto.ChatDocumentResponse;
import com.knowledgebase.platformspring.dto.HomepageResponse;
import com.knowledgebase.platformspring.dto.KnowledgeSearchRequest;
import com.knowledgebase.platformspring.dto.KnowledgeSearchResponse;
import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.dto.TagCloudResponse;
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
    @PostMapping(value = "/upload", consumes = "multipart/form-data")
    public Mono<ApiResponse<Document>> uploadDocument(
            @RequestPart("file") FilePart filePart,
            @RequestPart("space_id") String spaceIdStr,
            @RequestPart("sub_space_id") String subSpaceIdStr,
            @RequestPart("class_id") String classIdStr,
            @RequestPart(value = "file_name", required = false) String fileName,
            @RequestPart(value = "tags", required = false) String tags,
            @RequestPart(value = "summary", required = false) String summary,
            @RequestPart(value = "department", required = false) String department,
            @RequestPart(value = "need_approval", required = false) String needApprovalStr,
            @RequestPart(value = "version", required = false) String version,
            @RequestPart(value = "use_type", required = false) String useType,
            Authentication authentication) {
        
        Long userId = (Long) authentication.getPrincipal();
        
        // 转换必需的 ID 字段
        Long spaceId = Long.parseLong(spaceIdStr);
        Long subSpaceId = Long.parseLong(subSpaceIdStr);
        Long classId = Long.parseLong(classIdStr);
        
        // 如果没有提供file_name，使用原始文件名
        String actualFileName = (fileName != null && !fileName.isEmpty()) ? 
                fileName : filePart.filename();
        
        // 处理 Boolean 类型转换
        Boolean needApproval = (needApprovalStr != null && !needApprovalStr.isEmpty()) ? 
                Boolean.parseBoolean(needApprovalStr) : false;
        
        // 设置默认值
        String finalVersion = (version != null && !version.isEmpty()) ? version : "v1.0.0";
        String finalUseType = (useType != null && !useType.isEmpty()) ? useType : "viewable";
        
        return userRepository.findById(userId)
                .flatMap(user -> documentService.uploadDocument(
                        filePart, spaceId, subSpaceId, classId, userId, user.getNickname(),
                        actualFileName, tags, summary, 
                        department != null ? department : user.getDepartment(),
                        needApproval, finalVersion, finalUseType))
                .map(doc -> ApiResponse.<Document>success("文档上传成功", doc));
    }
    
    @Operation(summary = "获取文档详情", description = "根据ID获取文档的详细信息")
    @GetMapping("/{id}/info")
    public Mono<ApiResponse<Document>> getDocument(@PathVariable Long id) {
        return documentService.getDocumentById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取空间文档列表", description = "根据空间ID获取该空间的所有文档")
    @GetMapping("/space/{spaceId}")
    public Mono<ApiResponse<PaginationResponse<List<Document>>>> getDocumentsBySpaceId(
            @PathVariable Long spaceId,
            @RequestParam(defaultValue = "1") Integer page,
            @RequestParam(defaultValue = "10") Integer pageSize) {
        return documentService.getDocumentsBySpaceIdPaginated(spaceId, page, pageSize)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取空间文档", description = "获取指定文档所属空间的所有文档")
    @GetMapping("/{id}/space")
    public Mono<ApiResponse<PaginationResponse<List<Document>>>> getDocumentsBySpace(
            @PathVariable Long id,
            @RequestParam(defaultValue = "1") Integer page,
            @RequestParam(defaultValue = "10") Integer pageSize) {
        return documentService.getDocumentById(id)
                .flatMap(doc -> documentService.getDocumentsBySpaceIdPaginated(doc.getSpaceId(), page, pageSize))
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "预览文档", description = "在线预览文档内容")
    @GetMapping("/{id}/preview")
    public Mono<Void> previewDocument(@PathVariable Long id, 
                                      org.springframework.http.server.reactive.ServerHttpResponse response) {
        return documentService.previewDocument(id, response);
    }
    
    @Operation(summary = "获取首页文档", description = "获取首页展示的知识库和文档")
    @GetMapping("/homepage")
    public Mono<ApiResponse<HomepageResponse>> getHomepage() {
        return documentService.getHomepageDocuments()
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取标签云", description = "获取文档标签云数据")
    @GetMapping("/tag-cloud")
    public Mono<ApiResponse<TagCloudResponse>> getTagCloud(
            @RequestParam(required = false) Long spaceId,
            @RequestParam(required = false) Long subSpaceId,
            @RequestParam(defaultValue = "50") Integer limit) {
        return documentService.getTagCloud(spaceId, subSpaceId, limit)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "知识搜索", description = "基于向量搜索的知识检索")
    @PostMapping("/search")
    public Mono<ApiResponse<KnowledgeSearchResponse>> searchKnowledge(
            @RequestBody KnowledgeSearchRequest request) {
        return documentService.searchKnowledge(request)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "文档问答", description = "基于指定文档进行智能问答")
    @PostMapping("/{id}/chat")
    public Mono<ApiResponse<ChatDocumentResponse>> chat(
            @PathVariable Long id,
            @RequestBody ChatDocumentRequest request) {
        return documentService.chatWithDocument(id, request)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "文档流式问答", description = "基于知识库进行流式问答（SSE）")
    @PostMapping(value = "/chat/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<ServerSentEvent<String>> chatStream(@RequestBody ChatStreamRequest request) {
        return documentService.chatWithDocumentsStream(request.getQuestion(), request.getSpaceId())
                .map(chunk -> ServerSentEvent.<String>builder()
                        .event("message")
                        .data(chunk)
                        .build())
                .concatWith(Flux.just(ServerSentEvent.<String>builder()
                        .event("done")
                        .data("[DONE]")
                        .build()))
                .delayElements(Duration.ofMillis(10));
    }
    
    @Operation(summary = "重试处理文档", description = "重新处理失败的文档（OCR和向量化）")
    @PostMapping("/retry-process")
    public Mono<ApiResponse<Document>> retryProcess(@RequestBody RetryProcessRequest request) {
        return documentService.retryProcessDocument(request.getDocumentId(), request.isForceRetry())
                .map(doc -> ApiResponse.success("重试处理已开始", doc));
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
    public static class RetryProcessRequest {
        private Long documentId;
        private boolean forceRetry;
    }
    
    @Data
    public static class ChatStreamRequest {
        private String question;
        private Long spaceId;
    }
}

