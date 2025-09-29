# 先创建namespace
kubectl apply -f ./namespace/namespace.yaml

# 然后创建secret
kubectl apply -f ./secret/iam-secret.yaml
kubectl apply -f ./secret/workflow-secret.yaml
kubectl apply -f ./secret/kb-service-secret.yaml

# 然后创建configmap
kubectl apply -f ./configmap/iam-configmap.yaml
kubectl apply -f ./configmap/workflow-configmap.yaml
kubectl apply -f ./configmap/kb-service-configmap.yaml

# 然后创建statefulset
kubectl apply -f ./statefulset/postgres-master-stateful.yaml
kubectl apply -f ./statefulset/postgres-replica-stateful.yaml

# 然后创建deployment
kubectl apply -f ./deployment/iam-deployment.yaml
kubectl apply -f ./deployment/workflow-deployment.yaml
kubectl apply -f ./deployment/kb-service-deployment.yaml
kubectl apply -f ./deployment/gateway-deployment.yaml

# 初始化数据库
POSTGRES_POD=$(kubectl get pods -n kb-platform -l app=postgres-master -o jsonpath='{.items[0].metadata.name}')
if [ -z "$POSTGRES_POD" ]; then
    echo "❌ 未找到 PostgreSQL master pod, 请确保 StatefulSet 已部署"
    exit 1
fi

kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "CREATE DATABASE kb_platform;"
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb_platform -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"
kubectl exec -i $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb_platform < ../database/init-database.sql
