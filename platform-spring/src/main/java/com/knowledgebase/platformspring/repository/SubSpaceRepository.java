package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.SubSpace;

import reactor.core.publisher.Flux;

@Repository
public interface SubSpaceRepository extends R2dbcRepository<SubSpace, Long> {
    
    Flux<SubSpace> findBySpaceId(Long spaceId);
    
    Flux<SubSpace> findByStatus(Integer status);
}

