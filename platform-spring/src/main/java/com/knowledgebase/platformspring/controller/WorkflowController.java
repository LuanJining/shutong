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

import jakarta.validation.Valid;
import lombok.Data;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v1/workflows")
@RequiredArgsConstructor
public class WorkflowController {
    
    private final WorkflowService workflowService;
    
    @PostMapping
    public Mono<ApiResponse<Workflow>> createWorkflow(@Valid @RequestBody Workflow workflow,
                                                       Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.createWorkflow(workflow, userId)
                .map(created -> ApiResponse.success("Workflow created successfully", created));
    }
    
    @GetMapping("/{id}")
    public Mono<ApiResponse<Workflow>> getWorkflow(@PathVariable Long id) {
        return workflowService.getWorkflowById(id)
                .map(ApiResponse::success);
    }
    
    @GetMapping("/space/{spaceId}")
    public Flux<Workflow> getWorkflowsBySpace(@PathVariable Long spaceId) {
        return workflowService.getWorkflowsBySpaceId(spaceId);
    }
    
    @GetMapping("/tasks/my")
    public Flux<Task> getMyTasks(Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.getTasksByApproverId(userId);
    }
    
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

