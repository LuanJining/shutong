package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.Task;
import com.knowledgebase.platformspring.model.Workflow;
import com.knowledgebase.platformspring.repository.TaskRepository;
import com.knowledgebase.platformspring.repository.WorkflowRepository;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
public class WorkflowService {
    
    private final WorkflowRepository workflowRepository;
    private final TaskRepository taskRepository;
    
    public Mono<Workflow> createWorkflow(Workflow workflow, Long userId) {
        workflow.setCreatedBy(userId);
        workflow.setStatus(Workflow.STATUS_PROCESSING);
        workflow.setCreatedAt(LocalDateTime.now());
        workflow.setUpdatedAt(LocalDateTime.now());
        
        return workflowRepository.save(workflow);
    }
    
    public Mono<Workflow> getWorkflowById(Long id) {
        return workflowRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Workflow not found")));
    }
    
    public Flux<Workflow> getWorkflowsBySpaceId(Long spaceId) {
        return workflowRepository.findBySpaceId(spaceId);
    }
    
    public Flux<Task> getTasksByApproverId(Long approverId) {
        return taskRepository.findByApproverId(approverId);
    }
    
    public Mono<Task> approveTask(Long taskId, Long approverId, String comment, boolean approved) {
        return taskRepository.findById(taskId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Task not found")))
                .flatMap(task -> {
                    if (!task.getApproverId().equals(approverId)) {
                        return Mono.error(BusinessException.forbidden("You are not the approver of this task"));
                    }
                    
                    if (!Task.STATUS_PROCESSING.equals(task.getStatus())) {
                        return Mono.error(new BusinessException("Task is not in processing status"));
                    }
                    
                    task.setStatus(approved ? Task.STATUS_APPROVED : Task.STATUS_REJECTED);
                    task.setComment(comment);
                    task.setUpdatedAt(LocalDateTime.now());
                    
                    return taskRepository.save(task)
                            .flatMap(savedTask -> updateWorkflowStatus(savedTask.getWorkflowId())
                                    .thenReturn(savedTask));
                });
    }
    
    private Mono<Void> updateWorkflowStatus(Long workflowId) {
        return taskRepository.findByWorkflowId(workflowId)
                .collectList()
                .flatMap(tasks -> {
                    boolean allCompleted = tasks.stream()
                            .allMatch(task -> !Task.STATUS_PROCESSING.equals(task.getStatus()));
                    
                    boolean anyRejected = tasks.stream()
                            .anyMatch(task -> Task.STATUS_REJECTED.equals(task.getStatus()));
                    
                    if (allCompleted) {
                        return workflowRepository.findById(workflowId)
                                .flatMap(workflow -> {
                                    workflow.setStatus(anyRejected ? 
                                            Workflow.STATUS_CANCELLED : Workflow.STATUS_COMPLETED);
                                    workflow.setUpdatedAt(LocalDateTime.now());
                                    return workflowRepository.save(workflow);
                                })
                                .then();
                    }
                    return Mono.empty();
                });
    }
}

