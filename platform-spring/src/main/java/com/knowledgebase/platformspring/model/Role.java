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
@Table("roles")
public class Role {
    
    @Id
    private Long id;
    
    @Column("name")
    private String name;
    
    @Column("display_name")
    private String displayName;
    
    @Column("description")
    private String description;
    
    @Column("status")
    @Builder.Default
    private Integer status = 1;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    @Column("deleted_at")
    private LocalDateTime deletedAt;
    
    // Role names constants
    public static final String SUPER_ADMIN = "super_admin";
    public static final String CORP_ADMIN = "corp_admin";
}

