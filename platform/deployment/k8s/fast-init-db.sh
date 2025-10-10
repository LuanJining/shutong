
# 初始化数据库
POSTGRES_POD=$(kubectl get pods -n kb-platform -l app=postgres-master -o jsonpath='{.items[0].metadata.name}')
if [ -z "$POSTGRES_POD" ]; then
    echo "❌ 未找到 PostgreSQL master pod, 请确保 StatefulSet 已部署"
    exit 1
fi

kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "CREATE DATABASE kb_platform;"
kubectl exec -it $POSTGRES_POD -n kb-platform -- psql -U postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"
kubectl exec -i $POSTGRES_POD -n kb-platform -- psql -U postgres -d kb_platform < ../database/init-database.sql
