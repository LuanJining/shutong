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
@Table("spaces")
public class Space {
    
    @Id
    private Long id;
    
    @Column("name")
    private String name;
    
    @Column("description")
    private String description;
    
    @Column("type")
    private String type; // department, project, team
    
    @Column("status")
    @Builder.Default
    private Integer status = 1;
    
    @Column("created_by")
    private Long createdBy;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    @Column("deleted_at")
    private LocalDateTime deletedAt;
    
    // Space type constants
    public static final String TYPE_DEPARTMENT = "department";
    public static final String TYPE_PROJECT = "project";
    public static final String TYPE_TEAM = "team";
}

