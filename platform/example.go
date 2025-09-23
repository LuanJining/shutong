package platform

type RoleType int

const (
	RoleTypeSuperAdmin RoleType = iota
	RoleTypeCorpAdmin
	RoleTypeSpaceAdmin
	RoleTypeContentReviewer
	RoleTypeContentEditor
	RoleTypeContentViewer
)

var RoleTypeNames = map[RoleType]string{
	RoleTypeSuperAdmin: "super_admin",
	RoleTypeCorpAdmin: "corp_admin",
	RoleTypeSpaceAdmin: "space_admin",
	RoleTypeContentReviewer: "content_reviewer",
	RoleTypeContentEditor: "content_editor",
	RoleTypeContentViewer: "content_viewer",
}

type PermissionType int

const (
	PermissionTypeSuper PermissionType = iota
	PermissionTypeReviewDoc
	PermissionTypeCreateDoc
	PermissionTypeDeleteDoc
	PermissionTypeEditDoc
	PermissionTypeSetDocPermission
	PermissionTypeCreateSpace
	PermissionTypeManageUser
	PermissionTypeManageWorkflow
	PermissionTypeExportData
	PermissionTypeExportAllData
	PermissionTypeReviewLog
	PermissionTypeCreateDeleteUser
)

var PermissionTypeNames = map[PermissionType]string{
	PermissionTypeSuper: "*",
	PermissionTypeReviewDoc: "review_doc",
	PermissionTypeCreateDoc: "create_doc",
	PermissionTypeDeleteDoc: "delete_doc",
	PermissionTypeEditDoc: "edit_doc",
	PermissionTypeSetDocPermission: "set_doc_permission",
	PermissionTypeCreateSpace: "create_space",
	PermissionTypeManageUser: "manage_user",
	PermissionTypeManageWorkflow: "manage_workflow",
	PermissionTypeExportData: "export_data",
	PermissionTypeExportAllData: "export_all_data",
	PermissionTypeReviewLog: "review_log",
	PermissionTypeCreateDeleteUser: "create_delete_user",
}

