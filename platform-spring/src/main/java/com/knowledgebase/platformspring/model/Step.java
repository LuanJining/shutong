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
@Table("steps")
public class Step {
    
    @Id
    private Long id;
    
    @Column("workflow_id")
    private Long workflowId;
    
    @Column("step_name")
    private String stepName;
    
    @Column("step_order")
    private Integer stepOrder;
    
    @Column("step_role")
    private String stepRole;
    
    @Column("is_required")
    private Boolean isRequired;
    
    @Column("timeout_hours")
    private Integer timeoutHours;
    
    @Column("status")
    @Builder.Default
    private String status = STATUS_PROCESSING;
    
    // 注意：Go版本的Step没有CreatedAt和UpdatedAt字段
    
    // Step status constants
    public static final String STATUS_PROCESSING = "processing";
    public static final String STATUS_APPROVED = "approved";
    public static final String STATUS_REJECTED = "rejected";
}

