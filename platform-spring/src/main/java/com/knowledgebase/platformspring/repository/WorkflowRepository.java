package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Workflow;

import reactor.core.publisher.Flux;

@Repository
public interface WorkflowRepository extends R2dbcRepository<Workflow, Long> {
    
    Flux<Workflow> findBySpaceId(Long spaceId);
    
    Flux<Workflow> findByStatus(String status);
    
    Flux<Workflow> findByResourceTypeAndResourceId(String resourceType, Long resourceId);
    
    Flux<Workflow> findByCreatedBy(Long createdBy);
}

