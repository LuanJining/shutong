package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.KnowledgeClass;

import reactor.core.publisher.Flux;

@Repository
public interface ClassRepository extends R2dbcRepository<KnowledgeClass, Long> {
    
    Flux<KnowledgeClass> findBySubSpaceId(Long subSpaceId);
    
    Flux<KnowledgeClass> findByStatus(Integer status);
}

