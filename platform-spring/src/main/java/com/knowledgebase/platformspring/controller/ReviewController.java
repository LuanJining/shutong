package com.knowledgebase.platformspring.controller;

import java.time.LocalDateTime;

import org.springframework.http.MediaType;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.http.codec.multipart.FilePart;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.dto.ReviewRequest;
import com.knowledgebase.platformspring.dto.ReviewResponse;
import com.knowledgebase.platformspring.dto.ReviewSuggestion;
import com.knowledgebase.platformspring.service.ReviewService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

/**
 * 智能审查控制器
 */
@Slf4j
@RestController
@RequestMapping("/api/v1/review")
@RequiredArgsConstructor
@Tag(name = "智能审查", description = "文档智能审查相关接口")
public class ReviewController {
    
    private final ReviewService reviewService;
    
    /**
     * 上传文档用于审查
     */
    @PostMapping(value = "/upload", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @Operation(summary = "上传文档", description = "上传文档到临时存储，准备进行审查")
    public Mono<ApiResponse<String>> uploadDocument(
            @RequestPart("file") FilePart file) {
        
        log.info("Upload document for review: {}", file.filename());
        
        return reviewService.uploadForReview(file, file.filename())
                .map(sessionId -> ApiResponse.success(sessionId))
                .onErrorResume(e -> {
                    log.error("Upload failed", e);
                    return Mono.just(ApiResponse.error(e.getMessage()));
                });
    }
    
    /**
     * 开始审查文档（SSE流式返回）
     */
    @GetMapping(value = "/{sessionId}/suggestions", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    @Operation(summary = "获取审查建议", description = "实时流式返回文档审查建议")
    public Flux<ServerSentEvent<ReviewSuggestion>> getReviewSuggestions(
            @PathVariable String sessionId,
            @RequestParam String fileName,
            @RequestParam String fileType,
            @RequestParam(required = false) Long spaceId,
            @RequestParam(required = false, defaultValue = "true") Boolean checkFormat,
            @RequestParam(required = false, defaultValue = "true") Boolean verifyReferences,
            @RequestParam(required = false, defaultValue = "true") Boolean suggestContent) {
        
        log.info("Start review: sessionId={}, fileName={}", sessionId, fileName);
        
        // 构造请求
        ReviewRequest request = ReviewRequest.builder()
                .tempFilePath(reviewService.getTempFilePath(sessionId, fileName))
                .fileName(fileName)
                .fileType(fileType)
                .spaceId(spaceId)
                .options(ReviewRequest.ReviewOptions.builder()
                        .checkFormat(checkFormat)
                        .verifyReferences(verifyReferences)
                        .suggestContent(suggestContent)
                        .build())
                .build();
        
        return reviewService.reviewDocument(request)
                .map(suggestion -> ServerSentEvent.<ReviewSuggestion>builder()
                        .data(suggestion)
                        .build())
                .concatWith(Mono.just(ServerSentEvent.<ReviewSuggestion>builder()
                        .comment("DONE")
                        .build()))
                .doOnComplete(() -> log.info("Review completed: sessionId={}", sessionId))
                .doOnError(e -> log.error("Review failed: sessionId={}, error={}", sessionId, e.getMessage()));
    }
    
    /**
     * 获取审查结果摘要
     */
    @GetMapping("/{sessionId}/summary")
    @Operation(summary = "获取审查摘要", description = "获取文档审查结果的统计信息")
    public Mono<ApiResponse<ReviewResponse>> getReviewSummary(
            @PathVariable String sessionId,
            @RequestParam String fileName,
            @RequestParam String fileType,
            @RequestParam(required = false) Long spaceId) {
        
        log.info("Get review summary: sessionId={}", sessionId);
        
        ReviewRequest request = ReviewRequest.builder()
                .tempFilePath(reviewService.getTempFilePath(sessionId, fileName))
                .fileName(fileName)
                .fileType(fileType)
                .spaceId(spaceId)
                .build();
        
        LocalDateTime startTime = LocalDateTime.now();
        
        return reviewService.reviewDocument(request)
                .collectList()
                .map(suggestions -> {
                    // 统计各类型建议数量
                    int errorCount = 0;
                    int warningCount = 0;
                    int infoCount = 0;
                    
                    for (ReviewSuggestion suggestion : suggestions) {
                        switch (suggestion.getSeverity()) {
                            case ReviewSuggestion.SEVERITY_ERROR -> errorCount++;
                            case ReviewSuggestion.SEVERITY_WARNING -> warningCount++;
                            case ReviewSuggestion.SEVERITY_INFO -> infoCount++;
                        }
                    }
                    
                    ReviewResponse response = ReviewResponse.builder()
                            .sessionId(sessionId)
                            .fileName(fileName)
                            .status(ReviewResponse.STATUS_COMPLETED)
                            .suggestions(suggestions)
                            .totalSuggestions(suggestions.size())
                            .errorCount(errorCount)
                            .warningCount(warningCount)
                            .infoCount(infoCount)
                            .startTime(startTime)
                            .endTime(LocalDateTime.now())
                            .build();
                    
                    return ApiResponse.success(response);
                })
                .onErrorResume(e -> {
                    log.error("Get summary failed", e);
                    return Mono.just(ApiResponse.error(e.getMessage()));
                });
    }
}

