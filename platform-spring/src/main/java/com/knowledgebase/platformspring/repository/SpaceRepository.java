package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Space;

import reactor.core.publisher.Flux;

@Repository
public interface SpaceRepository extends R2dbcRepository<Space, Long> {
    
    Flux<Space> findByCreatedBy(Long createdBy);
    
    Flux<Space> findByType(String type);
    
    Flux<Space> findByStatus(Integer status);
}

