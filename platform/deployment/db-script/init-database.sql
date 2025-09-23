-- 知识库平台数据库初始化脚本
-- 创建数据库和基础数据

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS kbase CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE kbase;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号',
    email VARCHAR(100) COMMENT '邮箱',
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    department VARCHAR(100) COMMENT '所属部门',
    company VARCHAR(100) COMMENT '所属企业',
    status TINYINT DEFAULT 1 COMMENT '1-正常 0-禁用',
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_phone (phone),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);

-- 创建角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1-正常 0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);

-- 创建权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description VARCHAR(255),
    resource VARCHAR(50) COMMENT '资源类型',
    action VARCHAR(50) COMMENT '操作类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_resource (resource),
    INDEX idx_action (action),
    INDEX idx_deleted_at (deleted_at)
);

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- 创建角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- 创建知识空间表
CREATE TABLE IF NOT EXISTS spaces (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    type VARCHAR(50) COMMENT '空间类型:department,project,team',
    status TINYINT DEFAULT 1 COMMENT '1-正常 0-禁用',
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_name (name),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_deleted_at (deleted_at)
);

-- 创建空间成员关联表
CREATE TABLE IF NOT EXISTS space_members (
    space_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    role VARCHAR(50) COMMENT '在空间中的角色:admin,editor,viewer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (space_id, user_id),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_role (role)
);

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
ON DUPLICATE KEY UPDATE 
display_name = VALUES(display_name),
description = VALUES(description),
resource = VALUES(resource),
action = VALUES(action);

-- 插入基础角色数据
INSERT INTO roles (name, display_name, description) VALUES
('super_admin', '超级管理员', '拥有系统最高权限'),
('corp_admin', '企业管理员', '负责日常运维、空间管理、用户权限分配'),
('space_admin', '空间管理员', '在特定知识空间内拥有完全控制权'),
('content_reviewer', '内容审核员', '负责审阅和发布重要文档'),
('content_editor', '内容编辑者', '可创建、编辑、删除本空间内的文档内容'),
('content_viewer', '只读用户', '仅能查看、应用知识内容，不能修改')
ON DUPLICATE KEY UPDATE 
display_name = VALUES(display_name),
description = VALUES(description);

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'super_admin'
ON DUPLICATE KEY UPDATE role_id = VALUES(role_id);

-- 创建超级管理员用户（密码：admin123）
INSERT INTO users (username, phone, email, password, nickname, department, company, status) VALUES
('admin', '13800000000', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '超级管理员', 'IT部门', '示例公司', 1)
ON DUPLICATE KEY UPDATE 
phone = VALUES(phone),
email = VALUES(email),
password = VALUES(password),
nickname = VALUES(nickname),
department = VALUES(department),
company = VALUES(company);

-- 为超级管理员用户分配角色
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin' AND r.name = 'super_admin'
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);

-- 创建示例知识空间
INSERT INTO spaces (name, description, type, created_by) VALUES
('技术文档空间', '存放技术相关文档', 'department', (SELECT id FROM users WHERE username = 'admin')),
('项目文档空间', '存放项目相关文档', 'project', (SELECT id FROM users WHERE username = 'admin')),
('团队协作空间', '团队协作相关文档', 'team', (SELECT id FROM users WHERE username = 'admin'))
ON DUPLICATE KEY UPDATE 
description = VALUES(description),
type = VALUES(type);

-- 将超级管理员添加到所有空间
INSERT INTO space_members (space_id, user_id, role)
SELECT s.id, u.id, 'admin'
FROM spaces s, users u
WHERE u.username = 'admin'
ON DUPLICATE KEY UPDATE role = 'admin';

-- 显示初始化结果
SELECT 'Database initialization completed successfully!' as message;
SELECT COUNT(*) as user_count FROM users;
SELECT COUNT(*) as role_count FROM roles;
SELECT COUNT(*) as permission_count FROM permissions;
SELECT COUNT(*) as space_count FROM spaces;
