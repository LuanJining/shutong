## 启动所有服务 命令：make start
start:
	@echo "启动所有服务..."
	@nohup go run cmd/gateway/main.go > gateway.log 2>&1 & echo $$! > gateway.pid
	@nohup go run cmd/iam/main.go > iam.log 2>&1 & echo $$! > iam.pid
	@nohup go run cmd/workflow/main.go > workflow.log 2>&1 & echo $$! > workflow.pid
	@nohup go run cmd/kb_service/main.go > kb_service.log 2>&1 & echo $$! > kb_service.pid
	@echo "服务已启动，查看日志: make logs"

## 停止所有服务 命令：make stop
stop:
	@echo "停止所有服务..."
	@if [ -f gateway.pid ]; then kill -9 $$(cat gateway.pid) 2>/dev/null && rm gateway.pid; fi
	@if [ -f iam.pid ]; then kill -9 $$(cat iam.pid) 2>/dev/null && rm iam.pid; fi
	@if [ -f workflow.pid ]; then kill -9 $$(cat workflow.pid) 2>/dev/null && rm workflow.pid; fi
	@if [ -f kb_service.pid ]; then kill -9 $$(cat kb_service.pid) 2>/dev/null && rm kb_service.pid; fi
	@echo "服务已停止"

## 查看服务日志 命令：make logs
logs:
	@echo "=== Gateway 日志 ==="
	@tail -f gateway.log &
	@echo "=== IAM 日志 ==="
	@tail -f iam.log &
	@echo "=== Workflow 日志 ==="
	@tail -f workflow.log &
	@echo "=== KB Service 日志 ==="
	@tail -f kb_service.log

## 构建所有服务 命令：make build
build:
	go build -o bin/gateway cmd/gateway/main.go
	go build -o bin/iam cmd/iam/main.go
	go build -o bin/workflow cmd/workflow/main.go
	go build -o bin/kb_service cmd/kb_service/main.go

## 清理构建文件 命令：make clean
clean:
	rm -f bin/gateway bin/iam bin/workflow bin/kb_service *.log *.pid

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
	@if [ -f workflow.pid ]; then \
		if ps -p $$(cat workflow.pid) > /dev/null 2>&1; then \
			echo "✅ Workflow: 运行中 (PID: $$(cat workflow.pid))"; \
		else \
			echo "❌ Workflow: 已停止"; \
		fi; \
	else \
		echo "❌ Workflow: 未启动"; \
	fi
	@if [ -f kb_service.pid ]; then \
		if ps -p $$(cat kb_service.pid) > /dev/null 2>&1; then \
			echo "✅ KB Service: 运行中 (PID: $$(cat kb_service.pid))"; \
		else \
			echo "❌ KB Service: 已停止"; \
		fi; \
	else \
		echo "❌ KB Service: 未启动"; \
	fi

## 完整测试流程 命令：make test-all
test-all: start
	@echo "等待服务启动..."
	@sleep 3
	@make test-gateway
	@make stop
