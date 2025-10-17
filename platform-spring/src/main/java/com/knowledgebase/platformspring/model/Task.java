package com.knowledgebase.platformspring.model;

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
@Table("tasks")
public class Task {
    
    @Id
    private Long id;
    
    @Column("workflow_id")
    private Long workflowId;
    
    @Column("step_id")
    private Long stepId;
    
    @Column("task_name")
    private String taskName;
    
    @Column("is_required")
    private Boolean isRequired;
    
    @Column("timeout_hours")
    private Integer timeoutHours;
    
    @Column("status")
    @Builder.Default
    private String status = STATUS_PROCESSING;
    
    @Column("approver_id")
    private Long approverId;
    
    @Column("approver_nick_name")
    private String approverNickName;
    
    @Column("comment")
    private String comment;
    
    // 注意：Go版本的Task没有CreatedAt和UpdatedAt字段
    
    // Task status constants
    public static final String STATUS_PROCESSING = "processing";
    public static final String STATUS_APPROVED = "approved";
    public static final String STATUS_APPROVED_BY_OTHER = "approved_by_others";
    public static final String STATUS_REJECTED = "rejected";
    public static final String STATUS_REJECTED_BY_OTHER = "rejected_by_others";
}

