# Knowledge Base Platform

企业级知识库管理平台，基于Gin和PostgreSQL构建。

## 项目结构

```
platform/
├── cmd/                    # 应用程序入口
│   ├── iam/               # IAM服务
│   ├── workflow/          # 工作流服务
│   └── kbservice/         # 知识库服务
├── internal/              # 内部包
│   ├── iam/              # 身份认证与授权
│   ├── workflow/         # 工作流管理
│   └── kbservice/        # 知识库服务
├── scripts/              # 脚本文件
│   ├── init-db.go        # 数据库初始化
│   └── test-iam.sh       # IAM测试脚本
├── deployment/           # 部署配置
│   ├── db-script/        # 数据库脚本
│   ├── docker/           # Docker配置
│   └── k8s/              # Kubernetes配置
├── env.example           # 环境变量示例
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

## 快速开始

### 1. 环境准备

确保已安装以下软件：
- Go 1.25+
- PostgreSQL 12+
- Make (可选)

### 2. 配置环境变量

复制环境变量示例文件：
```bash
cp env.example .env
```

编辑 `.env` 文件，配置数据库连接信息：
```env
# 环境配置
KBASE_ENV=localtest

# 数据库配置
KBASE_DATABASE_HOST=localhost
KBASE_DATABASE_PORT=5432
KBASE_DATABASE_USER=postgres
KBASE_DATABASE_PASSWORD=your_password
KBASE_DATABASE_DBNAME=kbase
```

或者直接使用环境变量：
```bash
export KBASE_ENV=localtest
export KBASE_DATABASE_PASSWORD=your_password
```

### 3. 安装依赖

```bash
make deps
```

### 4. 初始化数据库

```bash
make init-db
```

这将创建：
- 基础权限数据
- 6种角色（超级管理员、企业管理员、空间管理员、内容审核员、内容编辑者、只读用户）
- 超级管理员用户（用户名：admin，密码：admin123）

### 5. 启动服务

```bash
# 启动IAM服务（本地测试环境）
make run-iam-local

# 启动IAM服务（生产环境）
make run-iam-prod

# 或者直接运行
make run-iam
```

### 6. 测试API

```bash
# 运行API测试
make test-api
```

## 服务说明

### IAM服务 (Identity and Access Management)

提供完整的身份认证与授权功能：

- **认证接口**：
  - `POST /api/v1/auth/login` - 用户登录（支持用户名/手机号/邮箱）
  - `POST /api/v1/auth/logout` - 用户登出
  - `POST /api/v1/auth/refresh` - 刷新Token
  - `PATCH /api/v1/auth/change-password` - 修改密码

- **用户管理**：
  - `GET /api/v1/users` - 获取用户列表
  - `GET /api/v1/users/:id` - 获取用户详情
  - `PUT /api/v1/users/:id` - 更新用户
  - `POST /api/v1/users` - 创建用户（仅超级管理员）
  - `DELETE /api/v1/users/:id` - 删除用户（仅超级管理员）

- **角色管理**：
  - `GET /api/v1/roles` - 获取角色列表
  - `GET /api/v1/roles/:id` - 获取角色详情
  - `POST /api/v1/roles` - 创建角色
  - `PUT /api/v1/roles/:id` - 更新角色
  - `DELETE /api/v1/roles/:id` - 删除角色

- **权限管理**：
  - `GET /api/v1/permissions` - 获取权限列表
  - `GET /api/v1/permissions/:id` - 获取权限详情

### 权限模型

系统采用RBAC（基于角色的访问控制）模型，支持以下角色：

| 角色 | 权限范围 |
|------|----------|
| 超级管理员 | 所有权限，包括用户管理 |
| 企业管理员 | 除用户管理外的所有权限 |
| 空间管理员 | 特定空间内的完全控制权 |
| 内容审核员 | 审阅和发布文档 |
| 内容编辑者 | 创建、编辑、删除文档 |
| 只读用户 | 仅查看内容 |

## API使用示例

### 1. 登录获取Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }'
```

支持用户名、手机号、邮箱登录。

### 2. 使用Token访问受保护的接口

```bash
# 获取用户列表
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <your-token>"

# 创建用户（需要超级管理员权限）
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "phone": "13800138000",
    "email": "newuser@example.com",
    "password": "123456",
    "nickname": "新用户",
    "department": "技术部",
    "company": "示例公司"
  }'

# 修改密码
curl -X PATCH http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldpassword",
    "new_password": "newpassword123"
  }'
```

## 配置管理

系统支持多种配置方式，优先级从高到低：

1. **环境变量**：`KBASE_*` 或传统环境变量
2. **YAML配置文件**：根据 `KBASE_ENV` 环境变量选择
3. **默认值**：内置默认配置

### 配置文件

- `internal/iam/config/localtest.yaml` - 本地测试环境配置
- `internal/iam/config/production.yaml` - 生产环境配置

### 环境变量

支持两种环境变量格式：

```bash
# 新格式（推荐）
export KBASE_ENV=production
export KBASE_DATABASE_HOST=localhost
export KBASE_DATABASE_PASSWORD=secret

# 传统格式（兼容）
export DB_HOST=localhost
export DB_PASSWORD=secret
```

## 开发命令

```bash
# 构建项目
make build

# 运行测试
make test

# 格式化代码
make fmt

# 代码检查
make vet

# 清理构建文件
make clean
```

## 部署说明

### Docker部署

使用Docker Compose一键部署所有服务：

```bash
# 一键部署
make docker-deploy

# 或者手动部署
cd deployment/docker
./deploy.sh
```

部署的服务：
- **PostgreSQL**: 端口5432
- **IAM服务**: 端口8080
- **KBService**: 端口8081  
- **Workflow**: 端口8082
- **Nginx**: 端口80（反向代理）

管理命令：
```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
make docker-stop

# 重启服务
docker-compose restart
```

### Kubernetes部署

使用Kubernetes部署到集群：

```bash
# 一键部署
make k8s-deploy

# 或者手动部署
cd deployment/k8s
./deploy.sh
```

部署的资源：
- **Namespace**: kb-platform
- **StatefulSet**: PostgreSQL数据库
- **Deployment**: IAM服务（2个副本）
- **Service**: 内部服务发现
- **ConfigMap**: 配置文件
- **Secret**: 敏感信息

管理命令：
```bash
# 查看Pod状态
kubectl get pods -n kb-platform

# 查看服务
kubectl get services -n kb-platform

# 查看日志
kubectl logs -f deployment/iam-deployment -n kb-platform

# 停止部署
make k8s-stop
```

### 环境要求

**Docker部署**：
- Docker 20.10+
- Docker Compose 2.0+
- 至少2GB可用内存

**Kubernetes部署**：
- Kubernetes 1.20+
- kubectl 1.20+
- 至少4GB可用内存
- 支持LoadBalancer或NodePort

## 安全注意事项

1. **JWT Secret**: 生产环境必须使用强密钥
2. **密码加密**: 使用bcrypt进行密码哈希
3. **Token过期**: 设置合理的token过期时间
4. **HTTPS**: 生产环境必须使用HTTPS
5. **输入验证**: 所有输入都经过验证和清理

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
