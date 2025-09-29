# 后端
## 代码拉取
```bash
git clone https://gitee.com/sichuan-shutong-zhihui-data/k-base.git
```
## 升级推包
```bash
cd k-base/platform/deployment/dockerfile/
sh docker-build.sh vX.X.X # 自带版本号
```
## k8s部署
修改各个deployment的image的版本号之后执行
```bash
cd k-base/platform/deployment/k8s/
make deploy
#状态检查完之后
make check
```
## k8s卸载
```bash
cd k-base/platform/deployment/k8s/
make uninstall
```