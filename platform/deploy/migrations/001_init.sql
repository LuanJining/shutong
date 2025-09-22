CREATE TABLE IF NOT EXISTS iam_users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    roles TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    spaces TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS iam_roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    permissions TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS iam_spaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS iam_policies (
    id TEXT PRIMARY KEY,
    space_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    resource TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_policy_space FOREIGN KEY (space_id) REFERENCES iam_spaces(id),
    CONSTRAINT fk_policy_role FOREIGN KEY (role_id) REFERENCES iam_roles(id)
);

CREATE TABLE IF NOT EXISTS wf_definitions (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    nodes JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS wf_instances (
    id TEXT PRIMARY KEY,
    definition_id TEXT NOT NULL,
    business_id TEXT NOT NULL,
    space_id TEXT NOT NULL,
    status TEXT NOT NULL,
    current_node_id TEXT,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    CONSTRAINT fk_instance_definition FOREIGN KEY (definition_id) REFERENCES wf_definitions(id)
);

CREATE TABLE IF NOT EXISTS wf_history (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL,
    node_id TEXT,
    actor_id TEXT NOT NULL,
    action TEXT NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_history_instance FOREIGN KEY (instance_id) REFERENCES wf_instances(id)
);

CREATE INDEX IF NOT EXISTS idx_iam_policies_space ON iam_policies(space_id);
CREATE INDEX IF NOT EXISTS idx_wf_instances_definition ON wf_instances(definition_id);
CREATE INDEX IF NOT EXISTS idx_wf_history_instance ON wf_history(instance_id);


-- 初始化角色数据
INSERT INTO iam_roles (id, name, description, permissions, created_at, updated_at) VALUES 
('role_1', 'super_admin', '超级管理员', ARRAY['*'], NOW(), NOW()),
('role_2', 'enterprise_admin', '企业管理员', ARRAY['manage_spaces', 'manage_users', 'manage_roles', 'export_data'], NOW(), NOW()),
('role_3', 'space_admin', '空间管理员', ARRAY['manage_space_content', 'manage_space_members', 'configure_approval'], NOW(), NOW()),
('role_4', 'content_reviewer', '内容审核员', ARRAY['review_content', 'approve_content', 'export_data'], NOW(), NOW()),
('role_5', 'content_editor', '内容编辑者', ARRAY['create_content', 'edit_content', 'delete_content'], NOW(), NOW()),
('role_6', 'readonly_user', '只读用户', ARRAY['view_content', 'export_data'], NOW(), NOW());

-- 初始化空间数据
INSERT INTO iam_spaces (id, name, description, created_at, updated_at) VALUES 
('space_1', 'default', '默认知识空间', NOW(), NOW()),
('space_2', 'public', '公共知识空间', NOW(), NOW());

-- 初始化权限策略数据
INSERT INTO iam_policies (id, space_id, role_id, resource, action, created_at) VALUES 
-- 超级管理员权限
('policy_1', 'space_1', 'role_1', '*', '*', NOW()),
('policy_2', 'space_2', 'role_1', '*', '*', NOW()),

-- 企业管理员权限
('policy_3', 'space_1', 'role_2', 'spaces', 'manage', NOW()),
('policy_4', 'space_1', 'role_2', 'users', 'manage', NOW()),
('policy_5', 'space_1', 'role_2', 'content', 'manage', NOW()),
('policy_6', 'space_2', 'role_2', 'spaces', 'manage', NOW()),
('policy_7', 'space_2', 'role_2', 'users', 'manage', NOW()),
('policy_8', 'space_2', 'role_2', 'content', 'manage', NOW()),

-- 内容审核员权限
('policy_9', 'space_1', 'role_4', 'content', 'review', NOW()),
('policy_10', 'space_1', 'role_4', 'content', 'approve', NOW()),
('policy_11', 'space_2', 'role_4', 'content', 'review', NOW()),
('policy_12', 'space_2', 'role_4', 'content', 'approve', NOW()),

-- 内容编辑者权限
('policy_13', 'space_1', 'role_5', 'content', 'create', NOW()),
('policy_14', 'space_1', 'role_5', 'content', 'edit', NOW()),
('policy_15', 'space_1', 'role_5', 'content', 'delete', NOW()),
('policy_16', 'space_2', 'role_5', 'content', 'create', NOW()),
('policy_17', 'space_2', 'role_5', 'content', 'edit', NOW()),
('policy_18', 'space_2', 'role_5', 'content', 'delete', NOW()),

-- 只读用户权限
('policy_19', 'space_1', 'role_6', 'content', 'view', NOW()),
('policy_20', 'space_2', 'role_6', 'content', 'view', NOW());

-- 初始化工作流定义
INSERT INTO wf_definitions (id, code, name, description, nodes, created_at, updated_at) VALUES 
('wf_1', 'document_approval', '文档审批流程', '知识库文档发布审批流程', 
 '{"nodes": [{"id": "start", "type": "start", "name": "开始"}, {"id": "review", "type": "user_task", "name": "内容审核", "assignee_role": "content_reviewer"}, {"id": "approve", "type": "user_task", "name": "审批确认", "assignee_role": "space_admin"}, {"id": "end", "type": "end", "name": "结束"}]}',
 NOW(), NOW());

-- 初始化超级管理员账号，默认密码admin123
INSERT INTO iam_users (
    id, 
    name, 
    phone, 
    password_hash, 
    roles, 
    spaces, 
    created_at, 
    updated_at
    ) VALUES ('1', 'admin', '13800138000', '$2a$10$FDX66S5GcX1RSRsd/AVqpea0.QwX4OXIM5F66O45ncUBtnftXHI5G', ARRAY['super_admin'], ARRAY['space_1', 'space_2'], NOW(), NOW());