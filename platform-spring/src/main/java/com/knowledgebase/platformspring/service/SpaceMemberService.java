package com.knowledgebase.platformspring.service;

import java.util.List;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.SpaceMember;
import com.knowledgebase.platformspring.repository.SpaceMemberRepository;
import com.knowledgebase.platformspring.repository.SpaceRepository;
import com.knowledgebase.platformspring.repository.UserRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Slf4j
@Service
@RequiredArgsConstructor
public class SpaceMemberService {
    
    private final SpaceMemberRepository spaceMemberRepository;
    private final SpaceRepository spaceRepository;
    private final UserRepository userRepository;
    
    public Flux<SpaceMember> getSpaceMembers(Long spaceId) {
        return spaceMemberRepository.findBySpaceId(spaceId);
    }
    
    public Mono<SpaceMember> addSpaceMember(Long spaceId, Long userId, List<String> roles) {
        log.debug("Adding space member: spaceId={}, userId={}, roles={}", spaceId, userId, roles);
        
        // 验证空间存在
        return spaceRepository.findById(spaceId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("空间不存在")))
                .then(userRepository.findById(userId))
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不存在")))
                .then(spaceMemberRepository.findBySpaceIdAndUserId(spaceId, userId)
                        .flatMap(existing -> {
                            // 更新已存在的成员
                            log.debug("Updating existing space member: {}", existing.getId());
                            existing.setRoles(roles);
                            return spaceMemberRepository.save(existing);
                        })
                        .switchIfEmpty(Mono.defer(() -> {
                            // 创建新成员
                            log.debug("Creating new space member");
                            SpaceMember newMember = SpaceMember.builder()
                                    .spaceId(spaceId)
                                    .userId(userId)
                                    .build();
                            newMember.setRoles(roles);
                            return spaceMemberRepository.save(newMember);
                        }))
                )
                .doOnSuccess(member -> log.debug("Space member saved successfully: {}", member.getId()))
                .doOnError(error -> log.error("Failed to add space member: {}", error.getMessage()));
    }
    
    public Mono<Void> removeSpaceMember(Long spaceId, Long userId) {
        return spaceMemberRepository.findBySpaceIdAndUserId(spaceId, userId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不在该空间中")))
                .flatMap(spaceMemberRepository::delete);
    }
    
    public Mono<SpaceMember> updateSpaceMemberRoles(Long spaceId, Long userId, List<String> roles) {
        return spaceMemberRepository.findBySpaceIdAndUserId(spaceId, userId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不在该空间中")))
                .flatMap(member -> {
                    member.setRoles(roles);
                    return spaceMemberRepository.save(member);
                });
    }
    
    public Flux<SpaceMember> getUserRoles(Long userId) {
        return spaceMemberRepository.findByUserId(userId);
    }
}

