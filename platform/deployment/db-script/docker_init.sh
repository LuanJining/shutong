#!/bin/bash

# docker版本数据库初始化脚本

echo "🗄️ 开始初始化 PostgreSQL master 数据库..."

echo "创建数据库kb-platform..."
docker exec -i platform-postgres psql -U postgres -c "CREATE DATABASE kb_platform;"

echo "清理kb-platform数据库..."

# 1. 清理数据库（如果在Docker容器中）
docker exec -i platform-postgres psql -U postgres -d kb_platform -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"

echo "重新执行初始化脚本..."

# 2. 重新执行初始化脚本（使用重定向）
docker exec -i platform-postgres psql -U postgres -d kb_platform < ./init-database.sql

echo "数据库初始化完成！"