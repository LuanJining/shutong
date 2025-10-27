package com.knowledgebase.platformspring.dto;

import java.util.List;

import jakarta.validation.constraints.NotBlank;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 接受建议请求
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class AcceptSuggestionRequest {

    /**
     * 会话ID
     */
    @NotBlank(message = "会话ID不能为空")
    private String sessionId;

    /**
     * 文件名
     */
    @NotBlank(message = "文件名不能为空")
    private String fileName;

    /**
     * 文件类型
     */
    @NotBlank(message = "文件类型不能为空")
    private String fileType;

    /**
     * 要接受的建议ID列表
     */
    private List<String> acceptedSuggestionIds;

    /**
     * 是否应用所有建议
     */
    @Builder.Default
    private Boolean applyAll = false;
}
