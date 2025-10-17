# Spring Boot Knowledge Base Platform - Migration Guide

## 项目概述

本项目是从 Go 微服务架构迁移到 Spring Boot WebFlux 单体应用的知识库管理平台。

## 技术栈

- **Spring Boot 3.5.6** - 主框架
- **Spring WebFlux** - 响应式 Web 框架
- **Spring Data R2DBC** - 响应式数据库访问
- **PostgreSQL** - 关系型数据库
- **MinIO** - 对象存储
- **Qdrant** - 向量数据库
- **OpenAI API** - LLM 和 Embedding
- **JWT** - 身份认证
- **Lombok** - 简化代码

## 项目结构

```
src/main/java/com/knowledgebase/platformspring/
├── client/                 # 外部服务客户端
│   ├── MinioClientService.java
│   ├── QdrantClientService.java
│   ├── OpenAIClientService.java
│   └── PaddleOCRClientService.java
├── config/                 # 配置类
│   ├── JwtConfig.java
│   ├── MinioConfig.java
│   ├── OpenAIConfig.java
│   ├── QdrantConfig.java
│   └── PaddleOCRConfig.java
├── controller/            # REST API 控制器
│   └── AuthController.java
├── dto/                   # 数据传输对象
│   ├── ApiResponse.java
│   ├── LoginRequest.java
│   ├── LoginResponse.java
│   └── RegisterRequest.java
├── exception/             # 异常处理
│   ├── BusinessException.java
│   └── GlobalExceptionHandler.java
├── model/                 # 实体模型
│   ├── User.java
│   ├── Role.java
│   ├── Permission.java
│   ├── Space.java
│   ├── Document.java
│   ├── Workflow.java
│   └── ...
├── repository/            # 数据访问层
│   ├── UserRepository.java
│   ├── DocumentRepository.java
│   └── ...
├── security/              # 安全配置
│   ├── SecurityConfig.java
│   ├── JwtUtil.java
│   └── JwtAuthenticationFilter.java
├── service/               # 业务逻辑层
│   └── AuthService.java
└── PlatformSpringApplication.java

src/main/resources/
└── application.yaml       # 应用配置
```

## 配置说明

### application.yaml

主要配置项：

```yaml
spring:
  r2dbc:
    url: r2dbc:postgresql://localhost:5432/knowledge_base
    username: postgres
    password: postgres
  
  security:
    jwt:
      secret: your-secret-key
      expiration: 86400000  # 24小时
      refresh-expiration: 604800000  # 7天

app:
  minio:
    endpoint: http://localhost:9000
    access-key: minioadmin
    secret-key: minioadmin
    bucket-name: knowledge-base
  
  openai:
    api-key: your-openai-api-key
    base-url: https://api.openai.com/v1
    model: gpt-4
    embedding-model: text-embedding-3-small
  
  qdrant:
    host: localhost
    port: 6333
    collection-name: knowledge_base
    vector-size: 1536
```

## 核心功能

### 1. 用户认证 (AuthService)

- 用户登录
- 用户注册
- Token 刷新
- 获取当前用户信息

### 2. 文档管理 (待实现)

- 文档上传
- 文档解析（OCR）
- 文档向量化
- 文档搜索
- 文档聊天

### 3. 工作流管理 (待实现)

- 创建审批流
- 任务审批
- 流程查询

### 4. 空间管理 (待实现)

- 空间创建
- 成员管理
- 权限控制

## API 端点

### 认证相关

```
POST /api/v1/auth/login          # 登录
POST /api/v1/auth/register       # 注册
POST /api/v1/auth/refresh        # 刷新token
GET  /api/v1/auth/me             # 获取当前用户
```

### 其他端点 (待实现)

参考 Go 项目的 API 设计实现。

## 数据库迁移

使用 Go 项目中的 `deployment/database/init-database.sql` 初始化数据库表结构。

## 启动步骤

1. 启动 PostgreSQL
2. 启动 MinIO
3. 启动 Qdrant
4. （可选）启动 PaddleOCR 服务
5. 配置 `application.yaml`
6. 运行应用：`./mvnw spring-boot:run`

## 待完成工作

- [ ] 实现 DocumentService 和 DocumentController
- [ ] 实现 WorkflowService 和 WorkflowController
- [ ] 实现 SpaceService 和 SpaceController
- [ ] 添加权限控制注解
- [ ] 添加单元测试
- [ ] 添加 API 文档（Swagger/OpenAPI）
- [ ] 实现文件流式处理
- [ ] 添加监控和日志

## 与 Go 版本的主要差异

1. **架构变化**：从多个微服务合并为单个应用
2. **框架**：从 Gin 迁移到 Spring WebFlux
3. **响应式编程**：全面使用 Reactor（Mono/Flux）
4. **数据库访问**：从 GORM 迁移到 Spring Data R2DBC
5. **依赖注入**：使用 Spring 的 DI 容器

## 注意事项

1. R2DBC 不支持 JPA 的关联查询，需要手动处理关联关系
2. 所有数据库操作都返回 Mono 或 Flux
3. 避免在响应式链中使用阻塞操作
4. JWT secret 应使用环境变量配置，不要硬编码

