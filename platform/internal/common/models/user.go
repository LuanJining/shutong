package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Username   string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Phone      string         `json:"phone" gorm:"uniqueIndex;not null;size:20;comment:手机号"`
	Email      string         `json:"email" gorm:"size:100;comment:邮箱"`
	Password   string         `json:"-" gorm:"not null;size:255"`
	Nickname   string         `json:"nickname" gorm:"size:50"`
	Avatar     string         `json:"avatar" gorm:"size:255"`
	Department string         `json:"department" gorm:"size:100;comment:所属部门"`
	Company    string         `json:"company" gorm:"size:100;comment:所属企业"`
	Status     int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	LastLogin  *time.Time     `json:"last_login"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        RoleName       `json:"name" gorm:"uniqueIndex;not null;size:50"`
	DisplayName string         `json:"display_name" gorm:"size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        PermissionName `json:"name" gorm:"uniqueIndex;not null;size:50"`
	DisplayName string         `json:"display_name" gorm:"size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Resource    string         `json:"resource" gorm:"size:50;comment:资源类型"`
	Action      string         `json:"action" gorm:"size:50;comment:操作类型"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

type RoleName string

// 全局角色常量定义（系统级）
const (
	RoleSuperAdmin      RoleName = "super_admin" // 超级管理员 - 拥有所有权限
	RoleEnterpriseAdmin RoleName = "corp_admin"  // 企业管理员 - 企业级管理权限
)

type PermissionName string

// 权限常量定义
const (
	// 内容权限
	PermissionViewAllContent        PermissionName = "view_all_content"
	PermissionCreateDocument        PermissionName = "create_document"
	PermissionDeleteDocument        PermissionName = "delete_document"
	PermissionMoveDocument          PermissionName = "move_document"
	PermissionSetDocumentPermission PermissionName = "set_document_permission"

	// 空间权限
	PermissionCreateSpace       PermissionName = "create_space"
	PermissionManageSpaceMember PermissionName = "manage_space_member"

	// 审批流权限
	PermissionConfigureWorkflow PermissionName = "configure_workflow"

	// 数据权限
	PermissionExportData       PermissionName = "export_data"
	PermissionExportAllData    PermissionName = "export_all_data"
	PermissionViewOperationLog PermissionName = "view_operation_log"
	PermissionAddDeleteUser    PermissionName = "add_delete_user"
)
