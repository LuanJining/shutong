# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目信息
BINARY_NAME=kbase
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/iam/main.go
INIT_DB_PATH=./scripts/init-db.go

# 构建
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# 构建Linux版本
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# 清理
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# 运行IAM服务
.PHONY: run-iam
run-iam:
	$(GOCMD) run $(MAIN_PATH)

# 运行IAM服务（本地测试环境）
.PHONY: run-iam-local
run-iam-local:
	KBASE_ENV=localtest $(GOCMD) run $(MAIN_PATH)

# 运行IAM服务（生产环境）
.PHONY: run-iam-prod
run-iam-prod:
	KBASE_ENV=production $(GOCMD) run $(MAIN_PATH)

# 初始化数据库
.PHONY: init-db
init-db:
	$(GOCMD) run $(INIT_DB_PATH)

# 测试
.PHONY: test
test:
	$(GOTEST) -v ./...

# 依赖管理
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 安装依赖
.PHONY: install
install:
	$(GOGET) -d -v ./...

# 格式化代码
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: vet
vet:
	$(GOCMD) vet ./...

# 运行测试脚本
.PHONY: test-api
test-api:
	./scripts/test-iam.sh

# Docker部署
.PHONY: docker-deploy
docker-deploy:
	cd deployment/docker && ./deploy.sh

# Docker停止
.PHONY: docker-stop
docker-stop:
	cd deployment/docker && docker-compose down

# Kubernetes部署
.PHONY: k8s-deploy
k8s-deploy:
	cd deployment/k8s && ./deploy.sh

# Kubernetes停止
.PHONY: k8s-stop
k8s-stop:
	kubectl delete namespace kb-platform

# 帮助
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build       - 构建项目"
	@echo "  build-linux - 构建Linux版本"
	@echo "  clean       - 清理构建文件"
	@echo "  run-iam     - 运行IAM服务"
	@echo "  run-iam-local - 运行IAM服务（本地测试环境）"
	@echo "  run-iam-prod  - 运行IAM服务（生产环境）"
	@echo "  init-db     - 初始化数据库"
	@echo "  test        - 运行测试"
	@echo "  test-api    - 运行API测试"
	@echo "  docker-deploy - Docker部署"
	@echo "  docker-stop   - Docker停止"
	@echo "  k8s-deploy    - Kubernetes部署"
	@echo "  k8s-stop      - Kubernetes停止"
	@echo "  deps        - 下载依赖"
	@echo "  install     - 安装依赖"
	@echo "  fmt         - 格式化代码"
	@echo "  vet         - 代码检查"
	@echo "  help        - 显示帮助"
