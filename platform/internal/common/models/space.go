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
	Creator User   `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members []User `json:"members" gorm:"many2many:space_members;"`
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
	SpaceID uint  `json:"space_id" gorm:"not null"`
	Space   Space `json:"space" gorm:"foreignKey:SpaceID"`
}

// 知识分类
type Class struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	SubSpaceID uint     `json:"sub_space_id" gorm:"not null"`
	SubSpace   SubSpace `json:"sub_space" gorm:"foreignKey:SubSpaceID"`
}

// SpaceMember 空间成员关联表
type SpaceMember struct {
	SpaceID uint   `gorm:"primaryKey"`
	UserID  uint   `gorm:"primaryKey"`
	Role    string `json:"role" gorm:"size:50;comment:在空间中的角色:space_admin,content_reviewer,content_editor,read_only_user"`

	// 关联关系
	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Space Space `json:"space" gorm:"foreignKey:SpaceID"`
}

// SpaceType 空间类型
type SpaceType string

const (
	SpaceTypeDepartment SpaceType = "department"
	SpaceTypeProject    SpaceType = "project"
	SpaceTypeTeam       SpaceType = "team"
)
