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
@Table("permissions")
public class Permission {
    
    @Id
    private Long id;
    
    @Column("name")
    private String name;
    
    @Column("display_name")
    private String displayName;
    
    @Column("description")
    private String description;
    
    @Column("resource")
    private String resource;
    
    @Column("action")
    private String action;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    @Column("deleted_at")
    private LocalDateTime deletedAt;
    
    // Permission constants
    public static final String VIEW_ALL_CONTENT = "view_all_content";
    public static final String CREATE_DOCUMENT = "create_document";
    public static final String DELETE_DOCUMENT = "delete_document";
    public static final String MOVE_DOCUMENT = "move_document";
    public static final String SET_DOCUMENT_PERMISSION = "set_document_permission";
    public static final String CREATE_SPACE = "create_space";
    public static final String MANAGE_SPACE_MEMBER = "manage_space_member";
    public static final String CONFIGURE_WORKFLOW = "configure_workflow";
    public static final String EXPORT_DATA = "export_data";
    public static final String EXPORT_ALL_DATA = "export_all_data";
    public static final String VIEW_OPERATION_LOG = "view_operation_log";
    public static final String ADD_DELETE_USER = "add_delete_user";
}

