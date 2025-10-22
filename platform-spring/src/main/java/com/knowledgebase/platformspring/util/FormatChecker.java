package com.knowledgebase.platformspring.util;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.knowledgebase.platformspring.dto.DocumentSection;
import com.knowledgebase.platformspring.dto.ReviewSuggestion;

import lombok.extern.slf4j.Slf4j;

/**
 * 格式规范检查工具
 */
@Slf4j
public class FormatChecker {
    
    /**
     * 检查段落格式
     */
    public static List<ReviewSuggestion> check(DocumentSection section) {
        List<ReviewSuggestion> suggestions = new ArrayList<>();
        
        String content = section.getContent();
        
        // 1. 检查标点符号（全角/半角）
        suggestions.addAll(checkPunctuation(section, content));
        
        // 2. 检查章节编号格式
        suggestions.addAll(checkNumbering(section, content));
        
        // 3. 检查日期格式
        suggestions.addAll(checkDateFormat(section, content));
        
        // 4. 检查空格规范
        suggestions.addAll(checkSpacing(section, content));
        
        return suggestions;
    }
    
    /**
     * 检查标点符号
     */
    private static List<ReviewSuggestion> checkPunctuation(DocumentSection section, String content) {
        List<ReviewSuggestion> suggestions = new ArrayList<>();
        
        // 检查半角标点
        if (content.matches(".*[,.].*")) {
            String corrected = content.replace(",", "，").replace(".", "。");
            suggestions.add(ReviewSuggestion.builder()
                    .type(ReviewSuggestion.TYPE_PUNCTUATION)
                    .severity(ReviewSuggestion.SEVERITY_ERROR)
                    .position(section.getPosition())
                    .originalText(content)
                    .suggestedText(corrected)
                    .reason("公文应使用全角标点符号")
                    .build());
        }
        
        // 检查括号
        if (content.matches(".*[()].*")) {
            String corrected = content.replace("(", "（").replace(")", "）");
            suggestions.add(ReviewSuggestion.builder()
                    .type(ReviewSuggestion.TYPE_PUNCTUATION)
                    .severity(ReviewSuggestion.SEVERITY_WARNING)
                    .position(section.getPosition())
                    .originalText(content)
                    .suggestedText(corrected)
                    .reason("建议使用全角括号")
                    .build());
        }
        
        return suggestions;
    }
    
    /**
     * 检查编号格式
     */
    private static List<ReviewSuggestion> checkNumbering(DocumentSection section, String content) {
        List<ReviewSuggestion> suggestions = new ArrayList<>();
        
        // 检查章节编号后是否有空格
        Pattern chapterPattern = Pattern.compile("(第[一二三四五六七八九十百]+章)([^\\s])");
        Matcher matcher = chapterPattern.matcher(content);
        
        if (matcher.find()) {
            String corrected = content.replace(matcher.group(0), matcher.group(1) + " " + matcher.group(2));
            suggestions.add(ReviewSuggestion.builder()
                    .type(ReviewSuggestion.TYPE_NUMBERING_ERROR)
                    .severity(ReviewSuggestion.SEVERITY_WARNING)
                    .position(section.getPosition())
                    .originalText(content)
                    .suggestedText(corrected)
                    .reason("章节编号后应有空格")
                    .build());
        }
        
        return suggestions;
    }
    
    /**
     * 检查日期格式
     */
    private static List<ReviewSuggestion> checkDateFormat(DocumentSection section, String content) {
        List<ReviewSuggestion> suggestions = new ArrayList<>();
        
        // 检查是否使用了 YYYY-MM-DD 格式（应该用中文格式）
        Pattern datePattern = Pattern.compile("(\\d{4})-(\\d{1,2})-(\\d{1,2})");
        Matcher matcher = datePattern.matcher(content);
        
        if (matcher.find()) {
            String year = matcher.group(1);
            String month = matcher.group(2);
            String day = matcher.group(3);
            String corrected = content.replace(
                matcher.group(0), 
                year + "年" + month + "月" + day + "日"
            );
            
            suggestions.add(ReviewSuggestion.builder()
                    .type(ReviewSuggestion.TYPE_DATE_FORMAT)
                    .severity(ReviewSuggestion.SEVERITY_WARNING)
                    .position(section.getPosition())
                    .originalText(content)
                    .suggestedText(corrected)
                    .reason("公文日期应使用中文格式（YYYY年MM月DD日）")
                    .build());
        }
        
        return suggestions;
    }
    
    /**
     * 检查空格规范
     */
    private static List<ReviewSuggestion> checkSpacing(DocumentSection section, String content) {
        List<ReviewSuggestion> suggestions = new ArrayList<>();
        
        // 检查中文和英文/数字之间是否有空格
        Pattern pattern = Pattern.compile("([\\u4e00-\\u9fa5])([a-zA-Z0-9])|([a-zA-Z0-9])([\\u4e00-\\u9fa5])");
        Matcher matcher = pattern.matcher(content);
        
        if (matcher.find()) {
            suggestions.add(ReviewSuggestion.builder()
                    .type(ReviewSuggestion.TYPE_FORMAT_ERROR)
                    .severity(ReviewSuggestion.SEVERITY_INFO)
                    .position(section.getPosition())
                    .originalText(content)
                    .suggestedText(null)  // 具体修正较复杂，仅提示
                    .reason("建议在中文和英文/数字之间添加空格")
                    .build());
        }
        
        return suggestions;
    }
}

