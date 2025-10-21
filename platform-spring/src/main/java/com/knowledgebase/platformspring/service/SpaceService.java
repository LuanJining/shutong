package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.dto.SpaceWithHierarchy;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.KnowledgeClass;
import com.knowledgebase.platformspring.model.Space;
import com.knowledgebase.platformspring.model.SubSpace;
import com.knowledgebase.platformspring.repository.ClassRepository;
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
    private final ClassRepository classRepository;
    
    public Flux<Space> getAllSpaces() {
        return spaceRepository.findAll();
    }
    
    public Flux<SpaceWithHierarchy> getAllSpacesWithHierarchy() {
        return spaceRepository.findAll()
                .flatMap(space -> 
                    subSpaceRepository.findBySpaceId(space.getId())
                        .flatMap(subSpace ->
                            classRepository.findBySubSpaceId(subSpace.getId())
                                .map(knowledgeClass -> SpaceWithHierarchy.ClassInfo.builder()
                                        .id(knowledgeClass.getId())
                                        .name(knowledgeClass.getName())
                                        .description(knowledgeClass.getDescription())
                                        .subSpaceId(knowledgeClass.getSubSpaceId())
                                        .status(knowledgeClass.getStatus())
                                        .createdBy(knowledgeClass.getCreatedBy())
                                        .createdAt(knowledgeClass.getCreatedAt())
                                        .updatedAt(knowledgeClass.getUpdatedAt())
                                        .deletedAt(knowledgeClass.getDeletedAt())
                                        .build())
                                .collectList()
                                .map(classes -> SpaceWithHierarchy.SubSpaceWithClasses.builder()
                                        .id(subSpace.getId())
                                        .name(subSpace.getName())
                                        .description(subSpace.getDescription())
                                        .spaceId(subSpace.getSpaceId())
                                        .status(subSpace.getStatus())
                                        .createdBy(subSpace.getCreatedBy())
                                        .createdAt(subSpace.getCreatedAt())
                                        .updatedAt(subSpace.getUpdatedAt())
                                        .deletedAt(subSpace.getDeletedAt())
                                        .classes(classes)
                                        .build())
                        )
                        .collectList()
                        .map(subSpaces -> SpaceWithHierarchy.builder()
                                .id(space.getId())
                                .name(space.getName())
                                .description(space.getDescription())
                                .type(space.getType())
                                .status(space.getStatus())
                                .createdBy(space.getCreatedBy())
                                .createdAt(space.getCreatedAt())
                                .updatedAt(space.getUpdatedAt())
                                .deletedAt(space.getDeletedAt())
                                .subSpaces(subSpaces)
                                .build())
                );
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
    
    public Mono<KnowledgeClass> createKnowledgeClass(KnowledgeClass knowledgeClass, Long userId) {
        knowledgeClass.setCreatedBy(userId);
        knowledgeClass.setCreatedAt(LocalDateTime.now());
        knowledgeClass.setUpdatedAt(LocalDateTime.now());
        knowledgeClass.setStatus(1);
        
        return classRepository.save(knowledgeClass);
    }
}

