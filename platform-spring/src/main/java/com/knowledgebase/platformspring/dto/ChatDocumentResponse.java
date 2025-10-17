package com.knowledgebase.platformspring.dto;

import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ChatDocumentResponse {
    
    private String answer;
    private List<ChatDocumentSource> sources;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ChatDocumentSource {
        private Long documentId;
        private String title;
        private String filePath;
    }
}

