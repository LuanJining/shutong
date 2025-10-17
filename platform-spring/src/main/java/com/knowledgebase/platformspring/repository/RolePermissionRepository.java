package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.RolePermission;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface RolePermissionRepository extends R2dbcRepository<RolePermission, Long> {
    
    Flux<RolePermission> findByRoleId(Long roleId);
    
    Flux<RolePermission> findByPermissionId(Long permissionId);
    
    @Query("DELETE FROM role_permissions WHERE role_id = :roleId AND permission_id = :permissionId")
    Mono<Void> deleteByRoleIdAndPermissionId(Long roleId, Long permissionId);
}

