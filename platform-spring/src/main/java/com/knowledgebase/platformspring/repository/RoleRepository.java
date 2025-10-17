package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.Role;

import reactor.core.publisher.Mono;

@Repository
public interface RoleRepository extends R2dbcRepository<Role, Long> {
    
    Mono<Role> findByName(String name);
}

