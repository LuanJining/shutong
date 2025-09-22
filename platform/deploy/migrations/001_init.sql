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


-- 初始化管理员账号，默认密码admin123
INSERT INTO iam_users (
    id, 
    name, 
    phone, 
    password_hash, 
    roles, 
    spaces, 
    created_at, 
    updated_at
    ) VALUES ('1', 'admin', '13800138000', '$2a$10$FDX66S5GcX1RSRsd/AVqpea0.QwX4OXIM5F66O45ncUBtnftXHI5G', ARRAY['admin'], ARRAY['default'], NOW(), NOW());