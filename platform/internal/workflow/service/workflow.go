package service

import (
	"gorm.io/gorm"
)

type WorkflowService struct {
	db *gorm.DB
}

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{
		db: db,
	}
}
