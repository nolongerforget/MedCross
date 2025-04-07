# MedCross 后端开发细则

## 1. 项目概述

MedCross是一个医疗数据跨链共享平台，旨在通过区块链技术实现医疗数据的安全共享和跨机构协作。后端系统采用Go语言开发，基于Gin框架构建RESTful API，支持以太坊和Hyperledger Fabric两种区块链网络的数据交互。

## 2. 技术栈

- **编程语言**: Go 1.18+
- **Web框架**: Gin 1.9.1
- **认证**: JWT (golang-jwt/jwt/v5)
- **跨域**: gin-contrib/cors
- **环境变量**: godotenv
- **加密**: golang.org/x/crypto
- **区块链**: 以太坊、Hyperledger Fabric

## 3. 项目结构

```
backend/
├── controllers/        # 控制器层，处理HTTP请求
│   ├── auth_controller.go
│   └── data_controller.go
├── middleware/        # 中间件
│   └── auth_middleware.go
├── models/            # 数据模型
│   ├── medical_data.go
│   ├── transfer_record.go
│   └── user.go
├── services/          # 业务逻辑层
│   ├── data_service.go
│   ├── gateway_service.go
│   └── user_service.go
├── utils/             # 工具函数
│   ├── chain_converter.go
│   ├── jwt_utils.go
│   └── password_utils.go
├── .env               # 环境变量配置
├── go.mod             # Go模块定义
└── main.go            # 应用入口
```

## 4. 架构设计

### 4.1 MVC架构

项目采用MVC (Model-View-Controller) 架构模式：

- **Model**: 定义数据结构和业务规则
- **View**: 由前端负责，后端仅提供JSON数据
- **Controller**: 处理HTTP请求，调用相应的服务

### 4.2 分层设计

- **控制器层 (Controllers)**: 负责接收HTTP请求，参数验证，调用服务层，返回响应
- **服务层 (Services)**: 实现业务逻辑，处理数据转换，调用区块链接口
- **模型层 (Models)**: 定义数据结构和验证规则
- **中间件层 (Middleware)**: 提供认证、日志、错误处理等横切关注点
- **工具层 (Utils)**: 提供通用功能，如加密、JWT处理、区块链数据转换等

## 5. API设计

### 5.1 认证API

- **POST /api/register**: 用户注册
- **POST /api/login**: 用户登录
- **GET /api/user**: 获取当前用户信息

### 5.2 数据API

- **GET /api/data**: 查询医疗数据
- **POST /api/data**: 上传医疗数据
- **GET /api/data/:id**: 获取特定医疗数据
- **POST /api/data/transfer**: 跨链转移数据
- **GET /api/statistics**: 获取数据统计信息

### 5.3 API响应格式

成功响应：
```json
{
  "data": {...},  // 响应数据
  "message": "操作成功"
}
```

错误响应：
```json
{
  "error": "错误信息"
}
```

## 6. 数据模型

### 6.1 用户模型 (User)

```go
type User struct {
  ID         string    // 用户唯一标识
  Username   string    // 用户名
  Password   string    // 密码（加密存储）
  Name       string    // 真实姓名
  Role       string    // 角色（doctor, researcher, admin等）
  Hospital   string    // 所属医院
  Department string    // 所属科室
  CreatedAt  time.Time // 创建时间
  UpdatedAt  time.Time // 更新时间
}
```

### 6.2 医疗数据模型 (MedicalData)

```go
type MedicalData struct {
  ID        string    // 数据唯一标识
  Owner     string    // 数据所有者的用户ID
  DataHash  string    // IPFS或其他存储系统的哈希值
  DataType  string    // 数据类型（如：影像数据、电子病历等）
  Metadata  string    // JSON格式的元数据
  Timestamp time.Time // 上传时间戳
  Keywords  string    // 关键词，用于搜索
  Chain     string    // 标识数据来源的区块链: "ethereum" 或 "fabric"
}
```

### 6.3 转移记录模型 (TransferRecord)

```go
type TransferRecord struct {
  ID            string    // 记录唯一标识
  DataID        string    // 医疗数据ID
  SourceChain   string    // 源区块链
  TargetChain   string    // 目标区块链
  SourceTxHash  string    // 源链交易哈希
  TargetTxHash  string    // 目标链交易哈希
  Status        string    // 转移状态
  Timestamp     time.Time // 转移时间
  InitiatedBy   string    // 发起者用户ID
}
```

## 7. 跨链网关实现

### 7.1 网关服务 (GatewayService)

跨链网关服务负责处理不同区块链之间的数据交互，主要功能包括：

- 数据查询：跨链查询医疗数据
- 数据转移：将数据从一个区块链转移到另一个区块链
- 数据验证：验证跨链数据的一致性和完整性

### 7.2 链转换器 (ChainConverter)

链转换器负责处理不同区块链之间的数据格式转换：

- **EthereumToFabric**: 将以太坊格式的医疗数据转换为Fabric格式
- **FabricToEthereum**: 将Fabric格式的医疗数据转换为以太坊格式

### 7.3 跨链数据流程

1. 用户发起跨链数据转移请求
2. 后端验证用户权限和数据所有权
3. 调用源链接口获取数据
4. 使用链转换器转换数据格式
5. 调用目标链接口存储数据
6. 记录转移结果并返回响应

## 8. 安全措施

### 8.1 认证与授权

- 使用JWT进行用户认证
- 基于角色的访问控制 (RBAC)
- 中间件验证用户权限

### 8.2 数据安全

- 密码加密存储 (bcrypt)
- HTTPS传输加密
- 敏感信息脱敏处理

### 8.3 区块链安全

- 交易签名验证
- 数据哈希一致性检查
- 智能合约权限控制

## 9. 错误处理

### 9.1 错误类型

- 客户端错误 (4xx)
- 服务器错误 (5xx)
- 区块链错误

### 9.2 错误日志

- 记录错误详情、时间和请求信息
- 区分错误级别：警告、错误、致命错误
- 敏感信息脱敏

## 10. 开发规范

### 10.1 代码风格

- 遵循Go语言官方代码规范
- 使用gofmt格式化代码
- 包名使用小写单词
- 函数名使用驼峰命名法

### 10.2 注释规范

- 包注释：描述包的功能和用途
- 函数注释：描述函数的功能、参数和返回值
- 复杂逻辑注释：解释复杂算法或业务逻辑

### 10.3 错误处理规范

- 使用明确的错误信息
- 避免忽略错误返回值
- 在适当的层级处理错误

## 11. 测试策略

### 11.1 单元测试

- 测试各个函数和方法的功能
- 使用Go标准测试框架
- 模拟外部依赖

### 11.2 集成测试

- 测试API端点
- 测试服务之间的交互
- 测试与区块链的交互

### 11.3 性能测试

- 测试API响应时间
- 测试并发处理能力
- 测试区块链交互性能

## 12. 部署流程

### 12.1 环境配置

- 开发环境：本地开发和测试
- 测试环境：功能测试和集成测试
- 生产环境：正式部署

### 12.2 部署步骤

1. 准备环境变量配置
2. 构建Go应用
3. 部署区块链网络
4. 配置网络和防火墙
5. 启动应用服务
6. 监控应用状态

### 12.3 容器化部署

- 使用Docker容器化应用
- 使用Docker Compose管理多容器应用
- 配置容器网络和存储

## 13. 监控与日志

### 13.1 应用监控

- 监控API请求量和响应时间
- 监控系统资源使用情况
- 监控错误率和异常情况

### 13.2 区块链监控

- 监控区块链节点状态
- 监控交易确认状态
- 监控智能合约执行情况

### 13.3 日志管理

- 集中式日志收集
- 日志分级和过滤
- 日志分析和告警

## 14. 性能优化

### 14.1 API性能

- 使用缓存减少重复计算
- 优化数据库查询
- 实现分页和限流

### 14.2 区块链交互优化

- 批量处理交易
- 异步处理非关键操作
- 优化智能合约调用

## 15. 版本控制与发布

### 15.1 版本命名

- 遵循语义化版本 (SemVer)
- 主版本.次版本.修订号 (如1.2.3)

### 15.2 发布流程

1. 代码审查
2. 测试验证
3. 构建发布包
4. 部署到测试环境
5. 验收测试
6. 部署到生产环境
7. 监控和回滚准备

## 16. 文档维护

### 16.1 API文档

- 使用Swagger或类似工具生成API文档
- 包含请求参数、响应格式和示例
- 定期更新文档

### 16.2 代码文档

- 使用godoc生成代码文档
- 包含包、函数和类型的说明
- 提供使用示例

## 17. 常见问题与解决方案

### 17.1 区块链连接问题

- 检查网络连接和防火墙设置
- 验证节点状态和同步情况
- 检查账户权限和余额

### 17.2 数据一致性问题

- 实现数据校验和比对机制
- 定期同步和修复不一致数据
- 记录和报告数据异常

### 17.3 性能瓶颈

- 识别性能瓶颈点
- 优化代码和算法
- 考虑水平扩展和负载均衡