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
kubectl apply -f ./configmap/gateway-configmap.yaml

# 然后创建statefulset
kubectl apply -f ./statefulset/postgres-master-stateful.yaml
kubectl apply -f ./statefulset/postgres-replica-stateful.yaml

# 然后创建deployment
kubectl apply -f ./deployment/iam-deployment.yaml
kubectl apply -f ./deployment/workflow-deployment.yaml
kubectl apply -f ./deployment/kb-service-deployment.yaml
kubectl apply -f ./deployment/gateway-deployment.yaml
