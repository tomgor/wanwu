# Brave Search MCP Server

## 概述

Brave Search MCP Server是专为工业技术资料搜索设计的智能搜索服务器。该服务器提供基于Brave搜索引擎的技术资料检索能力，支持实时网络搜索、技术文档查找、行业资讯获取等功能，适用于工业技术研究、问题解决方案查找、市场信息收集等场景。

## 主要特性

- 基于Brave搜索引擎的强大搜索能力
- 隐私保护的搜索服务
- 实时技术资料和资讯获取
- 多语言搜索支持
- 搜索结果过滤和排序
- API限速和使用监控

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- Brave Search API密钥
- 稳定的网络连接

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-brave-search`
2. 获取Brave Search API密钥
3. 配置API密钥和搜索参数
4. 启动MCP服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-brave-search

# 环境变量配置
export BRAVE_API_KEY=your_brave_api_key_here
export SEARCH_LANGUAGE=zh-CN
export MAX_RESULTS=10
export SAFE_SEARCH=moderate
```

## 应用场景

### 工业技术研究
- 新技术和新工艺调研
- 设备技术规格查询
- 行业标准和规范搜索
- 专利和技术文献检索
- 竞品技术分析

### 问题解决支持
- 故障诊断方案搜索
- 技术问题解决方案
- 最佳实践案例查找
- 专家观点和建议
- 技术论坛和社区资源

### 市场和行业信息
- 行业发展趋势
- 市场价格信息
- 供应商和产品信息
- 行业新闻和动态
- 展会和会议信息

## 工具功能

### 搜索功能
- **web_search**: 执行网络搜索查询
- **search_news**: 搜索相关新闻和资讯
- **search_images**: 搜索相关图片和图表
- **search_videos**: 搜索技术视频和教程

### 搜索优化
- **filter_results**: 按条件过滤搜索结果
- **sort_results**: 对搜索结果排序
- **get_search_suggestions**: 获取搜索建议
- **search_history**: 查看搜索历史

## 使用示例

### 技术问题搜索
```json
{
  "query": "PLC与变频器通信Modbus配置方法",
  "language": "zh-CN",
  "max_results": 10,
  "safe_search": "moderate",
  "time_range": "past_year"
}
```

### 设备技术规格查询
```json
{
  "query": "西门子S7-1500 CPU技术参数 通信接口",
  "search_type": "web",
  "filter": {
    "site": "siemens.com",
    "file_type": "pdf"
  }
}
```

### 行业标准搜索
```json
{
  "query": "工业机器人安全标准 ISO 10218",
  "language": "zh-CN",
  "category": "technical_standards",
  "official_sources_only": true
}
```

### 故障解决方案搜索
```json
{
  "query": "伺服电机编码器故障诊断方法",
  "search_type": "comprehensive",
  "include": ["forums", "technical_docs", "videos"],
  "exclude_ads": true
}
```

### 供应商信息查询
```json
{
  "query": "工业传感器供应商 中国 温度传感器",
  "location": "中国",
  "category": "suppliers",
  "sort_by": "relevance"
}
```

## 搜索优化策略

### 关键词优化
- 使用专业术语和技术名词
- 包含品牌和型号信息
- 添加限定词提高精确度
- 使用同义词扩展搜索范围

### 搜索过滤
```json
{
  "filters": {
    "language": "zh-CN",
    "time_range": "past_year",
    "content_type": ["technical", "academic"],
    "domain_filter": ["industry", "automation"],
    "exclude_sites": ["advertisement", "spam"]
  }
}
```

### 结果排序
- 按相关度排序
- 按时间新旧排序  
- 按权威性排序
- 按地域相关性排序

## 工业应用优势

1. **实时信息**: 获取最新的技术信息和行业动态
2. **全面覆盖**: 搜索范围涵盖全球技术资源
3. **隐私保护**: Brave搜索引擎注重用户隐私
4. **高质量结果**: 过滤低质量和垃圾信息
5. **多语言支持**: 支持中英文等多种语言搜索
6. **专业化**: 针对工业技术领域优化

## 搜索最佳实践

### 查询构建技巧
1. 使用具体的技术术语
2. 包含设备型号和规格
3. 添加应用场景描述
4. 使用引号进行精确搜索
5. 利用布尔操作符组合条件

### 结果评估标准
- 信息来源的权威性
- 内容的时效性
- 技术描述的准确性
- 解决方案的可行性
- 相关性和完整性

## 使用限制

- API调用频率限制
- 每日搜索次数限制
- 搜索结果数量限制
- 特定内容类型限制
- 地域访问限制

## 隐私和安全

- 不记录个人搜索历史
- 加密传输搜索请求
- 不追踪用户行为
- 遵循数据保护法规
- 支持匿名搜索

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考Brave Search API文档。 