# Knowledge Base Platform - Spring Boot Edition

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
