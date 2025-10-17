package com.knowledgebase.platformspring.model;

import java.time.LocalDateTime;

import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("workflows")
public class Workflow {
    
    @Id
    private Long id;
    
    @Column("name")
    private String name;
    
    @Column("description")
    private String description;
    
    @Column("space_id")
    private Long spaceId;
    
    @Column("status")
    @Builder.Default
    private String status = STATUS_PROCESSING;
    
    @Column("current_step_id")
    private Long currentStepId;
    
    @Column("resource_type")
    private String resourceType;
    
    @Column("resource_id")
    private Long resourceId;
    
    @Column("created_by")
    private Long createdBy;
    
    @Column("creator_nick_name")
    private String creatorNickName;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    // Workflow status constants
    public static final String STATUS_PROCESSING = "processing";
    public static final String STATUS_COMPLETED = "completed";
    public static final String STATUS_CANCELLED = "cancelled";
}

