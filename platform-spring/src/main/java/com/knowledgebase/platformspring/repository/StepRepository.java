package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Step;

import reactor.core.publisher.Flux;

@Repository
public interface StepRepository extends R2dbcRepository<Step, Long> {
    
    Flux<Step> findByWorkflowId(Long workflowId);
    
    Flux<Step> findByWorkflowIdOrderByStepOrderAsc(Long workflowId);
    
    Flux<Step> findByStatus(String status);
}

