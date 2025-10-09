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
	userList, err := s.iamClient.GetSpaceMemebersByRole(user, workflow.SpaceID, currentStep.StepRole)
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

func (s *WorkflowService) GetTasks(user *model.User, page int, pageSize int) (model.PaginationResponse, error) {
	var tasks []model.Task
	var total int64

	// 构建查询条件
	query := s.db.Where("approver_id = ?", user.ID)

	// 获取总数
	if err := query.Model(&model.Task{}).Count(&total).Error; err != nil {
		return model.PaginationResponse{}, err
	}

	// 分页查询，预加载关联的 Workflow
	offset := (page - 1) * pageSize
	if err := query.Preload("Workflow").
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return model.PaginationResponse{}, err
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return model.PaginationResponse{
		Items:      tasks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *WorkflowService) ApproveTask(req *model.ApproveTaskRequest, user *model.User) (*model.Task, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 加载任务及其关联数据
	task := &model.Task{}
	err := tx.Preload("Workflow").Preload("Step").Where("id = ?", req.TaskID).First(task).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 检查权限
	if task.ApproverID != user.ID {
		tx.Rollback()
		return nil, errors.New("用户无权限审批任务")
	}

	// 检查任务状态
	if task.Status != model.TaskStatusProcessing {
		tx.Rollback()
		return nil, errors.New("任务状态不正确，无法审批")
	}

	// 更新当前任务状态和备注
	task.Status = req.Status
	task.Comment = req.Comment
	if err := tx.Save(task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if req.Status == model.TaskStatusApproved {
		// 批准：更新 Step 状态
		task.Step.Status = model.StepStatusApproved
		if err := tx.Save(&task.Step).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 更新该 Step 的其他任务为 ApprovedByOther（排除当前任务）
		if err := tx.Model(&model.Task{}).
			Where("step_id = ? AND id != ? AND status = ?", task.StepID, task.ID, model.TaskStatusProcessing).
			Update("status", model.TaskStatusApprovedByOther).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 查找下一个步骤
		nextStep := &model.Step{}
		err := tx.Where("workflow_id = ? AND step_order = ?", task.WorkflowID, task.Step.StepOrder+1).First(nextStep).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据库错误
			tx.Rollback()
			return nil, err
		}

		if err == nil && nextStep.ID != 0 {
			// 有下一步：更新 Workflow 的当前步骤
			task.Workflow.CurrentStepID = nextStep.ID
			if err := tx.Save(&task.Workflow).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			// 创建下一步的任务
			userList, err := s.iamClient.GetSpaceMemebersByRole(user, task.Workflow.SpaceID, nextStep.StepRole)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			if len(userList) == 0 {
				tx.Rollback()
				return nil, errors.New("未找到下一步骤的审批人")
			}

			tasks := make([]model.Task, len(userList))
			for i, approver := range userList {
				tasks[i] = model.Task{
					WorkflowID:       task.WorkflowID,
					StepID:           nextStep.ID,
					ApproverID:       approver.ID,
					ApproverNickName: approver.Nickname,
					TaskName:         nextStep.StepName,
					IsRequired:       nextStep.IsRequired,
					TimeoutHours:     nextStep.TimeoutHours,
					Status:           model.TaskStatusProcessing,
				}
			}
			if err := tx.Create(&tasks).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 没有下一步：工作流完成
			task.Workflow.Status = model.WorkflowStatusCompleted
			if err := tx.Save(&task.Workflow).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	} else {
		// 拒绝：更新 Step 状态
		task.Step.Status = model.StepStatusRejected
		if err := tx.Save(&task.Step).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 更新该 Step 的其他任务为 RejectedByOther（排除当前任务）
		if err := tx.Model(&model.Task{}).
			Where("step_id = ? AND id != ? AND status = ?", task.StepID, task.ID, model.TaskStatusProcessing).
			Update("status", model.TaskStatusRejectedByOther).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 更新 Workflow 状态为已取消
		task.Workflow.Status = model.WorkflowStatusCancelled
		if err := tx.Save(&task.Workflow).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return task, nil
}
