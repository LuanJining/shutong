package com.knowledgebase.platformspring.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class KnowledgeSearchRequest {
    
    @NotBlank(message = "Query is required")
    private String query;
    
    private Integer limit;
    private Long spaceId;
    private Long subSpaceId;
    private Long classId;
}

