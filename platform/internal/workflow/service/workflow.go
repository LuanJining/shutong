package service

import (
	"errors"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/client"
	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gorm.io/gorm"
)

// User 用户模型（简化版，用于workflow服务）
type User struct {
	ID     uint `gorm:"primaryKey"`
	Status int  `gorm:"default:1"`
}

type WorkflowService struct {
	db        *gorm.DB
	iamClient *client.IamClient
}

func NewWorkflowService(db *gorm.DB, iamClient *client.IamClient) *WorkflowService {
	return &WorkflowService{db: db, iamClient: iamClient}
}

// CreateWorkflow 创建审批流程
func (s *WorkflowService) CreateWorkflow(req *model.CreateWorkflowRequest, user *model.User) (*model.Workflow, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建主工作流记录
	result := model.Workflow{
		Name:            req.Name,
		Description:     req.Description,
		SpaceID:         req.SpaceID,
		Status:          model.WorkflowStatusProcessing,
		CreatedBy:       user.ID,
		CreatorNickName: user.Nickname,
	}
	if err := tx.Create(&result).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建步骤和任务
	for stepIndex, step := range req.Steps {
		// 设置步骤信息
		step.WorkflowID = result.ID
		step.Status = model.StepStatusProcessing

		// 创建步骤
		if err := tx.Create(&step).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 更新步骤引用
		req.Steps[stepIndex] = step
	}

	// 更新结果中的步骤信息
	result.Steps = req.Steps

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *WorkflowService) StartWorkflow(req *model.StartWorkflowRequest, user *model.User) (*model.Workflow, error) {
	workflow := &model.Workflow{}
	err := s.db.Where("id = ?", req.WorkflowID).First(workflow).Error

	if err != nil {
		return nil, err
	}

	if workflow.Status != model.WorkflowStatusProcessing {
		return nil, errors.New("工作流状态不正确")
	}

	if workflow.CreatedBy != user.ID {
		return nil, errors.New("用户无权限启动工作流")
	}

	currentStep := &model.Step{}
	err = s.db.Where("workflow_id = ?", req.WorkflowID).Where("step_order = ?", 1).First(currentStep).Error
	if err != nil {
		return nil, err
	}

	// 创建任务 先通过iam获取有权限的user列表
	userList, err := s.iamClient.GetSpaceMemebersByRole(user, workflow.SpaceID, string(model.SpaceMemberRoleOwner))
	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, len(userList))
	for i, approver := range userList {
		tasks[i] = model.Task{
			WorkflowID:       req.WorkflowID,
			StepID:           currentStep.ID,
			ApproverID:       approver.ID,
			ApproverNickName: approver.Nickname,
			TaskName:         currentStep.StepName,
			IsRequired:       currentStep.IsRequired,
			TimeoutHours:     currentStep.TimeoutHours,
			Status:           model.TaskStatusProcessing,
		}
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建任务
	if err := tx.Create(&tasks).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新工作流
	workflow.CurrentStepID = currentStep.ID
	if err := tx.Save(workflow).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return workflow, nil
}
