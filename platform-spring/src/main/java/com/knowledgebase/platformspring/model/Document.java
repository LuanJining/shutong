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
@Table("documents")
public class Document {
    
    @Id
    private Long id;
    
    @Column("title")
    private String title;
    
    @Column("content")
    private String content;
    
    @Column("file_name")
    private String fileName;
    
    @Column("file_path")
    private String filePath;
    
    @Column("file_size")
    private Long fileSize;
    
    @Column("file_type")
    private String fileType;
    
    @Column("status")
    @Builder.Default
    private String status = STATUS_UPLOADING;
    
    @Column("need_approval")
    @Builder.Default
    private Boolean needApproval = false;
    
    @Column("mime_type")
    private String mimeType;
    
    @Column("version")
    @Builder.Default
    private String version = "v1.0.0";
    
    @Column("use_type")
    @Builder.Default
    private String useType = USE_TYPE_VIEWABLE;
    
    @Column("space_id")
    private Long spaceId;
    
    @Column("sub_space_id")
    private Long subSpaceId;
    
    @Column("class_id")
    private Long classId;
    
    @Column("created_by")
    private Long createdBy;
    
    @Column("creator_nick_name")
    private String creatorNickName;
    
    @Column("department")
    private String department;
    
    @Column("workflow_id")
    private Long workflowId;
    
    @Column("tags")
    private String tags;
    
    @Column("summary")
    private String summary;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    @Column("deleted_at")
    private LocalDateTime deletedAt;
    
    @Column("parse_error")
    private String parseError;
    
    @Column("processed_at")
    private LocalDateTime processedAt;
    
    @Column("vector_count")
    @Builder.Default
    private Integer vectorCount = 0;
    
    @Column("process_progress")
    @Builder.Default
    private Integer processProgress = 0;
    
    @Column("retry_count")
    @Builder.Default
    private Integer retryCount = 0;
    
    @Column("last_retry_at")
    private LocalDateTime lastRetryAt;
    
    // Document status constants
    public static final String STATUS_UPLOADING = "uploading";
    public static final String STATUS_PROCESSING = "processing";
    public static final String STATUS_VECTORIZING = "vectorizing";
    public static final String STATUS_PENDING_APPROVAL = "pending_approval";
    public static final String STATUS_PENDING_PUBLISH = "pending_publish";
    public static final String STATUS_PUBLISHED = "published";
    public static final String STATUS_FAILED = "failed";
    public static final String STATUS_PROCESS_FAILED = "process_failed";
    
    // Use type constants
    public static final String USE_TYPE_VIEWABLE = "viewable";
    public static final String USE_TYPE_APPLICABLE = "applicable";
}

