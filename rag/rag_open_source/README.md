# RAG 模型服务部署指南

## 📦 打包流程

### 1. 进入源码容器
```bash
docker exec -it 容器ID /bin/bash
```
### 2. 进入工作目录
```bash
cd /model_extend
```
### 3. 激活Conda环境
```bash
conda activate rag-new
```
### 4. 执行打包脚本
```bash
# 如果脚本没有执行权限，先添加权限
chmod +x rag_pack.sh

# 更新包列表（可选，但推荐）
apt-get update

# 安装binutils，-y选项表示自动同意
apt-get install -y binutils

# 后台运行打包脚本（输出将保存到nohup.out）
nohup ./rag_pack.sh &
```
### 5. 复制生成文件
```bash
mkdir -p /model_extend/opt
mkdir -p /model_extend/Fonts
cp -r /opt/* /model_extend/opt
cp -r /usr/share/fonts/Fonts/* /model_extend/Fonts
```
### 6. 退出容器
```bash
exit
```
### 7. 构建Docker镜像
```bash
#进入model_extend（代码挂载目录，查看容器启动命令）
#ARM64架构：
make docker-image-rag-arm64
#AMD64架构：
make docker-image-rag-amd64
```
