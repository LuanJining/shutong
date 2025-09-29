
# 初始化数据库
POSTGRES_POD=$(kubectl get pods -n kb-platform -l app=postgres-master -o jsonpath='{.items[0].metadata.name}')
if [ -z "$POSTGRES_POD" ]; then
    echo "❌ 未找到 PostgreSQL master pod, 请确保 StatefulSet 已部署"
    exit 1
fi

kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "CREATE DATABASE kb_platform;"
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb_platform -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"
kubectl exec -i $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb_platform < ../database/init-database.sql

# 获取gateway的cluster ip
GATEWAY_CLUSTER_IP=$(kubectl get svc -n kb-platform -l app=gateway -o jsonpath='{.items[0].spec.clusterIP}')
if [ -z "$GATEWAY_CLUSTER_IP" ]; then
    echo "❌ 未找到 Gateway 的 ClusterIP, 请确保 Service 已部署"
    exit 1
fi

# 测试gateway
echo "=== Gateway测试 ==="
curl -s http://$GATEWAY_CLUSTER_IP/api/v1/health | jq .

echo "=== Gateway测试完成 ==="
echo "Gateway URL: http://$GATEWAY_CLUSTER_IP"

# 获取iam的cluster ip
IAM_CLUSTER_IP=$(kubectl get svc -n kb-platform -l app=iam -o jsonpath='{.items[0].spec.clusterIP}')
if [ -z "$IAM_CLUSTER_IP" ]; then
    echo "❌ 未找到 IAM 的 ClusterIP, 请确保 Service 已部署"
    exit 1
fi

# 测试iam
echo "=== IAM测试 ==="
curl -s http://$IAM_CLUSTER_IP/api/v1/health | jq .

echo "=== IAM测试完成 ==="
echo "IAM URL: http://$IAM_CLUSTER_IP"

# 获取workflow的cluster ip
WORKFLOW_CLUSTER_IP=$(kubectl get svc -n kb-platform -l app=workflow -o jsonpath='{.items[0].spec.clusterIP}')
if [ -z "$WORKFLOW_CLUSTER_IP" ]; then
    echo "❌ 未找到 Workflow 的 ClusterIP, 请确保 Service 已部署"
    exit 1
fi

# 测试workflow
echo "=== Workflow测试 ==="
curl -s http://$WORKFLOW_CLUSTER_IP/api/v1/health | jq .

echo "=== Workflow测试完成 ==="
echo "Workflow URL: http://$WORKFLOW_CLUSTER_IP"

# 获取kb-service的cluster ip
KB_SERVICE_CLUSTER_IP=$(kubectl get svc -n kb-platform -l app=kb-service -o jsonpath='{.items[0].spec.clusterIP}')
if [ -z "$KB_SERVICE_CLUSTER_IP" ]; then
    echo "❌ 未找到 KB Service 的 ClusterIP, 请确保 Service 已部署"
    exit 1
fi
