# SSH MCP Server

## 概述

SSH MCP Server是专为工业设备远程管理设计的安全连接服务器。该服务器提供与远程设备的SSH连接能力，支持命令执行、文件传输、系统监控等功能，适用于工业自动化设备管理、远程维护、系统配置和故障排除等场景。

## 主要特性

- 安全的SSH连接和认证
- 远程命令执行和脚本运行
- 文件上传和下载功能
- 系统状态监控和日志查看
- 多设备并发连接管理
- 连接状态监控和错误处理

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- SSH客户端库支持
- 目标设备SSH服务开启

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-ssh`
2. 配置SSH连接参数
3. 设置认证方式（密钥或密码）
4. 启动MCP服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-ssh

# 环境变量配置
export SSH_HOST=192.168.1.100
export SSH_PORT=22
export SSH_USERNAME=industrial
export SSH_PRIVATE_KEY_PATH=/path/to/private/key
# 或使用密码认证
export SSH_PASSWORD=secure_password
```

## 应用场景

### 工业设备远程管理
- PLC设备远程配置
- HMI系统维护和更新
- 传感器数据采集设备管理
- 网关设备配置和监控
- 边缘计算节点管理

### 系统维护和监控
- 远程故障诊断和排除
- 系统日志查看和分析
- 配置文件更新和备份
- 软件包安装和更新
- 系统性能监控

### 自动化运维
- 批量设备配置部署
- 定时任务和脚本执行
- 系统状态自动检查
- 远程重启和维护
- 安全补丁自动应用

## 工具功能

### 连接管理
- **ssh_connect**: 建立SSH连接到远程设备
- **ssh_disconnect**: 断开SSH连接
- **check_connection**: 检查连接状态

### 命令执行
- **execute_command**: 执行远程命令
- **run_script**: 运行批处理脚本
- **get_system_info**: 获取系统信息

### 文件操作
- **upload_file**: 上传文件到远程设备
- **download_file**: 从远程设备下载文件
- **list_directory**: 列出远程目录内容

## 使用示例

### 连接到工业设备
```json
{
  "host": "192.168.1.100",
  "port": 22,
  "username": "industrial",
  "auth_method": "key",
  "private_key_path": "/keys/industrial_key"
}
```

### 执行设备检查命令
```bash
# 检查PLC设备状态
systemctl status plc-service

# 查看网络连接
netstat -an | grep :502

# 检查磁盘空间
df -h

# 查看系统负载
uptime
```

### 批量设备配置
```bash
#!/bin/bash
# 更新设备配置
cp /tmp/new_config.xml /opt/industrial/config/
systemctl restart industrial-service
echo "配置更新完成: $(date)"
```

### 文件传输示例
```json
{
  "action": "upload",
  "local_path": "./config/device_config.xml",
  "remote_path": "/opt/industrial/config/device_config.xml",
  "backup": true
}
```

## 工业应用优势

1. **安全可靠**: SSH加密传输确保通信安全
2. **远程访问**: 支持跨网络的设备管理
3. **批量操作**: 同时管理多个工业设备
4. **故障诊断**: 快速远程排除设备故障
5. **自动化部署**: 支持配置和软件的自动化部署
6. **审计跟踪**: 完整的操作日志和审计记录

## 安全特性

- SSH密钥认证支持
- 连接超时和重试机制
- 操作权限控制
- 连接日志记录
- 防暴力破解保护

## 监控和日志

- 连接状态实时监控
- 命令执行历史记录
- 错误日志和异常处理
- 性能指标统计
- 安全事件记录

## 最佳实践

1. 使用SSH密钥认证而非密码
2. 定期更新和轮换访问凭据
3. 限制SSH访问的IP地址范围
4. 定期备份重要配置文件
5. 监控异常连接和操作

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考SSH协议相关文档。 