#!/bin/bash
VERSION=$1
if [ -z "$VERSION" ]; then
    echo "请输入版本号"
    exit 1
fi

NAMESPACE="kb-platform"

echo "开始更新镜像版本到: ${VERSION}"
kubectl set image deployment/gateway gateway=harbor.kunxiangtech.com:8443/kb-platform/gateway:${VERSION} -n ${NAMESPACE}
kubectl set image deployment/iam iam=harbor.kunxiangtech.com:8443/kb-platform/iam:${VERSION} -n ${NAMESPACE}
kubectl set image deployment/workflow workflow=harbor.kunxiangtech.com:8443/kb-platform/workflow:${VERSION} -n ${NAMESPACE}
kubectl set image deployment/kb-service kb-service=harbor.kunxiangtech.com:8443/kb-platform/kb-service:${VERSION} -n ${NAMESPACE}
kubectl set image deployment/frontend frontend=harbor.kunxiangtech.com:8443/kb-platform/frontend:${VERSION} -n ${NAMESPACE}

echo ""
echo "等待滚动更新完成..."
kubectl rollout status deployment/gateway -n ${NAMESPACE}
kubectl rollout status deployment/iam -n ${NAMESPACE}
kubectl rollout status deployment/workflow -n ${NAMESPACE}
kubectl rollout status deployment/kb-service -n ${NAMESPACE}
kubectl rollout status deployment/frontend -n ${NAMESPACE}

echo ""
echo "当前部署的版本信息："
echo "===================="
CURRENT_VERSION=$(kubectl get deployment/gateway -n ${NAMESPACE} -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $3}')
echo "Gateway版本号: ${CURRENT_VERSION}"

CURRENT_VERSION=$(kubectl get deployment/iam -n ${NAMESPACE} -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $3}')
echo "IAM版本号: ${CURRENT_VERSION}"

CURRENT_VERSION=$(kubectl get deployment/workflow -n ${NAMESPACE} -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $3}')
echo "Workflow版本号: ${CURRENT_VERSION}"

CURRENT_VERSION=$(kubectl get deployment/kb-service -n ${NAMESPACE} -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $3}')
echo "KB Service版本号: ${CURRENT_VERSION}"

CURRENT_VERSION=$(kubectl get deployment/frontend -n ${NAMESPACE} -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $3}')
echo "Frontend版本号: ${CURRENT_VERSION}"
echo "===================="