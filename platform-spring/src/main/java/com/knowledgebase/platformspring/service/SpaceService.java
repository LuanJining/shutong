package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.Space;
import com.knowledgebase.platformspring.model.SubSpace;
import com.knowledgebase.platformspring.repository.SpaceRepository;
import com.knowledgebase.platformspring.repository.SubSpaceRepository;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
public class SpaceService {
    
    private final SpaceRepository spaceRepository;
    private final SubSpaceRepository subSpaceRepository;
    
    public Flux<Space> getAllSpaces() {
        return spaceRepository.findAll();
    }
    
    public Mono<Space> getSpaceById(Long id) {
        return spaceRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Space not found")));
    }
    
    public Mono<Space> createSpace(Space space, Long userId) {
        space.setCreatedBy(userId);
        space.setCreatedAt(LocalDateTime.now());
        space.setUpdatedAt(LocalDateTime.now());
        space.setStatus(1);
        
        return spaceRepository.save(space);
    }
    
    public Mono<Space> updateSpace(Long id, Space space) {
        return spaceRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Space not found")))
                .flatMap(existing -> {
                    existing.setName(space.getName());
                    existing.setDescription(space.getDescription());
                    existing.setType(space.getType());
                    existing.setUpdatedAt(LocalDateTime.now());
                    return spaceRepository.save(existing);
                });
    }
    
    public Mono<Void> deleteSpace(Long id) {
        return spaceRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Space not found")))
                .flatMap(spaceRepository::delete);
    }
    
    public Flux<SubSpace> getSubSpacesBySpaceId(Long spaceId) {
        return subSpaceRepository.findBySpaceId(spaceId);
    }
    
    public Mono<SubSpace> createSubSpace(SubSpace subSpace, Long userId) {
        subSpace.setCreatedBy(userId);
        subSpace.setCreatedAt(LocalDateTime.now());
        subSpace.setUpdatedAt(LocalDateTime.now());
        subSpace.setStatus(1);
        
        return subSpaceRepository.save(subSpace);
    }
}

