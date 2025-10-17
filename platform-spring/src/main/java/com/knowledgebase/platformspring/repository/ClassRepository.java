package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import reactor.core.publisher.Flux;

@Repository
public interface ClassRepository extends R2dbcRepository<com.knowledgebase.platformspring.model.Class, Long> {
    
    Flux<com.knowledgebase.platformspring.model.Class> findBySubSpaceId(Long subSpaceId);
    
    Flux<com.knowledgebase.platformspring.model.Class> findByStatus(Integer status);
}

