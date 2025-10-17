package com.knowledgebase.platformspring.controller;

import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.model.Space;
import com.knowledgebase.platformspring.model.SubSpace;
import com.knowledgebase.platformspring.service.SpaceService;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v1/spaces")
@RequiredArgsConstructor
public class SpaceController {
    
    private final SpaceService spaceService;
    
    @GetMapping
    public Flux<Space> getAllSpaces() {
        return spaceService.getAllSpaces();
    }
    
    @GetMapping("/{id}")
    public Mono<ApiResponse<Space>> getSpaceById(@PathVariable Long id) {
        return spaceService.getSpaceById(id)
                .map(ApiResponse::success);
    }
    
    @PostMapping
    public Mono<ApiResponse<Space>> createSpace(@Valid @RequestBody Space space, 
                                                 Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return spaceService.createSpace(space, userId)
                .map(created -> ApiResponse.success("Space created successfully", created));
    }
    
    @PutMapping("/{id}")
    public Mono<ApiResponse<Space>> updateSpace(@PathVariable Long id, 
                                                 @Valid @RequestBody Space space) {
        return spaceService.updateSpace(id, space)
                .map(updated -> ApiResponse.success("Space updated successfully", updated));
    }
    
    @DeleteMapping("/{id}")
    public Mono<ApiResponse<Void>> deleteSpace(@PathVariable Long id) {
        return spaceService.deleteSpace(id)
                .then(Mono.just(ApiResponse.<Void>success("Space deleted successfully", null)));
    }
    
    @GetMapping("/{id}/sub-spaces")
    public Flux<SubSpace> getSubSpaces(@PathVariable Long id) {
        return spaceService.getSubSpacesBySpaceId(id);
    }
    
    @PostMapping("/sub-spaces")
    public Mono<ApiResponse<SubSpace>> createSubSpace(@Valid @RequestBody SubSpace subSpace,
                                                       Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return spaceService.createSubSpace(subSpace, userId)
                .map(created -> ApiResponse.success("SubSpace created successfully", created));
    }
}

