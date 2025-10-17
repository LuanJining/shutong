package com.knowledgebase.platformspring.controller;

import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.model.Task;
import com.knowledgebase.platformspring.model.Workflow;
import com.knowledgebase.platformspring.service.WorkflowService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.Data;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Tag(name = "工作流管理", description = "审批流程和任务管理相关接口")
@RestController
@RequestMapping("/api/v1/workflows")
@RequiredArgsConstructor
@SecurityRequirement(name = "bearerAuth")
public class WorkflowController {
    
    private final WorkflowService workflowService;
    
    @Operation(summary = "创建工作流", description = "创建新的审批工作流")
    @PostMapping
    public Mono<ApiResponse<Workflow>> createWorkflow(@Valid @RequestBody Workflow workflow,
                                                       Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.createWorkflow(workflow, userId)
                .map(created -> ApiResponse.success("Workflow created successfully", created));
    }
    
    @Operation(summary = "获取工作流详情", description = "根据ID获取工作流的详细信息")
    @GetMapping("/{id}")
    public Mono<ApiResponse<Workflow>> getWorkflow(@PathVariable Long id) {
        return workflowService.getWorkflowById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取空间工作流", description = "获取指定空间下的所有工作流")
    @GetMapping("/space/{spaceId}")
    public Flux<Workflow> getWorkflowsBySpace(@PathVariable Long spaceId) {
        return workflowService.getWorkflowsBySpaceId(spaceId);
    }
    
    @Operation(summary = "获取我的待办任务", description = "获取当前用户需要审批的所有任务")
    @GetMapping("/tasks/my")
    public Flux<Task> getMyTasks(Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.getTasksByApproverId(userId);
    }
    
    @Operation(summary = "审批任务", description = "对指定任务进行审批（通过或拒绝）")
    @PostMapping("/tasks/{id}/approve")
    public Mono<ApiResponse<Task>> approveTask(@PathVariable Long id,
                                                @RequestBody ApproveTaskRequest request,
                                                Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.approveTask(id, userId, request.getComment(), request.isApproved())
                .map(task -> ApiResponse.success("Task processed successfully", task));
    }
    
    @Data
    public static class ApproveTaskRequest {
        private String comment;
        private boolean approved;
    }
}

