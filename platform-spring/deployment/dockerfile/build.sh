VERSION=$1

docker build --platform linux/amd64 -f Dockerfile.platform-spring -t harbor.kunxiangtech.com:8443/kb-platform/platform-spring:${VERSION} ../../ 
docker push harbor.kunxiangtech.com:8443/kb-platform/platform-spring:${VERSION}

docker build --platform linux/amd64 -f Dockerfile.frontend -t harbor.kunxiangtech.com:8443/kb-platform/frontend:${VERSION} ../../../
docker push harbor.kunxiangtech.com:8443/kb-platform/frontend:${VERSION}