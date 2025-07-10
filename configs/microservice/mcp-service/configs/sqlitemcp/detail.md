# SQLite MCP Server

## 概述

SQLite MCP Server是专为工业数据管理设计的轻量级数据库集成服务器。该服务器提供完整的SQLite数据库操作能力，支持复杂的SQL查询和事务处理，适用于工业生产数据记录、设备状态监控、质量控制数据管理等场景。

## 主要特性

- 完整的SQLite数据库操作支持
- 复杂SQL查询和事务处理
- 轻量级部署，适合边缘计算
- 高性能的本地数据存储
- 支持并发访问和数据完整性
- 简单易用的数据库管理

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- SQLite数据库文件访问权限

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-sqlite`
2. 配置数据库路径
3. 通过stdio模式运行服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-sqlite
# 配置数据库路径后启动
./sqlite-mcp-server --db-path ./industrial.db
```

## 应用场景

### 工业生产数据记录
- 生产线数据采集和存储
- 设备运行参数记录
- 产品质量检测数据
- 生产统计和报表

### 设备状态监控
- 设备维护记录管理
- 故障历史数据存储
- 设备性能数据分析
- 预防性维护计划

### 质量控制数据管理
- 质量检测数据分析
- 工艺参数存储
- 不合格品追溯
- 质量趋势分析

## 工具功能

### 数据查询
- **query**: 执行SQL查询语句，支持复杂的数据检索
- **list_tables**: 列出数据库中的所有表
- **describe_table**: 获取表的结构信息

### 数据操作
- **execute**: 执行SQL命令（INSERT、UPDATE、DELETE等）
- **create_table**: 创建新的数据表

## 使用示例

### 创建生产数据表
```sql
CREATE TABLE production_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    equipment_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    temperature REAL,
    pressure REAL,
    production_count INTEGER
);
```

### 查询设备状态
```sql
SELECT equipment_id, AVG(temperature) as avg_temp, COUNT(*) as records
FROM production_data 
WHERE timestamp >= datetime('now', '-1 day')
GROUP BY equipment_id;
```

### 插入生产数据
```sql
INSERT INTO production_data (equipment_id, temperature, pressure, production_count)
VALUES ('EQ001', 75.5, 2.3, 150);
```

## 工业应用优势

1. **轻量级部署**: 适合工业现场和边缘计算环境
2. **高可靠性**: SQLite的ACID特性确保数据完整性
3. **零配置**: 无需复杂的数据库服务器配置
4. **高性能**: 优化的本地存储，读写速度快
5. **可移植性**: 数据库文件可以轻松备份和迁移

## 性能特点

- 支持高频数据写入
- 快速查询响应
- 低内存占用
- 适合嵌入式系统

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考SQLite官方文档。 