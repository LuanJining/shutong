package com.knowledgebase.platformspring.controller;

import java.util.List;

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
import com.knowledgebase.platformspring.dto.SpaceWithHierarchy;
import com.knowledgebase.platformspring.model.KnowledgeClass;
import com.knowledgebase.platformspring.model.Space;
import com.knowledgebase.platformspring.model.SpaceMember;
import com.knowledgebase.platformspring.model.SubSpace;
import com.knowledgebase.platformspring.service.SpaceMemberService;
import com.knowledgebase.platformspring.service.SpaceService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Tag(name = "空间管理", description = "知识空间管理相关接口")
@RestController
@RequestMapping("/api/v1/spaces")
@RequiredArgsConstructor
@SecurityRequirement(name = "bearerAuth")
public class SpaceController {
    
    private final SpaceService spaceService;
    private final SpaceMemberService spaceMemberService;
    
    @Operation(summary = "获取所有空间", description = "获取当前用户可访问的所有知识空间（包含子空间和分类）")
    @GetMapping
    public Mono<ApiResponse<List<SpaceWithHierarchy>>> getAllSpaces() {
        return spaceService.getAllSpacesWithHierarchy()
                .collectList()
                .map(spaces -> ApiResponse.success("Spaces retrieved successfully", spaces));
    }
    
    @Operation(summary = "获取空间详情", description = "根据ID获取指定空间的详细信息")
    @GetMapping("/{id}")
    public Mono<ApiResponse<Space>> getSpaceById(@PathVariable Long id) {
        return spaceService.getSpaceById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "创建空间", description = "创建新的知识空间")
    @PostMapping
    public Mono<ApiResponse<Space>> createSpace(@Valid @RequestBody Space space, 
                                                 Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return spaceService.createSpace(space, userId)
                .map(created -> ApiResponse.success("Space created successfully", created));
    }
    
    @Operation(summary = "更新空间", description = "更新指定空间的信息")
    @PutMapping("/{id}")
    public Mono<ApiResponse<Space>> updateSpace(@PathVariable Long id, 
                                                 @Valid @RequestBody Space space) {
        return spaceService.updateSpace(id, space)
                .map(updated -> ApiResponse.success("Space updated successfully", updated));
    }
    
    @Operation(summary = "删除空间", description = "删除指定的知识空间")
    @DeleteMapping("/{id}")
    public Mono<ApiResponse<Void>> deleteSpace(@PathVariable Long id) {
        return spaceService.deleteSpace(id)
                .then(Mono.just(ApiResponse.<Void>success("Space deleted successfully", null)));
    }
    
    @Operation(summary = "获取子空间列表", description = "获取指定空间下的所有子空间")
    @GetMapping("/{id}/sub-spaces")
    public Flux<SubSpace> getSubSpaces(@PathVariable Long id) {
        return spaceService.getSubSpacesBySpaceId(id);
    }
    
    @Operation(summary = "创建子空间", description = "在指定空间下创建新的子空间")
    @PostMapping("/sub-spaces")
    public Mono<ApiResponse<SubSpace>> createSubSpace(@Valid @RequestBody SubSpace subSpace,
                                                       Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return spaceService.createSubSpace(subSpace, userId)
                .map(created -> ApiResponse.success("SubSpace created successfully", created));
    }

    @Operation(summary = "创建知识分类", description = "在指定子空间下创建新的知识分类")
    @PostMapping("/classes")
    public Mono<ApiResponse<KnowledgeClass>> createClass(@Valid @RequestBody KnowledgeClass knowledgeClass,
                                                        Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return spaceService.createKnowledgeClass(knowledgeClass, userId)
                .map(created -> ApiResponse.success("知识分类创建成功", created));
    }
    
    @Operation(summary = "获取空间成员", description = "获取指定空间的所有成员")
    @GetMapping("/{id}/members")
    public Flux<SpaceMember> getSpaceMembers(@PathVariable Long id) {
        return spaceMemberService.getSpaceMembers(id);
    }
    
    @Operation(summary = "添加空间成员", description = "将用户添加到空间并分配角色")
    @PostMapping("/{id}/members")
    public Mono<ApiResponse<SpaceMember>> addSpaceMember(
            @PathVariable Long id,
            @RequestBody AddMemberRequest request) {
        return spaceMemberService.addSpaceMember(id, request.userId, request.roles)
                .map(member -> ApiResponse.success("成员添加成功", member));
    }
    
    @Operation(summary = "移除空间成员", description = "从空间中移除指定用户")
    @DeleteMapping("/{id}/members/{userId}")
    public Mono<ApiResponse<Void>> removeSpaceMember(@PathVariable Long id, 
                                                      @PathVariable Long userId) {
        return spaceMemberService.removeSpaceMember(id, userId)
                .then(Mono.just(ApiResponse.<Void>success("成员移除成功", null)));
    }
    
    @Operation(summary = "更新成员角色", description = "更新空间成员的角色")
    @PutMapping("/{id}/members/{userId}")
    public Mono<ApiResponse<SpaceMember>> updateMemberRoles(
            @PathVariable Long id,
            @PathVariable Long userId,
            @RequestBody UpdateMemberRolesRequest request) {
        return spaceMemberService.updateSpaceMemberRoles(id, userId, request.roles)
                .map(member -> ApiResponse.success("成员角色更新成功", member));
    }
    
    public record AddMemberRequest(Long userId, java.util.List<String> roles) {}
    public record UpdateMemberRolesRequest(java.util.List<String> roles) {}
}

