package com.knowledgebase.platformspring.dto;

import java.util.List;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class ChatDocumentRequest {
    
    @NotBlank(message = "Question is required")
    private String question;
    
    private List<Long> documentIds;
    private Integer limit;
}

