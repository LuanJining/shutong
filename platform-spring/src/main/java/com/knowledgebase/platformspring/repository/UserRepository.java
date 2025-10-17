package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.User;

import reactor.core.publisher.Mono;

@Repository
public interface UserRepository extends R2dbcRepository<User, Long> {
    
    Mono<User> findByUsername(String username);
    
    Mono<User> findByPhone(String phone);
    
    Mono<User> findByEmail(String email);
    
    Mono<Boolean> existsByUsername(String username);
    
    Mono<Boolean> existsByPhone(String phone);
}

