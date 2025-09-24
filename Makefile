## 直接本地测试 命令：make test
test:
	cd platform && go run cmd/gateway/main.go &
	cd platform && go run cmd/iam/main.go &

## 测试Gateway服务 命令：make test-gateway
test-gateway:
	cd platform && ./scripts/test-gateway.sh

## 测试IAM服务 命令：make test-iam
test-iam:
	cd platform && ./scripts/test-iam-enhanced.sh

## 构建所有服务 命令：make build
build:
	cd platform && go build -o bin/gateway cmd/gateway/main.go
	cd platform && go build -o bin/iam cmd/iam/main.go

## 清理构建文件 命令：make clean
clean:
	cd platform && rm -f bin/gateway bin/iam
