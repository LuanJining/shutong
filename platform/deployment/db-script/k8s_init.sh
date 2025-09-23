#!/bin/bash
# k8s版本数据库初始化脚本
echo "🗄️ 开始初始化 PostgreSQL master 数据库..."

# 获取 PostgreSQL master pod 名称
POSTGRES_POD=$(kubectl get pods -n kb-platform -l app=postgres-master -o jsonpath='{.items[0].metadata.name}')

if [ -z "$POSTGRES_POD" ]; then
    echo "❌ 未找到 PostgreSQL master pod，请确保 StatefulSet 已部署"
    exit 1
fi

echo "📋 找到 PostgreSQL master pod: $POSTGRES_POD"

# 等待 pod 就绪
echo "⏳ 等待 PostgreSQL master pod 就绪..."
kubectl wait --for=condition=ready pod/$POSTGRES_POD -n kb-platform --timeout=60s

# 创建数据库
echo "🗄️ 创建数据库..."
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "CREATE DATABASE kb-platform;"

# 创建用户
echo "🗄️ 创建用户..."
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "CREATE USER postgres WITH PASSWORD 'postgres123';"

# 清理并重新创建数据库结构
echo "🧹 清理数据库结构..."
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb-platform -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"

# 执行初始化脚本
echo "📝 执行数据库初始化脚本..."
kubectl exec -i $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb-platform < ./init-database.sql

echo "✅ PostgreSQL master 数据库初始化完成！"


kubectl get service postgres-master-service -n kb-platform -o jsonpath='{.spec.clusterIP}'