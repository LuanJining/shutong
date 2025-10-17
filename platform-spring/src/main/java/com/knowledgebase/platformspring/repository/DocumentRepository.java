package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Document;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface DocumentRepository extends R2dbcRepository<Document, Long> {
    
    Flux<Document> findBySpaceId(Long spaceId);
    
    Flux<Document> findBySubSpaceId(Long subSpaceId);
    
    Flux<Document> findByClassId(Long classId);
    
    Flux<Document> findByCreatedBy(Long createdBy);
    
    Flux<Document> findByStatus(String status);
    
    Flux<Document> findBySpaceIdAndStatus(Long spaceId, String status);
    
    Flux<Document> findByWorkflowId(Long workflowId);
    
    Mono<Long> countByStatus(String status);
}

