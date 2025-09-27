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

# rag_private_deploy

请补充说明内容

## 常用git命令

1.git global setup

git config --global user.name 'test'

git config --global user.email 'xxx@email.com'

2.克隆代码至本地

git clone <仓库克隆地址>

git clone https://gitlab.ai-yuanjing.cn/model_extend_new/rag_private_deploy.git

git clone ssh://git@gitlab.ai-yuanjing.cn:54322/model_extend_new/rag_private_deploy.git

3.克隆指定分支

git clone -b <指定分支名> <远程仓库地址>，如: git clone -b dev https://gitlab.ai-yuanjing.cn/model_extend_new/rag_private_deploy.git

4.查看分支

git branch  // 查看所有本地分支

git branch -a //查看本地和远程所有分支

git branch  -r //查看所有远程分支

5.切换分支

git checkout <指定分支名>，如：git checkout dev //切换到指定分支

git checkout -b <指定分支名>  //新建分支，并切换到该分支

6.拉取代码

git pull

7.将本地修改的文件xx添加到暂存区

git add <文件名称>， 如：git add test01

git add -A  提交所有变化

git add -u  提交被修改(modified)和被删除(deleted)文件，不包括新文件(new)

git add .  提交新文件(new)和被修改(modified)文件，不包括被删除(deleted)文件

8.提交暂存区的内容

git commit -m "注释"，如： git commit -m "feat: devops-xxxx, 新增文件xxx，完成xx功能"

9.推送代码

git push //将提交的文件推送到远端仓库

git push --set-upstream origin <分支名称>  // 若该分支远端gitlab中不存在，则使用该命令进行推送