#版本号参数传入
VERSION=$1

#构建kb_service
docker build --platform linux/amd64 -f Dockerfile.kb_service -t harbor.kunxiangtech.com:8443/kb-platform/kb-service:${VERSION} ../../
#构建iam
docker build --platform linux/amd64 -f Dockerfile.iam -t harbor.kunxiangtech.com:8443/kb-platform/iam:${VERSION} ../../

#构建workflow
docker build --platform linux/amd64 -f Dockerfile.workflow -t harbor.kunxiangtech.com:8443/kb-platform/workflow:${VERSION} ../../

#构建gateway
docker build --platform linux/amd64 -f Dockerfile.gateway -t harbor.kunxiangtech.com:8443/kb-platform/gateway:${VERSION} ../../

#构建frontend
docker build --platform linux/amd64 -f Dockerfile.frontend -t harbor.kunxiangtech.com:8443/kb-platform/frontend:${VERSION} ../../../

#推送
docker push harbor.kunxiangtech.com:8443/kb-platform/kb-service:${VERSION}
docker push harbor.kunxiangtech.com:8443/kb-platform/iam:${VERSION}
docker push harbor.kunxiangtech.com:8443/kb-platform/workflow:${VERSION}
docker push harbor.kunxiangtech.com:8443/kb-platform/gateway:${VERSION}
docker push harbor.kunxiangtech.com:8443/kb-platform/frontend:${VERSION}