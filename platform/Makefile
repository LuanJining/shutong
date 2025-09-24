## 启动所有服务 命令：make start
start:
	@echo "启动所有服务..."
	@nohup go run cmd/gateway/main.go > gateway.log 2>&1 & echo $$! > gateway.pid
	@nohup go run cmd/iam/main.go > iam.log 2>&1 & echo $$! > iam.pid
	@echo "服务已启动，查看日志: make logs"

## 停止所有服务 命令：make stop
stop:
	@echo "停止所有服务..."
	@if [ -f gateway.pid ]; then kill -9 $$(cat gateway.pid) 2>/dev/null && rm gateway.pid; fi
	@if [ -f iam.pid ]; then kill -9 $$(cat iam.pid) 2>/dev/null && rm iam.pid; fi
	@echo "服务已停止"

## 查看服务日志 命令：make logs
logs:
	@echo "=== Gateway 日志 ==="
	@tail -f gateway.log &
	@echo "=== IAM 日志 ==="
	@tail -f iam.log

## 测试Gateway服务 命令：make test-gateway
test-gateway:
	./scripts/test-gateway.sh

## 测试IAM服务 命令：make test-iam
test-iam:
	./scripts/test-iam-enhanced.sh

## 构建所有服务 命令：make build
build:
	go build -o bin/gateway cmd/gateway/main.go
	go build -o bin/iam cmd/iam/main.go

## 清理构建文件 命令：make clean
clean:
	rm -f bin/gateway bin/iam *.log *.pid

## 重启服务 命令：make restart
restart: stop start

## 检查服务状态 命令：make status
status:
	@echo "检查服务状态..."
	@if [ -f gateway.pid ]; then \
		if ps -p $$(cat gateway.pid) > /dev/null 2>&1; then \
			echo "✅ Gateway: 运行中 (PID: $$(cat gateway.pid))"; \
		else \
			echo "❌ Gateway: 已停止"; \
		fi; \
	else \
		echo "❌ Gateway: 未启动"; \
	fi
	@if [ -f iam.pid ]; then \
		if ps -p $$(cat iam.pid) > /dev/null 2>&1; then \
			echo "✅ IAM: 运行中 (PID: $$(cat iam.pid))"; \
		else \
			echo "❌ IAM: 已停止"; \
		fi; \
	else \
		echo "❌ IAM: 未启动"; \
	fi

## 完整测试流程 命令：make test-all
test-all: start
	@echo "等待服务启动..."
	@sleep 3
	@make test-gateway
	@make stop
