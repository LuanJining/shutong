package service

import (
	"errors"
	"fmt"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/models"
	"gorm.io/gorm"
)

// User 用户模型（简化版，用于workflow服务）
type User struct {
	ID     uint `gorm:"primaryKey"`
	Status int  `gorm:"default:1"`
}

type WorkflowService struct {
	db *gorm.DB
}

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

// CreateWorkflow 创建审批流程
func (s *WorkflowService) CreateWorkflow(req *models.CreateWorkflowRequest, createdBy uint) (*models.WorkflowDefinition, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建流程定义
	workflow := &models.WorkflowDefinition{
		Name:        req.Name,
		Description: req.Description,
		SpaceID:     req.SpaceID,
		Priority:    req.Priority,
		CreatedBy:   createdBy,
		IsActive:    true,
	}

	if err := tx.Create(workflow).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// 创建流程步骤
	for _, stepReq := range req.Steps {
		step := &models.WorkflowStep{
			WorkflowID:   workflow.ID,
			StepName:     stepReq.StepName,
			StepOrder:    stepReq.StepOrder,
			ApproverType: stepReq.ApproverType,
			ApproverID:   stepReq.ApproverID,
			IsRequired:   stepReq.IsRequired,
			TimeoutHours: stepReq.TimeoutHours,
		}

		if err := tx.Create(step).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create workflow step: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新查询包含步骤的完整数据
	var result models.WorkflowDefinition
	if err := s.db.Preload("Steps").First(&result, workflow.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load workflow with steps: %w", err)
	}

	return &result, nil
}

// GetWorkflows 获取流程列表
func (s *WorkflowService) GetWorkflows(spaceID uint, page, pageSize int) (*models.PaginationResponse, error) {
	var workflows []models.WorkflowDefinition
	var total int64

	query := s.db.Model(&models.WorkflowDefinition{})
	if spaceID > 0 {
		query = query.Where("space_id = ?", spaceID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count workflows: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Steps").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&workflows).Error; err != nil {
		return nil, fmt.Errorf("failed to get workflows: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.PaginationResponse{
		Items:      workflows,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetWorkflow 获取流程详情
func (s *WorkflowService) GetWorkflow(id uint) (*models.WorkflowDefinition, error) {
	var workflow models.WorkflowDefinition
	if err := s.db.Preload("Steps").First(&workflow, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}
	return &workflow, nil
}

// UpdateWorkflow 更新审批流程
func (s *WorkflowService) UpdateWorkflow(id uint, req *models.CreateWorkflowRequest, updatedBy uint) (*models.WorkflowDefinition, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查流程是否存在
	var workflow models.WorkflowDefinition
	if err := tx.First(&workflow, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// 更新流程基本信息
	workflow.Name = req.Name
	workflow.Description = req.Description
	workflow.SpaceID = req.SpaceID
	workflow.Priority = req.Priority
	workflow.UpdatedBy = &updatedBy

	if err := tx.Save(&workflow).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	// 删除原有步骤
	if err := tx.Where("workflow_id = ?", id).Delete(&models.WorkflowStep{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete old steps: %w", err)
	}

	// 创建新步骤
	for _, stepReq := range req.Steps {
		step := &models.WorkflowStep{
			WorkflowID:   workflow.ID,
			StepName:     stepReq.StepName,
			StepOrder:    stepReq.StepOrder,
			ApproverType: stepReq.ApproverType,
			ApproverID:   stepReq.ApproverID,
			IsRequired:   stepReq.IsRequired,
			TimeoutHours: stepReq.TimeoutHours,
		}

		if err := tx.Create(step).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create workflow step: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新查询包含步骤的完整数据
	var result models.WorkflowDefinition
	if err := s.db.Preload("Steps").First(&result, workflow.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load workflow with steps: %w", err)
	}

	return &result, nil
}

// DeleteWorkflow 删除审批流程
func (s *WorkflowService) DeleteWorkflow(id uint) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查流程是否存在
	var workflow models.WorkflowDefinition
	if err := tx.First(&workflow, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("workflow not found")
		}
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	// 检查是否有正在进行的实例
	var count int64
	if err := tx.Model(&models.WorkflowInstance{}).Where("workflow_id = ? AND status IN (?)", id, []string{models.StatusPending, models.StatusInProgress}).Count(&count).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check workflow instances: %w", err)
	}

	if count > 0 {
		tx.Rollback()
		return fmt.Errorf("cannot delete workflow with active instances")
	}

	// 删除相关数据
	if err := tx.Where("workflow_id = ?", id).Delete(&models.WorkflowStep{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete workflow steps: %w", err)
	}

	if err := tx.Delete(&workflow).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// StartWorkflow 启动流程实例
func (s *WorkflowService) StartWorkflow(req *models.StartWorkflowRequest, createdBy uint) (*models.WorkflowInstance, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查流程是否存在且启用
	var workflow models.WorkflowDefinition
	if err := tx.Preload("Steps").Where("id = ? AND is_active = ?", req.WorkflowID, true).First(&workflow).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow not found or inactive")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// 创建流程实例
	instance := &models.WorkflowInstance{
		WorkflowID:   req.WorkflowID,
		Title:        req.Title,
		Description:  req.Description,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		SpaceID:      req.SpaceID,
		Priority:     req.Priority,
		CreatedBy:    createdBy,
		Status:       models.StatusPending,
	}

	if err := tx.Create(instance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create workflow instance: %w", err)
	}

	// 创建审批任务
	for _, step := range workflow.Steps {
		// 根据审批人类型确定审批人ID
		assigneeID := s.determineAssigneeID(&step, req.SpaceID, tx)
		if assigneeID == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("failed to determine assignee for step: %s", step.StepName)
		}

		task := &models.WorkflowTask{
			InstanceID: instance.ID,
			StepID:     step.ID,
			AssigneeID: assigneeID,
			Status:     models.StatusPending,
			AssignedAt: time.Now(),
		}

		if err := tx.Create(task).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create workflow task: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新查询包含任务的完整数据
	var result models.WorkflowInstance
	if err := s.db.Preload("Tasks").Preload("Workflow").First(&result, instance.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load workflow instance with tasks: %w", err)
	}

	return &result, nil
}

// GetInstances 获取流程实例列表
func (s *WorkflowService) GetInstances(workflowID uint, spaceID uint, status string, page, pageSize int) (*models.PaginationResponse, error) {
	var instances []models.WorkflowInstance
	var total int64

	query := s.db.Model(&models.WorkflowInstance{})
	if workflowID > 0 {
		query = query.Where("workflow_id = ?", workflowID)
	}
	if spaceID > 0 {
		query = query.Where("space_id = ?", spaceID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count instances: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Workflow").
		Preload("Tasks").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&instances).Error; err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.PaginationResponse{
		Items:      instances,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetInstance 获取流程实例详情
func (s *WorkflowService) GetInstance(id uint) (*models.WorkflowInstance, error) {
	var instance models.WorkflowInstance
	if err := s.db.Preload("Workflow").
		Preload("Tasks").
		Preload("Tasks.Step").
		First(&instance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}
	return &instance, nil
}

// CancelInstance 取消流程实例
func (s *WorkflowService) CancelInstance(id uint, userID uint) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取实例
	var instance models.WorkflowInstance
	if err := tx.First(&instance, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("instance not found")
		}
		return fmt.Errorf("failed to get instance: %w", err)
	}

	// 检查权限（只有创建者可以取消）
	if instance.CreatedBy != userID {
		tx.Rollback()
		return fmt.Errorf("unauthorized to cancel this instance")
	}

	// 检查状态
	if instance.Status != models.StatusPending && instance.Status != models.StatusInProgress {
		tx.Rollback()
		return fmt.Errorf("instance cannot be cancelled in current status")
	}

	// 更新实例状态
	now := time.Now()
	if err := tx.Model(&instance).Updates(map[string]interface{}{
		"status":       models.StatusCancelled,
		"completed_at": &now,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update instance: %w", err)
	}

	// 取消所有待处理的任务
	if err := tx.Model(&models.WorkflowTask{}).
		Where("instance_id = ? AND status = ?", id, models.StatusPending).
		Updates(map[string]interface{}{
			"status":       models.StatusCancelled,
			"completed_at": &now,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to cancel tasks: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetMyTasks 获取我的待办任务
func (s *WorkflowService) GetMyTasks(userID uint, page, pageSize int) (*models.PaginationResponse, error) {
	var tasks []models.WorkflowTask
	var total int64

	query := s.db.Model(&models.WorkflowTask{}).
		Where("assignee_id = ? AND status = ?", userID, models.StatusPending).
		Joins("JOIN workflow_instances ON workflow_tasks.instance_id = workflow_instances.id").
		Where("workflow_instances.status = ?", models.StatusPending)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count tasks: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Instance").
		Preload("Step").
		Order("assigned_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.PaginationResponse{
		Items:      tasks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ApproveTask 审批通过
func (s *WorkflowService) ApproveTask(taskID uint, userID uint, comment string) error {
	return s.processTask(taskID, userID, models.StatusApproved, comment)
}

// RejectTask 审批拒绝
func (s *WorkflowService) RejectTask(taskID uint, userID uint, comment string) error {
	return s.processTask(taskID, userID, models.StatusRejected, comment)
}

// processTask 处理审批任务
func (s *WorkflowService) processTask(taskID uint, userID uint, status string, comment string) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取任务
	var task models.WorkflowTask
	if err := tx.Preload("Instance").Preload("Step").First(&task, taskID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task not found")
		}
		return fmt.Errorf("failed to get task: %w", err)
	}

	// 检查权限
	if task.AssigneeID != userID {
		tx.Rollback()
		return fmt.Errorf("unauthorized to process this task")
	}

	// 检查状态
	if task.Status != models.StatusPending {
		tx.Rollback()
		return fmt.Errorf("task is not pending")
	}

	// 更新任务状态
	now := time.Now()
	task.Status = status
	task.Comment = comment
	task.CompletedAt = &now

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update task: %w", err)
	}

	// 更新流程实例状态
	if status == models.StatusRejected {
		// 拒绝则整个流程结束
		if err := tx.Model(&task.Instance).Updates(map[string]interface{}{
			"status":       models.StatusRejected,
			"completed_at": &now,
		}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update workflow instance: %w", err)
		}
	} else {
		// 通过则检查是否还有后续步骤
		var nextStep models.WorkflowStep
		if err := tx.Where("workflow_id = ? AND step_order > ?", task.Step.WorkflowID, task.Step.StepOrder).
			Order("step_order ASC").
			First(&nextStep).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 没有后续步骤，流程完成
				if err := tx.Model(&task.Instance).Updates(map[string]interface{}{
					"status":       models.StatusApproved,
					"completed_at": &now,
				}).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to complete workflow instance: %w", err)
				}
			} else {
				tx.Rollback()
				return fmt.Errorf("failed to check next step: %w", err)
			}
		} else {
			// 创建下一个任务
			assigneeID := s.determineAssigneeID(&nextStep, task.Instance.SpaceID, tx)
			if assigneeID == 0 {
				tx.Rollback()
				return fmt.Errorf("failed to determine assignee for next step")
			}

			nextTask := &models.WorkflowTask{
				InstanceID: task.Instance.ID,
				StepID:     nextStep.ID,
				AssigneeID: assigneeID,
				Status:     models.StatusPending,
				AssignedAt: now,
			}

			if err := tx.Create(nextTask).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create next task: %w", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// determineAssigneeID 确定审批人ID
func (s *WorkflowService) determineAssigneeID(step *models.WorkflowStep, spaceID uint, tx *gorm.DB) uint {
	switch step.ApproverType {
	case models.ApproverTypeUser:
		// 检查指定用户是否在当前空间中，超级管理员和企业管理员可以跨空间
		var userID uint
		err := tx.Raw(`
			SELECT u.id FROM users u 
			LEFT JOIN space_members sm ON u.id = sm.user_id AND sm.space_id = ?
			JOIN user_roles ur ON u.id = ur.user_id 
			JOIN roles r ON ur.role_id = r.id 
			WHERE u.id = ? AND u.status = 1 
			AND (sm.user_id IS NOT NULL OR r.name IN ('super_admin', 'enterprise_admin'))
			LIMIT 1
		`, spaceID, step.ApproverID).Scan(&userID).Error
		if err != nil || userID == 0 {
			return 0
		}
		return userID
	case models.ApproverTypeRole:
		// 查询角色对应的用户，超级管理员和企业管理员可以跨空间
		var userID uint
		err := tx.Raw(`
			SELECT u.id FROM users u 
			JOIN user_roles ur ON u.id = ur.user_id 
			JOIN roles r ON ur.role_id = r.id 
			LEFT JOIN space_members sm ON u.id = sm.user_id AND sm.space_id = ?
			WHERE r.name = ? AND u.status = 1 
			AND (sm.user_id IS NOT NULL OR r.name IN ('super_admin', 'enterprise_admin'))
			LIMIT 1
		`, spaceID, step.ApproverID).Scan(&userID).Error
		if err != nil || userID == 0 {
			return 0
		}
		return userID
	case models.ApproverTypeSpaceAdmin:
		// 查询空间管理员，超级管理员和企业管理员也可以作为空间管理员
		var userID uint
		err := tx.Raw(`
			SELECT u.id FROM users u 
			LEFT JOIN space_members sm ON u.id = sm.user_id AND sm.space_id = ? AND sm.role = 'space_admin'
			JOIN user_roles ur ON u.id = ur.user_id 
			JOIN roles r ON ur.role_id = r.id 
			WHERE u.status = 1 
			AND (sm.user_id IS NOT NULL OR r.name IN ('super_admin', 'enterprise_admin'))
			LIMIT 1
		`, spaceID).Scan(&userID).Error
		if err != nil || userID == 0 {
			return 0
		}
		return userID
	default:
		return 0
	}
}

// TransferTask 转交任务
func (s *WorkflowService) TransferTask(taskID uint, fromUserID uint, toUserID uint, comment string) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取任务
	var task models.WorkflowTask
	if err := tx.Preload("Instance").First(&task, taskID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task not found")
		}
		return fmt.Errorf("failed to get task: %w", err)
	}

	// 检查权限
	if task.AssigneeID != fromUserID {
		tx.Rollback()
		return fmt.Errorf("unauthorized to transfer this task")
	}

	// 检查状态
	if task.Status != models.StatusPending {
		tx.Rollback()
		return fmt.Errorf("task is not pending")
	}

	// 检查目标用户是否存在
	var userCount int64
	if err := tx.Model(&User{}).Where("id = ? AND status = 1", toUserID).Count(&userCount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check target user: %w", err)
	}
	if userCount == 0 {
		tx.Rollback()
		return fmt.Errorf("target user not found or inactive")
	}

	// 更新任务
	now := time.Now()
	task.AssigneeID = toUserID
	task.Comment = comment
	task.TransferredAt = &now
	task.TransferredBy = &fromUserID

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update task: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
