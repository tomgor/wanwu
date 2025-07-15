# EverArt MCP Server

## 概述

EverArt MCP Server是专为工业设计可视化设计的AI图像生成服务器。该服务器提供基于人工智能的图像创建和编辑能力，支持工业产品设计、技术示意图生成、培训材料制作等功能，适用于工业设计、技术文档配图、产品展示、培训教育等可视化需求场景。

## 主要特性

- AI驱动的图像生成和编辑
- 工业风格图像定制
- 技术示意图自动生成
- 多种图像格式支持
- 批量图像处理能力
- 高质量图像输出

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- EverArt API访问密钥
- 稳定的网络连接
- 充足的存储空间

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-everart`
2. 获取EverArt API密钥
3. 配置API参数和图像设置
4. 启动MCP服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-everart

# 环境变量配置
export EVERART_API_KEY=your_everart_api_key_here
export DEFAULT_IMAGE_SIZE=1024x1024
export OUTPUT_FORMAT=PNG
export QUALITY_LEVEL=high
export MAX_CONCURRENT_REQUESTS=5
```

## 应用场景

### 工业产品设计
- 产品概念设计图
- 产品外观效果图
- 零部件设计草图
- 产品包装设计
- 工业设备外观设计

### 技术文档配图
- 工艺流程示意图
- 设备结构图解
- 安装指导图片
- 操作步骤配图
- 安全警示图标

### 培训教育材料
- 培训课件配图
- 安全教育图片
- 操作指南插图
- 设备认知图片
- 工艺原理图解

## 工具功能

### 图像生成
- **generate_image**: 根据文字描述生成图像
- **create_industrial_diagram**: 创建工业技术图解
- **generate_product_mockup**: 生成产品效果图
- **create_safety_icons**: 创建安全标识图标

### 图像编辑
- **edit_image**: 编辑现有图像
- **add_annotations**: 添加标注和说明
- **adjust_style**: 调整图像风格
- **resize_image**: 调整图像尺寸

### 批量处理
- **batch_generate**: 批量生成图像
- **apply_template**: 应用设计模板
- **create_variations**: 创建图像变体
- **export_formats**: 导出多种格式

## 使用示例

### 生成设备操作示意图
```json
{
  "prompt": "工业机器人在汽车生产线上进行焊接作业，展示机器人臂部动作和安全防护区域，工业风格，技术示意图",
  "style": "technical_illustration",
  "size": "1024x768",
  "include_annotations": true,
  "color_scheme": "industrial_blue"
}
```

### 创建安全警示图标
```json
{
  "type": "safety_icon",
  "warning_type": "high_voltage",
  "text": "高压危险",
  "style": "ISO_standard",
  "background_color": "yellow",
  "text_color": "black",
  "size": "256x256"
}
```

### 生成产品展示图
```json
{
  "product_type": "工业控制器",
  "features": ["触摸屏", "多个接口", "LED指示灯", "金属外壳"],
  "viewing_angle": "45度斜视",
  "background": "clean_white",
  "lighting": "professional",
  "resolution": "high"
}
```

### 制作工艺流程图
```json
{
  "process_name": "钢材热处理工艺",
  "steps": [
    {"name": "加热", "temperature": "850°C", "time": "2小时"},
    {"name": "保温", "temperature": "850°C", "time": "1小时"},
    {"name": "淬火", "medium": "油冷", "temperature": "室温"},
    {"name": "回火", "temperature": "400°C", "time": "3小时"}
  ],
  "style": "flowchart",
  "include_parameters": true
}
```

### 创建培训材料配图
```json
{
  "training_topic": "PLC编程基础",
  "image_type": "educational_diagram",
  "content": "显示PLC输入输出连接示意图，包含传感器、执行器和控制逻辑",
  "annotation_level": "detailed",
  "target_audience": "初学者"
}
```

## 工业应用优势

1. **专业定制**: 针对工业领域的专门优化
2. **快速生成**: 大幅减少设计制作时间
3. **一致性**: 保持企业视觉风格统一
4. **成本节约**: 减少外包设计费用
5. **灵活修改**: 便于快速调整和优化
6. **标准化**: 符合工业标准和规范

## 图像风格选项

### 技术图解风格
- 工程制图风格
- 等轴测图风格
- 爆炸图风格
- 剖面图风格
- 三维渲染风格

### 安全标识风格
- ISO国际标准
- 国家标准（GB）
- 企业定制标准
- 行业专用标准
- 多语言版本

### 产品展示风格
- 产品摄影风格
- 概念设计风格
- CAD渲染风格
- 手绘草图风格
- 科技感风格

## 质量控制

- 图像分辨率检查
- 色彩准确性验证
- 文字清晰度确保
- 比例准确性检验
- 风格一致性维护

## 版权和合规

- 确保生成内容的原创性
- 遵守知识产权法规
- 符合行业标准要求
- 避免敏感内容
- 提供使用授权

## 性能优化

- 图像生成速度优化
- 批量处理效率
- 内存使用管理
- 网络传输优化
- 缓存机制应用

## 集成特性

- 与CAD软件集成
- 文档系统集成
- 培训平台集成
- 内容管理系统集成
- 版本控制支持

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考EverArt API文档。 