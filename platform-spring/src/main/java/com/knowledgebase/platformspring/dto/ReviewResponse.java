package com.knowledgebase.platformspring.dto;

import java.time.LocalDateTime;
import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 文档审查响应
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ReviewResponse {
    
    /**
     * 会话ID
     */
    private String sessionId;
    
    /**
     * 文件名
     */
    private String fileName;
    
    /**
     * 审查状态
     */
    private String status;
    
    /**
     * 建议列表
     */
    private List<ReviewSuggestion> suggestions;
    
    /**
     * 总建议数
     */
    private Integer totalSuggestions;
    
    /**
     * 错误数量
     */
    private Integer errorCount;
    
    /**
     * 警告数量
     */
    private Integer warningCount;
    
    /**
     * 信息数量
     */
    private Integer infoCount;
    
    /**
     * 审查开始时间
     */
    private LocalDateTime startTime;
    
    /**
     * 审查完成时间
     */
    private LocalDateTime endTime;
    
    // 状态常量
    public static final String STATUS_PROCESSING = "PROCESSING";
    public static final String STATUS_COMPLETED = "COMPLETED";
    public static final String STATUS_FAILED = "FAILED";
}

