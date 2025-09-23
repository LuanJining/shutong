#!/bin/bash

# Docker部署脚本

set -e

echo "🚀 开始使用Docker Compose部署知识库平台..."

# 检查Docker和Docker Compose是否可用
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 进入部署目录
cd "$(dirname "$0")"

# 停止现有服务
echo "🛑 停止现有服务..."
docker-compose down

# 构建镜像
echo "🔨 构建Docker镜像..."
docker-compose build

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 检查服务状态
echo "📋 服务状态："
docker-compose ps

# 等待数据库就绪
echo "⏳ 等待数据库就绪..."
until docker exec platform-postgres pg_isready -U postgres; do
    echo "等待PostgreSQL启动..."
    sleep 2
done

# 初始化数据库
echo "🗄️ 初始化数据库..."
if [ -f "../db-script/init-database.sql" ]; then
    docker exec -i platform-postgres psql -U postgres -d kb-platform < ../db-script/init-database.sql
    echo "✅ 数据库初始化完成"
else
    echo "⚠️ 数据库初始化脚本未找到，请手动初始化"
fi

echo "✅ 部署完成！"
echo ""
echo "📋 服务状态："
docker-compose ps

echo ""
echo "🌐 访问地址："
echo "  IAM服务: http://localhost:8080"
echo "  KBService: http://localhost:8081"
echo "  Workflow: http://localhost:8082"
echo "  Nginx: http://localhost"

echo ""
echo "🔧 管理命令："
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo "  进入容器: docker exec -it platform-iam /bin/sh"
