# 快速删除所有资源
kubectl delete -f ./deployment/gateway-deployment.yaml --ignore-not-found
kubectl delete -f ./deployment/kb-service-deployment.yaml --ignore-not-found
kubectl delete -f ./deployment/workflow-deployment.yaml --ignore-not-found
kubectl delete -f ./deployment/iam-deployment.yaml --ignore-not-found

kubectl delete -f ./statefulset/postgres-replica-stateful.yaml --ignore-not-found
kubectl delete -f ./statefulset/postgres-master-stateful.yaml --ignore-not-found

kubectl delete -f ./configmap/gateway-configmap.yaml --ignore-not-found
kubectl delete -f ./configmap/kb-service-configmap.yaml --ignore-not-found
kubectl delete -f ./configmap/workflow-configmap.yaml --ignore-not-found
kubectl delete -f ./configmap/iam-configmap.yaml --ignore-not-found

kubectl delete -f ./secret/kb-service-secret.yaml --ignore-not-found
kubectl delete -f ./secret/workflow-secret.yaml --ignore-not-found
kubectl delete -f ./secret/iam-secret.yaml --ignore-not-found

kubectl delete -f ./namespace/namespace.yaml --ignore-not-found
