#!/bin/bash

set -e

NAMESPACE="kb-platform"

echo "========================================="
echo "部署 Spring Boot 应用到 Kubernetes"
echo "命名空间: ${NAMESPACE}"
echo "========================================="

# 创建命名空间（如果不存在）
echo ""
echo "=== 1. 创建命名空间 ==="
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

# 应用 ConfigMap
echo ""
echo "=== 2. 应用 ConfigMap ==="
kubectl apply -f platform-configmap.yaml

# 应用 Secret
echo ""
echo "=== 3. 应用 Secret ==="
kubectl apply -f platform-secret.yaml

# 应用 Deployment 和 Service
echo ""
echo "=== 4. 应用 Deployment 和 Service ==="
kubectl apply -f platform-deployment.yaml

# 等待部署完成
echo ""
echo "=== 5. 等待 Pod 启动 ==="
kubectl rollout status deployment/platform-spring -n ${NAMESPACE} --timeout=300s

# 查看 Pod 状态
echo ""
echo "=== 6. Pod 状态 ==="
kubectl get pods -n ${NAMESPACE} -l app=platform-spring

# 查看 Service
echo ""
echo "=== 7. Service 信息 ==="
kubectl get svc -n ${NAMESPACE} platform-spring-service

echo ""
echo "========================================="
echo "✅ 部署完成！"
echo "========================================="
echo ""
echo "查看日志："
echo "  kubectl logs -f -n ${NAMESPACE} -l app=platform-spring"
echo ""
echo "进入容器："
echo "  kubectl exec -it -n ${NAMESPACE} <pod-name> -- sh"
echo ""
echo "查看服务："
echo "  kubectl get all -n ${NAMESPACE} -l app=platform-spring"

