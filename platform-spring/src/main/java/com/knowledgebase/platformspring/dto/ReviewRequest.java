package com.knowledgebase.platformspring.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 文档审查请求
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ReviewRequest {
    
    /**
     * 临时文件路径（MinIO）
     */
    private String tempFilePath;
    
    /**
     * 文件名
     */
    private String fileName;
    
    /**
     * 文件类型
     */
    private String fileType;
    
    /**
     * 知识库范围（可选，限定检索范围）
     */
    private Long spaceId;
    
    /**
     * 审查选项
     */
    private ReviewOptions options;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ReviewOptions {
        /**
         * 是否检查格式
         */
        @Builder.Default
        private Boolean checkFormat = true;
        
        /**
         * 是否验证引用
         */
        @Builder.Default
        private Boolean verifyReferences = true;
        
        /**
         * 是否提供内容建议
         */
        @Builder.Default
        private Boolean suggestContent = true;
        
        /**
         * AI建议的详细程度 (1-3)
         */
        @Builder.Default
        private Integer detailLevel = 2;
    }
}

