package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Permission;

import reactor.core.publisher.Mono;

@Repository
public interface PermissionRepository extends R2dbcRepository<Permission, Long> {
    
    Mono<Permission> findByName(String name);
}

