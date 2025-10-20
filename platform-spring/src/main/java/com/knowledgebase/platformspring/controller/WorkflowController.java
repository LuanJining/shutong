package com.knowledgebase.platformspring.controller;

import java.util.List;

import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.model.Task;
import com.knowledgebase.platformspring.model.Workflow;
import com.knowledgebase.platformspring.service.WorkflowService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.Data;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Tag(name = "工作流管理", description = "审批流程和任务管理相关接口")
@RestController
@RequestMapping("/api/v1/workflow")
@RequiredArgsConstructor
@SecurityRequirement(name = "bearerAuth")
public class WorkflowController {
    
    private final WorkflowService workflowService;
    
    @Operation(summary = "创建工作流", description = "创建新的审批工作流")
    @PostMapping("/workflows")
    public Mono<ApiResponse<Workflow>> createWorkflow(@Valid @RequestBody Workflow workflow,
                                                       Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.createWorkflow(workflow, userId)
                .map(created -> ApiResponse.success("工作流创建成功", created));
    }
    
    @Operation(summary = "启动工作流", description = "启动指定的审批工作流")
    @PostMapping("/workflows/{id}/start")
    public Mono<ApiResponse<Workflow>> startWorkflow(@PathVariable Long id,
                                                      Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.startWorkflow(id, userId)
                .map(workflow -> ApiResponse.success("工作流启动成功", workflow));
    }
    
    @Operation(summary = "获取工作流详情", description = "根据ID获取工作流的详细信息")
    @GetMapping("/workflows/{id}")
    public Mono<ApiResponse<Workflow>> getWorkflow(@PathVariable Long id) {
        return workflowService.getWorkflowById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取我的待办任务", description = "获取当前用户需要审批的所有任务（分页）")
    @GetMapping("/tasks")
    public Mono<ApiResponse<PaginationResponse<List<Task>>>> getTasks(
            Authentication authentication,
            @RequestParam(defaultValue = "1") Integer page,
            @RequestParam(defaultValue = "10") Integer pageSize) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.getTasksPaginated(userId, page, pageSize)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "审批任务", description = "对指定任务进行审批（通过或拒绝）")
    @PostMapping("/tasks/approve")
    public Mono<ApiResponse<Task>> approveTask(@RequestBody ApproveTaskRequest request,
                                                Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return workflowService.approveTask(request.getTaskId(), userId, request.getComment(), 
                                          Task.STATUS_APPROVED.equals(request.getStatus()))
                .map(task -> ApiResponse.success("审批任务成功", task));
    }
    
    @Data
    public static class ApproveTaskRequest {
        private Long taskId;
        private String comment;
        private String status; // approved, rejected
    }
}

