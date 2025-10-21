package com.knowledgebase.platformspring.dto;

import java.time.LocalDateTime;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SpaceWithHierarchy {
    
    private Long id;
    private String name;
    private String description;
    private String type;
    private Integer status;
    
    @JsonProperty("created_by")
    private Long createdBy;
    
    @JsonProperty("created_at")
    private LocalDateTime createdAt;
    
    @JsonProperty("updated_at")
    private LocalDateTime updatedAt;
    
    @JsonProperty("deleted_at")
    private LocalDateTime deletedAt;
    
    @JsonProperty("sub_spaces")
    private List<SubSpaceWithClasses> subSpaces;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class SubSpaceWithClasses {
        private Long id;
        private String name;
        private String description;
        
        @JsonProperty("space_id")
        private Long spaceId;
        
        private Integer status;
        
        @JsonProperty("created_by")
        private Long createdBy;
        
        @JsonProperty("created_at")
        private LocalDateTime createdAt;
        
        @JsonProperty("updated_at")
        private LocalDateTime updatedAt;
        
        @JsonProperty("deleted_at")
        private LocalDateTime deletedAt;
        
        private List<ClassInfo> classes;
    }
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ClassInfo {
        private Long id;
        private String name;
        private String description;
        
        @JsonProperty("sub_space_id")
        private Long subSpaceId;
        
        private Integer status;
        
        @JsonProperty("created_by")
        private Long createdBy;
        
        @JsonProperty("created_at")
        private LocalDateTime createdAt;
        
        @JsonProperty("updated_at")
        private LocalDateTime updatedAt;
        
        @JsonProperty("deleted_at")
        private LocalDateTime deletedAt;
    }
}

