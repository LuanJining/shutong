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
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	DisplayName string         `json:"display_name" gorm:"size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
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

// 角色常量定义
const (
	RoleSuperAdmin      = "super_admin"
	RoleEnterpriseAdmin = "enterprise_admin"
	RoleSpaceAdmin      = "space_admin"
	RoleContentReviewer = "content_reviewer"
	RoleContentEditor   = "content_editor"
	RoleReadOnlyUser    = "read_only_user"
)

// 权限常量定义
const (
	// 内容权限
	PermissionViewAllContent        = "view_all_content"
	PermissionCreateDocument        = "create_document"
	PermissionDeleteDocument        = "delete_document"
	PermissionMoveDocument          = "move_document"
	PermissionSetDocumentPermission = "set_document_permission"

	// 空间权限
	PermissionCreateSpace       = "create_space"
	PermissionManageSpaceMember = "manage_space_member"

	// 审批流权限
	PermissionConfigureWorkflow = "configure_workflow"

	// 数据权限
	PermissionExportData       = "export_data"
	PermissionExportAllData    = "export_all_data"
	PermissionViewOperationLog = "view_operation_log"
	PermissionAddDeleteUser    = "add_delete_user"
)
