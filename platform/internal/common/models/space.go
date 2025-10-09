package model

import (
	"time"

	"gorm.io/gorm"
)

// Space 一级知识库模型
type Space struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Type        SpaceType      `json:"type" gorm:"size:50;comment:空间类型:department,project,team"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedBy   uint           `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Creator   User       `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members   []User     `json:"members" gorm:"many2many:space_members;"`
	SubSpaces []SubSpace `json:"sub_spaces" gorm:"foreignKey:SpaceID"`
}

// SubSpace 二级知识库模型
type SubSpace struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedBy   uint           `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	SpaceID uint    `json:"space_id" gorm:"not null"`
	Classes []Class `json:"classes" gorm:"foreignKey:SubSpaceID"`
}

// 知识分类
type Class struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedBy   uint           `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	SubSpaceID uint `json:"sub_space_id" gorm:"not null"`
}

// SpaceMember 空间成员关联表
type SpaceMember struct {
	SpaceID uint            `json:"space_id" gorm:"primaryKey"`
	UserID  uint            `json:"user_id" gorm:"primaryKey"`
	Role    SpaceMemberRole `json:"role" gorm:"not null;size:20;comment:空间角色:owner,admin,editor,viewer"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SpaceMemberRole 空间成员角色
type SpaceMemberRole string

const (
	SpaceMemberRoleAdmin    SpaceMemberRole = "admin"    // 空间管理员
	SpaceMemberRoleApprover SpaceMemberRole = "approver" // 空间审批者
	SpaceMemberRoleEditor   SpaceMemberRole = "editor"   // 空间编辑者
	SpaceMemberRoleReader   SpaceMemberRole = "reader"   // 空间只读者
)

var SpaceMemberRoleMap = map[SpaceMemberRole]int{
	SpaceMemberRoleAdmin:    1,
	SpaceMemberRoleApprover: 2,
	SpaceMemberRoleEditor:   3,
	SpaceMemberRoleReader:   4,
}

// SpaceType 空间类型
type SpaceType string

const (
	SpaceTypeDepartment SpaceType = "department"
	SpaceTypeProject    SpaceType = "project"
	SpaceTypeTeam       SpaceType = "team"
)

type CreateSubSpaceRequest struct {
	SpaceID     uint   `json:"space_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateClassRequest struct {
	SubSpaceID  uint   `json:"sub_space_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
