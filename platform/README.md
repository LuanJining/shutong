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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kbase
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
# 启动IAM服务
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

```bash
# 构建镜像
docker build -t kbase-platform .

# 运行容器
docker run -p 8080:8080 kbase-platform
```

### Kubernetes部署

```bash
# 应用Kubernetes配置
kubectl apply -f deployment/k8s/
```

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
