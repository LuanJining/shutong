package com.knowledgebase.platformspring.util;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.knowledgebase.platformspring.dto.DocumentSection;

import lombok.extern.slf4j.Slf4j;

/**
 * 文档解析工具
 */
@Slf4j
public class DocumentParser {
    
    // 章节标题正则
    private static final Pattern CHAPTER_PATTERN = Pattern.compile("^第[一二三四五六七八九十百]+章\\s+.+$");
    private static final Pattern SECTION_PATTERN = Pattern.compile("^第[一二三四五六七八九十百]+节\\s+.+$");
    private static final Pattern ARTICLE_PATTERN = Pattern.compile("^第[一二三四五六七八九十百0-9]+条\\s+.+$");
    
    // 日期格式正则
    private static final Pattern DATE_PATTERN = Pattern.compile("\\d{4}年\\d{1,2}月\\d{1,2}日");
    
    /**
     * 解析文档为段落
     */
    public static List<DocumentSection> parse(String content) {
        if (content == null || content.trim().isEmpty()) {
            return new ArrayList<>();
        }
        
        List<DocumentSection> sections = new ArrayList<>();
        String[] lines = content.split("\n");
        
        int position = 0;
        int lineNumber = 0;
        
        for (String line : lines) {
            lineNumber++;
            String trimmedLine = line.trim();
            
            if (trimmedLine.isEmpty()) {
                continue;
            }
            
            DocumentSection section = DocumentSection.builder()
                    .position(position++)
                    .content(trimmedLine)
                    .startLine(lineNumber)
                    .endLine(lineNumber)
                    .build();
            
            // 识别段落类型
            if (CHAPTER_PATTERN.matcher(trimmedLine).matches()) {
                section.setType(DocumentSection.TYPE_CHAPTER);
                section.setLevel(1);
            } else if (SECTION_PATTERN.matcher(trimmedLine).matches()) {
                section.setType(DocumentSection.TYPE_SECTION);
                section.setLevel(2);
            } else if (ARTICLE_PATTERN.matcher(trimmedLine).matches()) {
                section.setType(DocumentSection.TYPE_ARTICLE);
                section.setLevel(3);
            } else if (DATE_PATTERN.matcher(trimmedLine).find()) {
                section.setType(DocumentSection.TYPE_DATE);
            } else {
                section.setType(DocumentSection.TYPE_PARAGRAPH);
            }
            
            sections.add(section);
        }
        
        log.debug("Parsed document into {} sections", sections.size());
        return sections;
    }
    
    /**
     * 提取法规引用
     */
    public static List<String> extractReferences(String content) {
        List<String> references = new ArrayList<>();
        
        // 匹配《XX法》《XX条例》等
        Pattern refPattern = Pattern.compile("《[^》]+[法条例规定办法准则]》(?:第[\\d一二三四五六七八九十百]+条)?");
        Matcher matcher = refPattern.matcher(content);
        
        while (matcher.find()) {
            references.add(matcher.group());
        }
        
        return references;
    }
    
    /**
     * 合并相邻段落（用于上下文分析）
     */
    public static String mergeContext(List<DocumentSection> sections, int centerIndex, int contextWindow) {
        StringBuilder context = new StringBuilder();
        
        int start = Math.max(0, centerIndex - contextWindow);
        int end = Math.min(sections.size(), centerIndex + contextWindow + 1);
        
        for (int i = start; i < end; i++) {
            if (i > start) {
                context.append("\n");
            }
            context.append(sections.get(i).getContent());
        }
        
        return context.toString();
    }
}

