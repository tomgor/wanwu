# Memory MCP Server

## 概述

Memory MCP Server是专为工业知识管理设计的智能记忆服务器。该服务器提供持久化知识存储、语义检索、知识关联等功能，适用于工业技术文档管理、专家知识库、故障诊断经验积累、工艺参数优化等知识密集型工业应用场景。

## 主要特性

- 智能知识存储和索引
- 语义相似度检索
- 知识图谱构建和关联
- 多模态内容支持
- 版本控制和历史追踪
- 分布式知识同步

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- 向量数据库支持（如Chroma、Pinecone等）
- 充足的存储空间

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-memory`
2. 配置向量数据库连接
3. 设置知识库存储路径
4. 启动MCP服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-memory

# 环境变量配置
export MEMORY_STORAGE_PATH=/opt/industrial/knowledge
export VECTOR_DB_URL=http://localhost:8000
export EMBEDDING_MODEL=text-embedding-ada-002
export MAX_MEMORY_SIZE=10GB
```

## 应用场景

### 工业技术文档管理
- 设备操作手册存储
- 技术规范和标准管理
- 工艺流程文档
- 安全操作规程
- 维护手册和指南

### 专家知识库
- 故障诊断经验库
- 设备维护知识
- 工艺优化案例
- 质量控制经验
- 安全事故分析

### 智能问答系统
- 技术问题快速检索
- 相似问题推荐
- 解决方案匹配
- 最佳实践推荐
- 历史案例查询

## 工具功能

### 知识存储
- **store_memory**: 存储新的知识条目
- **update_memory**: 更新现有知识
- **delete_memory**: 删除知识条目
- **list_memories**: 列出所有存储的知识

### 知识检索
- **search_memories**: 语义搜索相关知识
- **get_similar**: 获取相似知识条目
- **get_related**: 获取关联知识
- **query_knowledge**: 复杂知识查询

### 知识管理
- **categorize_memories**: 知识分类管理
- **tag_memories**: 添加标签和元数据
- **export_knowledge**: 导出知识库
- **import_knowledge**: 导入外部知识

## 使用示例

### 存储故障诊断知识
```json
{
  "title": "PLC通信故障诊断",
  "content": "当PLC与上位机通信中断时，首先检查网络连接状态，确认IP地址配置正确，然后检查通信电缆是否损坏，最后验证通信协议参数设置。",
  "category": "故障诊断",
  "tags": ["PLC", "通信", "网络", "故障"],
  "equipment": "西门子S7-1200",
  "severity": "高",
  "solution_time": "15分钟"
}
```

### 检索相关维护知识
```json
{
  "query": "变频器过热故障处理",
  "category": "设备维护", 
  "max_results": 5,
  "similarity_threshold": 0.8
}
```

### 工艺参数优化记录
```json
{
  "title": "注塑工艺温度优化",
  "process": "塑料注塑成型",
  "parameters": {
    "barrel_temperature": "220°C",
    "mold_temperature": "60°C",
    "injection_pressure": "80MPa",
    "cooling_time": "15秒"
  },
  "result": "产品缺陷率从3%降低到0.5%",
  "optimization_date": "2024-01-15",
  "engineer": "张工程师"
}
```

### 安全事故案例存储
```json
{
  "incident_type": "设备安全事故",
  "description": "操作员在设备运行时进行维护，导致手部受伤",
  "root_cause": "未执行锁定挂牌程序",
  "corrective_actions": [
    "加强安全培训",
    "完善锁定挂牌程序",
    "增加安全防护装置"
  ],
  "lessons_learned": "严格执行安全操作程序的重要性",
  "date": "2024-01-10"
}
```

## 工业应用优势

1. **知识积累**: 系统性积累和保存工业专家知识
2. **快速检索**: 通过语义搜索快速找到相关知识
3. **经验传承**: 防止专家知识流失，促进知识传承
4. **决策支持**: 为工程决策提供历史数据和经验支持
5. **持续改进**: 通过知识分析发现改进机会
6. **培训支持**: 为新员工培训提供知识资源

## 知识组织结构

```
工业知识库/
├── 设备管理/
│   ├── 故障诊断/
│   ├── 预防维护/
│   └── 性能优化/
├── 工艺技术/
│   ├── 生产工艺/
│   ├── 质量控制/
│   └── 工艺优化/
├── 安全管理/
│   ├── 安全规程/
│   ├── 事故案例/
│   └── 应急处理/
└── 技术标准/
    ├── 行业标准/
    ├── 企业标准/
    └── 操作规范/
```

## 智能特性

- 自动知识分类和标签
- 重复知识检测和去重
- 知识有效性评估
- 使用频率统计
- 知识更新提醒

## 数据安全

- 知识库访问权限控制
- 敏感信息加密存储
- 定期数据备份
- 版本控制和恢复
- 审计日志记录

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考相关技术文档。 