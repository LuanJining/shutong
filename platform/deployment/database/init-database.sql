-- 知识库平台数据库初始化脚本
-- 初始化基础数据（表结构由代码自动migration创建）
-- 注意：此脚本假设表已经通过代码的migration创建
-- 插入基础权限数据
INSERT INTO permissions (name, display_name, description, resource, action) VALUES
('view_all_content', '查看所有内容', '查看所有内容权限', 'content', 'view'),
('create_doc', '创建文档', '创建文档权限', 'document', 'create'),
('delete_doc', '删除文档', '删除文档权限', 'document', 'delete'),
('move_doc', '移动文档', '移动文档权限', 'document', 'move'),
('set_doc_permission', '设置文档权限', '设置文档权限权限', 'document', 'set_permission'),
('create_space', '创建知识空间', '创建知识空间权限', 'space', 'create'),
('manage_space_members', '管理空间成员', '管理空间成员权限', 'space', 'manage_members'),
('configure_workflow', '配置审批流', '配置审批流权限', 'workflow', 'configure'),
('export_data', '导出数据', '导出数据权限', 'data', 'export'),
('export_all_data', '导出全部数据', '导出全部数据权限', 'data', 'export_all'),
('view_operation_log', '查看操作日志', '查看操作日志权限', 'log', 'view'),
('manage_users', '管理用户', '管理用户权限', 'user', 'manage')
ON CONFLICT (name) DO UPDATE SET
display_name = EXCLUDED.display_name,
description = EXCLUDED.description,
resource = EXCLUDED.resource,
action = EXCLUDED.action;

-- 插入基础角色数据（仅全局角色）
INSERT INTO roles (name, display_name, description) VALUES
('super_admin', '超级管理员', '拥有系统最高权限'),
('corp_admin', '企业管理员', '负责日常运维、空间管理、用户权限分配')
ON CONFLICT (name) DO UPDATE SET
display_name = EXCLUDED.display_name,
description = EXCLUDED.description;

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'super_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 为企业管理员角色分配管理权限（不包含用户管理权限）
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'corp_admin' AND p.name NOT IN ('manage_users')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 创建超级管理员用户（密码：admin123）
INSERT INTO users (username, phone, email, password, nickname, department, company, status) VALUES
('admin', '13800000000', 'admin@example.com', '$2a$10$qBikQOmbPSkLawCnswHqBuDsWEtcPKPEq0KSmj4opC2UgA2.qsWYq', '超级管理员', 'IT部门', '示例公司', 1)
ON CONFLICT (username) DO UPDATE SET
phone = EXCLUDED.phone,
email = EXCLUDED.email,
password = EXCLUDED.password,
nickname = EXCLUDED.nickname,
department = EXCLUDED.department,
company = EXCLUDED.company;

-- 为超级管理员用户分配角色
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin' AND r.name = 'super_admin'
ON CONFLICT (user_id, role_id) DO NOTHING;