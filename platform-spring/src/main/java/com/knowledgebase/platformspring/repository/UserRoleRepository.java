package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.UserRole;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface UserRoleRepository extends R2dbcRepository<UserRole, Long> {
    
    Flux<UserRole> findByUserId(Long userId);
    
    Flux<UserRole> findByRoleId(Long roleId);
    
    @Query("DELETE FROM user_roles WHERE user_id = :userId AND role_id = :roleId")
    Mono<Void> deleteByUserIdAndRoleId(Long userId, Long roleId);
}

