package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Task;

import reactor.core.publisher.Flux;

@Repository
public interface TaskRepository extends R2dbcRepository<Task, Long> {
    
    Flux<Task> findByWorkflowId(Long workflowId);
    
    Flux<Task> findByStepId(Long stepId);
    
    Flux<Task> findByApproverId(Long approverId);
    
    Flux<Task> findByStatus(String status);
    
    Flux<Task> findByApproverIdAndStatus(Long approverId, String status);
}

