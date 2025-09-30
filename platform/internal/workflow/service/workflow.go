package service

import (
	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
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
