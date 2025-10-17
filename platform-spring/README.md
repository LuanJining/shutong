# Knowledge Base Platform - Spring Boot Edition

## 📋 项目简介

这是一个基于 Spring Boot WebFlux 的企业级知识库管理平台，从 Go 微服务架构迁移而来，整合了文档管理、工作流审批、向量搜索和 AI 问答等功能。

## 🏗️ 技术架构

### 核心技术栈
- **Spring Boot 3.5.6** - 应用框架
- **Spring WebFlux** - 响应式 Web 框架
- **Spring Data R2DBC** - 响应式数据库访问
- **Spring Security** - 安全认证与授权
- **PostgreSQL** - 关系型数据库
- **R2DBC PostgreSQL** - 响应式数据库驱动

### 集成服务
- **MinIO** - 对象存储服务
- **Qdrant** - 向量数据库
- **OpenAI API** - 大语言模型与 Embedding
- **PaddleOCR** - 文档 OCR 识别
- **JWT** - 无状态身份认证

## 📁 项目结构

```
platform-spring/
├── src/main/java/com/knowledgebase/platformspring/
│   ├── client/              # 外部服务客户端
│   │   ├── MinioClientService.java
│   │   ├── QdrantClientService.java
│   │   ├── OpenAIClientService.java
│   │   └── PaddleOCRClientService.java
│   ├── config/              # 配置类
│   ├── controller/          # REST API 控制器
│   │   ├── AuthController.java
│   │   ├── DocumentController.java
│   │   ├── SpaceController.java
│   │   └── WorkflowController.java
│   ├── dto/                 # 数据传输对象
│   ├── exception/           # 异常处理
│   ├── model/              # 实体模型（14个实体）
│   ├── repository/         # 数据访问层（14个Repository）
│   ├── security/           # 安全配置
│   └── service/            # 业务逻辑层
│       ├── AuthService.java
│       ├── DocumentService.java
│       ├── SpaceService.java
│       └── WorkflowService.java
├── src/main/resources/
│   └── application.yaml    # 应用配置
├── pom.xml                 # Maven 依赖配置
└── README.md
```

## 🚀 快速开始

### 前置要求

1. **Java 17+**
2. **Maven 3.8+**
3. **PostgreSQL 12+**
4. **MinIO** (可选，文档存储)
5. **Qdrant** (可选，向量搜索)
6. **PaddleOCR** (可选，文档 OCR)

### 环境搭建

#### 1. 启动 PostgreSQL

```bash
# 使用 Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=knowledge_base \
  -p 5432:5432 \
  postgres:15
```

#### 2. 初始化数据库

使用 Go 项目中的数据库初始化脚本：

```bash
psql -h localhost -U postgres -d knowledge_base \
  -f ../platform/deployment/database/init-database.sql
```

#### 3. 启动 MinIO

```bash
docker run -d \
  --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

#### 4. 启动 Qdrant

```bash
docker run -d \
  --name qdrant \
  -p 6333:6333 \
  qdrant/qdrant
```

### 配置应用

编辑 `src/main/resources/application.yaml`：

```yaml
spring:
  r2dbc:
    url: r2dbc:postgresql://localhost:5432/knowledge_base
    username: postgres
    password: postgres
  
  security:
    jwt:
      secret: your-secret-key-at-least-256-bits-long
      expiration: 86400000
      refresh-expiration: 604800000

app:
  openai:
    api-key: your-openai-api-key
    base-url: https://api.openai.com/v1
    model: gpt-4
    embedding-model: text-embedding-3-small
```

### 启动应用

```bash
# 方式 1: 使用 Maven
./mvnw spring-boot:run

# 方式 2: 使用启动脚本
chmod +x start.sh
./start.sh

# 方式 3: 打包后运行
./mvnw clean package
java -jar target/platform-spring-0.0.1-SNAPSHOT.jar
```

应用将在 `http://localhost:8080` 启动。

## 📚 API 文档

### Swagger UI (推荐)

启动应用后，访问 Swagger UI 查看完整的交互式 API 文档：

```
http://localhost:8080/swagger-ui.html
```

或访问 OpenAPI JSON：

```
http://localhost:8080/v3/api-docs
```

### API 端点示例

### 认证 API

#### 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "phone": "13800138000",
  "email": "test@example.com",
  "password": "password123",
  "nickname": "测试用户",
  "department": "技术部",
  "company": "示例公司"
}
```

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户"
    }
  }
}
```

#### 刷新 Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 获取当前用户
```http
GET /api/v1/auth/me
Authorization: Bearer {accessToken}
```

### 空间管理 API

#### 获取所有空间
```http
GET /api/v1/spaces
Authorization: Bearer {accessToken}
```

#### 创建空间
```http
POST /api/v1/spaces
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "name": "技术文档",
  "description": "技术相关文档空间",
  "type": "department"
}
```

#### 创建子空间
```http
POST /api/v1/spaces/sub-spaces
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "spaceId": 1,
  "name": "后端开发",
  "description": "后端开发文档"
}
```

### 文档管理 API

#### 上传文档
```http
POST /api/v1/documents/upload
Authorization: Bearer {accessToken}
Content-Type: multipart/form-data

file: [文件]
spaceId: 1
subSpaceId: 1
classId: 1
```

#### 获取空间文档
```http
GET /api/v1/documents/space/{spaceId}
Authorization: Bearer {accessToken}
```

#### 文档问答
```http
POST /api/v1/documents/chat
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "question": "如何部署应用？",
  "spaceId": 1
}
```

#### 发布文档
```http
POST /api/v1/documents/{id}/publish
Authorization: Bearer {accessToken}
```

### 工作流 API

#### 获取我的待办任务
```http
GET /api/v1/workflows/tasks/my
Authorization: Bearer {accessToken}
```

#### 审批任务
```http
POST /api/v1/workflows/tasks/{id}/approve
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "comment": "审批通过",
  "approved": true
}
```

## 🔐 安全说明

### JWT 认证

所有 API（除了登录/注册）都需要在 Header 中携带 JWT Token：

```
Authorization: Bearer {accessToken}
```

Token 过期时间：
- Access Token: 24 小时
- Refresh Token: 7 天

### 密码加密

用户密码使用 BCrypt 算法加密存储。

## 🧪 测试

```bash
# 运行所有测试
./mvnw test

# 运行特定测试类
./mvnw test -Dtest=AuthServiceTest
```

## 📊 监控

应用暴露了 Actuator 端点用于监控：

```
http://localhost:8080/actuator/health
http://localhost:8080/actuator/info
http://localhost:8080/actuator/metrics
```

## 🔧 配置说明

### 数据库连接池配置

```yaml
spring:
  r2dbc:
    pool:
      initial-size: 10
      max-size: 50
      max-idle-time: 30m
```

### 日志配置

```yaml
logging:
  level:
    root: INFO
    com.knowledgebase: DEBUG
  file:
    name: logs/application.log
```

## 🚨 故障排查

### 应用无法启动

1. 检查 PostgreSQL 是否运行
2. 检查数据库连接配置
3. 检查端口 8080 是否被占用

### 文档上传失败

1. 检查 MinIO 是否运行
2. 检查 MinIO 配置（endpoint、credentials）
3. 检查文件大小限制

### 向量搜索不工作

1. 检查 Qdrant 是否运行
2. 检查 OpenAI API Key 是否配置
3. 查看日志中的错误信息

## 📈 性能优化

### 响应式编程最佳实践

1. 避免在响应式链中使用阻塞操作
2. 使用 `subscribeOn()` 和 `publishOn()` 控制执行线程
3. 合理使用背压策略

### 数据库优化

1. 添加合适的索引
2. 使用连接池
3. 避免 N+1 查询问题

## 🔄 与 Go 版本的差异

| 特性 | Go 版本 | Spring Boot 版本 |
|------|---------|-----------------|
| 架构 | 4个微服务 | 单体应用 |
| 并发模型 | Goroutine | Reactor (Mono/Flux) |
| ORM | GORM | Spring Data R2DBC |
| 路由 | Gin | Spring WebFlux |
| 配置 | Viper + YAML | Spring Boot + YAML |

## 📝 待完善功能

- [ ] 添加 Swagger API 文档
- [ ] 实现完整的权限控制（RBAC）
- [ ] 添加文件流式上传/下载
- [ ] 实现文档全文搜索
- [ ] 添加单元测试和集成测试
- [ ] 实现 WebSocket 实时通知
- [ ] 添加缓存支持（Redis）
- [ ] 实现分布式追踪
- [ ] 添加性能指标监控

## 📧 联系方式

如有问题，请提交 Issue 或联系开发团队。

## 📄 许可证

[Your License Here]

