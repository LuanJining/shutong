package com.knowledgebase.platformspring.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 文档段落
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DocumentSection {
    
    /**
     * 段落位置（索引）
     */
    private Integer position;
    
    /**
     * 段落类型
     */
    private String type;
    
    /**
     * 段落内容
     */
    private String content;
    
    /**
     * 起始行号
     */
    private Integer startLine;
    
    /**
     * 结束行号
     */
    private Integer endLine;
    
    /**
     * 层级（如章节层级）
     */
    private Integer level;
    
    // 段落类型常量
    public static final String TYPE_TITLE = "TITLE";          // 标题
    public static final String TYPE_HEADER = "HEADER";        // 红头
    public static final String TYPE_CHAPTER = "CHAPTER";      // 章
    public static final String TYPE_SECTION = "SECTION";      // 节
    public static final String TYPE_ARTICLE = "ARTICLE";      // 条
    public static final String TYPE_PARAGRAPH = "PARAGRAPH";  // 段落
    public static final String TYPE_SIGNATURE = "SIGNATURE";  // 落款
    public static final String TYPE_DATE = "DATE";           // 日期
}

