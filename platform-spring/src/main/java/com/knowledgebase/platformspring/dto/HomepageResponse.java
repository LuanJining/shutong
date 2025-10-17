package com.knowledgebase.platformspring.dto;

import java.time.LocalDateTime;
import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class HomepageResponse {
    private List<HomepageSpace> spaces;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class HomepageSpace {
        private Long id;
        private String name;
        private String description;
        private List<HomepageSubSpace> subSpaces;
    }
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class HomepageSubSpace {
        private Long id;
        private String name;
        private String description;
        private List<HomepageDocument> documents;
    }
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class HomepageDocument {
        private Long id;
        private String title;
        private String fileName;
        private Long fileSize;
        private String fileType;
        private String status;
        private String creatorNickName;
        private String summary;
        private LocalDateTime createdAt;
        private LocalDateTime updatedAt;
    }
}

