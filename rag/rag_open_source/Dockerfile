# 使用 rag-base 镜像作为基础镜像
FROM crpi-6pj79y7ddzdpexs8.cn-hangzhou.personal.cr.aliyuncs.com/wanwulite/rag-base:v1.1.0-20250919

# 设置工作目录
WORKDIR /model_extend

# 1. 创建目录结构
RUN mkdir -p /model_extend/

# 2. 复制本地文件
COPY ./langchain_rag-master /model_extend/langchain_rag-master
COPY ./minio_project /model_extend/minio_project
COPY ./rag-es-server-unify /model_extend/rag-es-server-unify
COPY ./start.sh /model_extend/start.sh

# 4. 设置脚本权限
RUN chmod +x /model_extend/start.sh

# 5. 容器启动命令
CMD ["/bin/bash", "/model_extend/start.sh"]
