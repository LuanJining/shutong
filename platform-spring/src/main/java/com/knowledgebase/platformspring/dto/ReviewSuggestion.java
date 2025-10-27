package com.knowledgebase.platformspring.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 审查建议
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ReviewSuggestion {

    /**
     * 建议ID（用于标识和接受建议）
     */
    private String id;

    /**
     * 建议类型
     */
    private String type;

    /**
     * 严重程度: ERROR, WARNING, INFO
     */
    private String severity;

    /**
     * 问题位置（段落索引或行号）
     */
    private Integer position;

    /**
     * 原始文本
     */
    private String originalText;

    /**
     * 建议修改后的文本
     */
    private String suggestedText;

    /**
     * 建议原因/说明
     */
    private String reason;

    /**
     * 知识来源（如果是基于知识库的建议）
     */
    private String knowledgeSource;

    /**
     * 知识来源ID（用于溯源）
     */
    private Long knowledgeDocumentId;

    /**
     * 文档内容（仅在第一条建议中返回，用于前端显示）
     */
    private String documentContent;

    // 建议类型常量
    public static final String TYPE_FORMAT_ERROR = "FORMAT_ERROR";
    public static final String TYPE_PUNCTUATION = "PUNCTUATION";
    public static final String TYPE_REFERENCE_OUTDATED = "REFERENCE_OUTDATED";
    public static final String TYPE_REFERENCE_MISSING = "REFERENCE_MISSING";
    public static final String TYPE_CONTENT_ENHANCEMENT = "CONTENT_ENHANCEMENT";
    public static final String TYPE_NUMBERING_ERROR = "NUMBERING_ERROR";
    public static final String TYPE_DATE_FORMAT = "DATE_FORMAT";

    // 严重程度常量
    public static final String SEVERITY_ERROR = "ERROR";
    public static final String SEVERITY_WARNING = "WARNING";
    public static final String SEVERITY_INFO = "INFO";
}

