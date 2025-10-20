package com.knowledgebase.platformspring.service;

import java.util.List;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.Document;
import com.knowledgebase.platformspring.model.Step;
import com.knowledgebase.platformspring.model.Task;
import com.knowledgebase.platformspring.model.Workflow;
import com.knowledgebase.platformspring.repository.SpaceMemberRepository;
import com.knowledgebase.platformspring.repository.StepRepository;
import com.knowledgebase.platformspring.repository.TaskRepository;
import com.knowledgebase.platformspring.repository.WorkflowRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Slf4j
@Service
@RequiredArgsConstructor
public class WorkflowService {
    
    private final WorkflowRepository workflowRepository;
    private final TaskRepository taskRepository;
    private final StepRepository stepRepository;
    private final SpaceMemberRepository spaceMemberRepository;
    private final com.knowledgebase.platformspring.repository.DocumentRepository documentRepository;
    
    /**
     * 创建包含 Step 的完整工作流（对齐 Go 版本）
     */
    public Mono<Workflow> createWorkflowWithStep(Long spaceId, Long resourceId, String resourceType, Long createdBy) {
        // 1. 创建 Workflow
        Workflow workflow = Workflow.builder()
                .name("文档发布审批流程")
                .description("用于文档发布的审批流程")
                .spaceId(spaceId)
                .status(Workflow.STATUS_PROCESSING)
                .resourceType(resourceType)
                .resourceId(resourceId)
                .createdBy(createdBy)
                .build();
        
        return workflowRepository.save(workflow)
                .flatMap(savedWorkflow -> {
                    // 2. 创建 Step
                    Step step = Step.builder()
                            .workflowId(savedWorkflow.getId())
                            .stepName("文档发布审批")
                            .stepOrder(1)
                            .stepRole("approver")  // 对应 SpaceMemberRoleApprover
                            .isRequired(true)
                            .timeoutHours(24 * 7)  // 7天
                            .status(Step.STATUS_PROCESSING)
                            .build();
                    
                    return stepRepository.save(step)
                            .map(savedStep -> {
                                savedWorkflow.setCurrentStepId(savedStep.getId());
                                return savedWorkflow;
                            });
                })
                .flatMap(workflowRepository::save);
    }
    
    /**
     * 启动工作流：创建审批任务
     */
    public Mono<Workflow> startWorkflow(Long workflowId, Long spaceId, Long createdBy) {
        return workflowRepository.findById(workflowId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Workflow not found")))
                .flatMap(workflow -> {
                    if (!Workflow.STATUS_PROCESSING.equals(workflow.getStatus())) {
                        return Mono.error(new BusinessException("工作流状态不正确"));
                    }
                    
                    // 查找第一个 Step
                    return stepRepository.findByWorkflowIdOrderByStepOrderAsc(workflowId)
                            .next()
                            .flatMap(firstStep -> {
                                // 查找空间中具有 approver 角色的成员
                                return spaceMemberRepository.findBySpaceId(spaceId)
                                        .filter(member -> member.getRoles() != null && 
                                                member.getRoles().contains("approver"))
                                        .flatMap(approver -> {
                                            // 为每个 approver 创建一个任务
                                            Task task = Task.builder()
                                                    .workflowId(workflowId)
                                                    .stepId(firstStep.getId())
                                                    .approverId(approver.getUserId())
                                                    .status(Task.STATUS_PROCESSING)
                                                    .build();
                                            return taskRepository.save(task);
                                        })
                                        .collectList()
                                        .map(tasks -> {
                                            log.info("Created {} approval tasks for workflow {}", tasks.size(), workflowId);
                                            return workflow;
                                        });
                            });
                });
    }
    
    public Mono<Workflow> createWorkflow(Workflow workflow, Long userId) {
        workflow.setCreatedBy(userId);
        workflow.setStatus(Workflow.STATUS_PROCESSING);
        
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
    
    public Mono<PaginationResponse<List<Task>>> getTasksPaginated(Long approverId, Integer page, Integer pageSize) {
        return taskRepository.findByApproverId(approverId)
                .collectList()
                .map(allTasks -> {
                    long total = allTasks.size();
                    int offset = (page - 1) * pageSize;
                    List<Task> pagedTasks = allTasks.stream()
                            .skip(offset)
                            .limit(pageSize)
                            .collect(java.util.stream.Collectors.toList());
                    
                    return PaginationResponse.of(pagedTasks, total, page, pageSize);
                });
    }
    
    public Mono<Workflow> startWorkflow(Long workflowId, Long userId) {
        return workflowRepository.findById(workflowId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Workflow not found")))
                .flatMap(workflow -> {
                    if (!Workflow.STATUS_PROCESSING.equals(workflow.getStatus())) {
                        return Mono.error(new BusinessException("工作流状态不正确"));
                    }
                    // 这里可以添加启动逻辑，比如创建第一批任务等
                    return Mono.just(workflow);
                });
    }
    
    /**
     * 审批任务
     */
    public Mono<Task> approveTask(Long taskId, Long approverId, String comment, String status) {
        return taskRepository.findById(taskId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("Task not found")))
                .flatMap(task -> {
                    // 1. 检查权限
                    if (!task.getApproverId().equals(approverId)) {
                        return Mono.error(BusinessException.forbidden("用户无权限审批任务"));
                    }
                    
                    // 2. 检查任务状态
                    if (!Task.STATUS_PROCESSING.equals(task.getStatus())) {
                        return Mono.error(new BusinessException("任务状态不正确，无法审批"));
                    }
                    
                    // 3. 更新当前任务状态和备注
                    task.setStatus(status);
                    task.setComment(comment);
                    
                    return taskRepository.save(task)
                            .flatMap(savedTask -> {
                                // 加载 Workflow 和 Step
                                return workflowRepository.findById(savedTask.getWorkflowId())
                                        .zipWith(stepRepository.findById(savedTask.getStepId()))
                                        .flatMap(tuple -> {
                                            Workflow workflow = tuple.getT1();
                                            Step step = tuple.getT2();
                                            
                                            if (Task.STATUS_APPROVED.equals(status)) {
                                                return handleApproval(savedTask, workflow, step);
                                            } else {
                                                return handleRejection(savedTask, workflow, step);
                                            }
                                        });
                            });
                });
    }
    
    /**
     * 处理审批通过逻辑
     */
    private Mono<Task> handleApproval(Task task, Workflow workflow, Step step) {
        // 1. 更新 Step 状态为 approved
        step.setStatus(Step.STATUS_APPROVED);
        
        return stepRepository.save(step)
                .flatMap(savedStep -> {
                    // 2. 更新同 Step 的其他 processing 任务为 approved_by_others
                    return taskRepository.findByStepId(task.getStepId())
                            .filter(t -> !t.getId().equals(task.getId()) && 
                                    Task.STATUS_PROCESSING.equals(t.getStatus()))
                            .flatMap(t -> {
                                t.setStatus(Task.STATUS_APPROVED_BY_OTHER);
                                return taskRepository.save(t);
                            })
                            .collectList();
                })
                .flatMap(updatedTasks -> {
                    // 3. 查找下一个步骤
                    return stepRepository.findByWorkflowIdOrderByStepOrderAsc(workflow.getId())
                            .filter(s -> s.getStepOrder() > step.getStepOrder())
                            .next()
                            .flatMap(nextStep -> {
                                // 有下一步：创建下一步的任务
                                log.info("Found next step {}, creating tasks", nextStep.getId());
                                return createTasksForStep(workflow, nextStep)
                                        .flatMap(createdTasks -> {
                                            // 更新 workflow 的 current_step_id
                                            workflow.setCurrentStepId(nextStep.getId());
                                            return workflowRepository.save(workflow);
                                        })
                                        .thenReturn(task);
                            })
                            .switchIfEmpty(Mono.defer(() -> {
                                // 没有下一步：工作流完成
                                log.info("No next step, completing workflow {}", workflow.getId());
                                workflow.setStatus(Workflow.STATUS_COMPLETED);
                                return workflowRepository.save(workflow)
                                        .flatMap(savedWorkflow -> {
                                            // 更新文档状态为 pending_publish
                                            return updateDocumentStatus(
                                                    savedWorkflow.getResourceId(),
                                                    Document.STATUS_PENDING_PUBLISH
                                            ).thenReturn(task);
                                        });
                            }));
                });
    }
    
    /**
     * 处理审批拒绝逻辑
     */
    private Mono<Task> handleRejection(Task task, Workflow workflow, Step step) {
        // 1. 更新 Step 状态为 rejected
        step.setStatus(Step.STATUS_REJECTED);
        
        return stepRepository.save(step)
                .flatMap(savedStep -> {
                    // 2. 更新同 Step 的其他 processing 任务为 rejected_by_others
                    return taskRepository.findByStepId(task.getStepId())
                            .filter(t -> !t.getId().equals(task.getId()) && 
                                    Task.STATUS_PROCESSING.equals(t.getStatus()))
                            .flatMap(t -> {
                                t.setStatus(Task.STATUS_REJECTED_BY_OTHER);
                                return taskRepository.save(t);
                            })
                            .collectList();
                })
                .flatMap(updatedTasks -> {
                    // 3. 更新 Workflow 状态为 cancelled
                    workflow.setStatus(Workflow.STATUS_CANCELLED);
                    return workflowRepository.save(workflow);
                })
                .flatMap(savedWorkflow -> {
                    // 4. 更新文档状态为 failed
                    return updateDocumentStatus(
                            savedWorkflow.getResourceId(),
                            Document.STATUS_FAILED
                    ).thenReturn(task);
                });
    }
    
    /**
     * 为指定 Step 创建任务
     */
    private Mono<List<Task>> createTasksForStep(Workflow workflow, Step step) {
        // 查找空间中具有指定角色的成员
        return spaceMemberRepository.findBySpaceId(workflow.getSpaceId())
                .filter(member -> member.getRoles() != null && 
                        member.getRoles().contains(step.getStepRole()))
                .flatMap(approver -> {
                    Task newTask = Task.builder()
                            .workflowId(workflow.getId())
                            .stepId(step.getId())
                            .approverId(approver.getUserId())
                            .taskName(step.getStepName())
                            .isRequired(step.getIsRequired())
                            .timeoutHours(step.getTimeoutHours())
                            .status(Task.STATUS_PROCESSING)
                            .build();
                    return taskRepository.save(newTask);
                })
                .collectList()
                .doOnSuccess(tasks -> 
                    log.info("Created {} tasks for step {}", tasks.size(), step.getId())
                );
    }
    
    /**
     * 更新文档状态（仅当资源类型为 document 时）
     */
    private Mono<Void> updateDocumentStatus(Long resourceId, String status) {
        return documentRepository.findById(resourceId)
                .flatMap(document -> {
                    document.setStatus(status);
                    return documentRepository.save(document);
                })
                .doOnSuccess(doc -> 
                    log.info("Updated document {} status to {}", resourceId, status)
                )
                .doOnError(e -> 
                    log.error("Failed to update document {} status: {}", resourceId, e.getMessage())
                )
                .then();
    }
}

