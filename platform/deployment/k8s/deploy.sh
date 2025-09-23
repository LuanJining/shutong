#!/bin/bash

# Kubernetes部署脚本

set -e

echo "🚀 开始部署知识库平台到Kubernetes..."

# 检查kubectl是否可用
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl未安装，请先安装kubectl"
    exit 1
fi

# 创建命名空间
echo "📦 创建命名空间..."
kubectl create namespace kb-platform --dry-run=client -o yaml | kubectl apply -f -

# 创建Secret
echo "🔐 创建Secret..."
kubectl create secret generic postgres-secret \
  --from-literal=password=postgres123 \
  --namespace=kb-platform \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic iam-secret \
  --from-literal=jwt-secret=your-secret-key-change-this-in-production \
  --namespace=kb-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# 创建PostgreSQL初始化脚本ConfigMap
echo "📝 创建PostgreSQL初始化脚本..."
kubectl create configmap postgres-init-script \
  --from-file=init-database.sql=../db-script/init-database.sql \
  --namespace=kb-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# 部署PostgreSQL
echo "🗄️ 部署PostgreSQL..."
kubectl apply -f statefulset/postgres-stateful-replica.yaml

# 等待PostgreSQL就绪
echo "⏳ 等待PostgreSQL就绪..."
kubectl wait --for=condition=ready pod -l app=postgres-master -n kb-platform --timeout=300s

# 部署IAM服务
echo "🔐 部署IAM服务..."
kubectl apply -f configmap/iam-configmap.yaml
kubectl apply -f deployment/iam-deployment.yaml

# 创建Service
echo "🌐 创建Service..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: postgres-master-service
  namespace: kb-platform
spec:
  selector:
    app: postgres-master
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: iam-service
  namespace: kb-platform
spec:
  selector:
    app: iam
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: iam-ingress
  namespace: kb-platform
spec:
  selector:
    app: iam
  ports:
  - port: 8080
    targetPort: 8080
  type: LoadBalancer
EOF

# 等待IAM服务就绪
echo "⏳ 等待IAM服务就绪..."
kubectl wait --for=condition=available deployment/iam-deployment -n kb-platform --timeout=300s

# 初始化数据库
echo "🗄️ 初始化数据库..."
kubectl exec -i deployment/iam-deployment -n kb-platform -- /bin/bash -c "
  export KBASE_DATABASE_HOST=postgres-master-service
  export KBASE_DATABASE_PORT=5432
  export KBASE_DATABASE_USER=postgres
  export KBASE_DATABASE_PASSWORD=postgres123
  export KBASE_DATABASE_DBNAME=kb-platform
  export KBASE_DATABASE_SSLMODE=disable
  export KBASE_JWT_SECRET=your-secret-key-change-this-in-production
  export KBASE_JWT_EXPIRE_TIME=24
  go run scripts/init-db.go
"

echo "✅ 部署完成！"
echo ""
echo "📋 服务状态："
kubectl get pods -n kb-platform
kubectl get services -n kb-platform

echo ""
echo "🌐 访问地址："
kubectl get service iam-ingress -n kb-platform -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
echo ""

echo "🔧 管理命令："
echo "  查看日志: kubectl logs -f deployment/iam-deployment -n kb-platform"
echo "  进入Pod: kubectl exec -it deployment/iam-deployment -n kb-platform -- /bin/bash"
echo "  删除部署: kubectl delete namespace kb-platform"
