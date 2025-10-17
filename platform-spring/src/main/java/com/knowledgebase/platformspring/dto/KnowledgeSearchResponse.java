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
public class KnowledgeSearchResponse {
    
    private List<KnowledgeSearchResult> items;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class KnowledgeSearchResult {
        private Long documentId;
        private Long chunkId;
        private String title;
        private String content;
        private String snippet;
        private Double score;
        private String fileName;
    }
}

