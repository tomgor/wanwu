# Sequelize MCP Server

## 概述

Sequelize MCP Server是专为工业数据库ORM操作设计的对象关系映射服务器。该服务器提供基于Sequelize的数据库抽象层，支持多种数据库系统、复杂查询、事务管理等功能，适用于工业数据建模、企业级应用开发、数据库迁移、多数据源集成等场景。

## 主要特性

- 多数据库支持（PostgreSQL、MySQL、SQLite、MSSQL）
- 强大的ORM功能和查询构建器
- 数据模型定义和关系映射
- 事务管理和并发控制
- 数据验证和约束
- 数据库迁移和同步

## 安装配置

### 环境要求
- Node.js 16.0或更高版本
- 支持的数据库系统
- 相应的数据库驱动程序
- 数据库连接权限

### 安装步骤
1. 通过npm安装：`npm install @modelcontextprotocol/server-sequelize`
2. 安装相应的数据库驱动
3. 配置数据库连接参数
4. 启动MCP服务器

### 配置示例
```bash
npm install @modelcontextprotocol/server-sequelize

# 安装数据库驱动
npm install pg pg-hstore          # PostgreSQL
npm install mysql2               # MySQL
npm install sqlite3              # SQLite
npm install tedious              # SQL Server

# 环境变量配置
export DB_DIALECT=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=industrial_db
export DB_USER=industrial_user
export DB_PASSWORD=secure_password
export DB_POOL_MIN=5
export DB_POOL_MAX=30
```

## 应用场景

### 工业数据建模
- 生产数据模型设计
- 设备管理数据结构
- 质量控制数据模型
- 库存管理系统
- 人员组织架构模型

### 企业级应用开发
- ERP系统数据层
- MES制造执行系统
- WMS仓储管理系统
- CRM客户关系管理
- SCM供应链管理

### 数据集成和迁移
- 多数据源集成
- 数据库版本迁移
- 历史数据归档
- 数据同步和复制
- 跨平台数据迁移

## 工具功能

### 模型管理
- **define_model**: 定义数据模型
- **create_associations**: 创建模型关联
- **sync_models**: 同步模型到数据库
- **migrate_schema**: 执行数据库迁移

### 数据操作
- **create_record**: 创建记录
- **find_records**: 查询记录
- **update_records**: 更新记录
- **delete_records**: 删除记录
- **bulk_operations**: 批量操作

### 查询构建
- **build_query**: 构建复杂查询
- **join_tables**: 表连接查询
- **aggregate_data**: 数据聚合
- **raw_query**: 执行原生SQL

### 事务管理
- **begin_transaction**: 开始事务
- **commit_transaction**: 提交事务
- **rollback_transaction**: 回滚事务
- **managed_transaction**: 托管事务

## 使用示例

### 定义工业设备模型
```javascript
const Equipment = sequelize.define('Equipment', {
  id: {
    type: DataTypes.INTEGER,
    primaryKey: true,
    autoIncrement: true
  },
  equipmentCode: {
    type: DataTypes.STRING(50),
    allowNull: false,
    unique: true
  },
  name: {
    type: DataTypes.STRING(100),
    allowNull: false
  },
  type: {
    type: DataTypes.ENUM('PLC', 'HMI', 'Robot', 'Sensor', 'Actuator'),
    allowNull: false
  },
  location: {
    type: DataTypes.STRING(100)
  },
  status: {
    type: DataTypes.ENUM('Running', 'Stopped', 'Maintenance', 'Error'),
    defaultValue: 'Stopped'
  },
  installDate: {
    type: DataTypes.DATE,
    allowNull: false
  },
  lastMaintenanceDate: {
    type: DataTypes.DATE
  }
});
```

### 创建生产订单模型
```javascript
const ProductionOrder = sequelize.define('ProductionOrder', {
  orderNumber: {
    type: DataTypes.STRING(20),
    primaryKey: true
  },
  productCode: {
    type: DataTypes.STRING(50),
    allowNull: false
  },
  plannedQuantity: {
    type: DataTypes.INTEGER,
    allowNull: false,
    validate: {
      min: 1
    }
  },
  actualQuantity: {
    type: DataTypes.INTEGER,
    defaultValue: 0
  },
  startDate: {
    type: DataTypes.DATE,
    allowNull: false
  },
  endDate: {
    type: DataTypes.DATE
  },
  status: {
    type: DataTypes.ENUM('Planned', 'InProgress', 'Completed', 'Cancelled'),
    defaultValue: 'Planned'
  },
  priority: {
    type: DataTypes.INTEGER,
    defaultValue: 5,
    validate: {
      min: 1,
      max: 10
    }
  }
});
```

### 建立模型关联
```javascript
// 设备和维护记录关联
Equipment.hasMany(MaintenanceRecord, {
  foreignKey: 'equipmentId',
  as: 'maintenanceRecords'
});

MaintenanceRecord.belongsTo(Equipment, {
  foreignKey: 'equipmentId',
  as: 'equipment'
});

// 生产订单和产品关联
Product.hasMany(ProductionOrder, {
  foreignKey: 'productCode',
  sourceKey: 'productCode'
});

ProductionOrder.belongsTo(Product, {
  foreignKey: 'productCode',
  targetKey: 'productCode'
});
```

### 复杂查询示例
```javascript
// 查询本月设备维护统计
const maintenanceStats = await Equipment.findAll({
  attributes: [
    'type',
    [sequelize.fn('COUNT', sequelize.col('maintenanceRecords.id')), 'maintenanceCount'],
    [sequelize.fn('AVG', sequelize.col('maintenanceRecords.cost')), 'avgCost']
  ],
  include: [{
    model: MaintenanceRecord,
    as: 'maintenanceRecords',
    where: {
      maintenanceDate: {
        [Op.gte]: new Date(new Date().getFullYear(), new Date().getMonth(), 1)
      }
    },
    attributes: []
  }],
  group: ['Equipment.type'],
  order: [[sequelize.fn('COUNT', sequelize.col('maintenanceRecords.id')), 'DESC']]
});
```

### 事务处理示例
```javascript
// 生产订单完工处理
const completeProductionOrder = async (orderNumber, actualQuantity) => {
  const transaction = await sequelize.transaction();
  
  try {
    // 更新生产订单状态
    await ProductionOrder.update(
      {
        actualQuantity: actualQuantity,
        endDate: new Date(),
        status: 'Completed'
      },
      {
        where: { orderNumber: orderNumber },
        transaction: transaction
      }
    );
    
    // 更新库存
    await Inventory.increment(
      'quantity',
      {
        by: actualQuantity,
        where: { productCode: productCode },
        transaction: transaction
      }
    );
    
    // 记录生产历史
    await ProductionHistory.create({
      orderNumber: orderNumber,
      completedDate: new Date(),
      actualQuantity: actualQuantity
    }, { transaction: transaction });
    
    await transaction.commit();
    return { success: true, message: '生产订单完工处理成功' };
  } catch (error) {
    await transaction.rollback();
    throw error;
  }
};
```

## 工业应用优势

1. **数据一致性**: 强大的事务支持确保数据完整性
2. **跨数据库**: 支持多种数据库系统无缝切换
3. **开发效率**: 简化数据库操作，提高开发速度
4. **类型安全**: TypeScript支持提供类型检查
5. **查询优化**: 智能查询优化和缓存机制
6. **扩展性**: 支持复杂的企业级应用需求

## 数据模型设计最佳实践

### 命名规范
- 使用清晰描述性的模型名称
- 遵循驼峰命名法则
- 建立统一的字段命名规范
- 使用有意义的关联名称

### 索引策略
- 为经常查询的字段建立索引
- 复合索引优化多字段查询
- 避免过度索引影响写入性能
- 定期分析索引使用情况

### 数据验证
- 在模型级别定义验证规则
- 使用数据库约束作为最后防线
- 实施业务逻辑验证
- 错误处理和用户友好提示

## 性能优化

- 查询优化和索引策略
- 连接池配置调优
- 缓存机制应用
- 批量操作优化
- 分页查询实现

## 监控和调试

- SQL查询日志记录
- 性能指标监控
- 慢查询分析
- 连接池状态监控
- 错误日志收集

## 数据库迁移

- 版本化迁移脚本
- 自动化迁移工具
- 回滚策略制定
- 数据备份和恢复
- 生产环境迁移流程

## 技术支持

如需技术支持或报告问题，请联系开发团队或参考Sequelize官方文档。 