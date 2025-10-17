package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.SpaceMember;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface SpaceMemberRepository extends R2dbcRepository<SpaceMember, Long> {
    
    Flux<SpaceMember> findBySpaceId(Long spaceId);
    
    Flux<SpaceMember> findByUserId(Long userId);
    
    Mono<SpaceMember> findBySpaceIdAndUserId(Long spaceId, Long userId);
    
    @Query("DELETE FROM space_members WHERE space_id = :spaceId AND user_id = :userId")
    Mono<Void> deleteBySpaceIdAndUserId(Long spaceId, Long userId);
    
    @Query("SELECT * FROM space_members WHERE space_id = :spaceId AND roles::jsonb @> :role::jsonb")
    Flux<SpaceMember> findBySpaceIdAndRole(Long spaceId, String role);
}

