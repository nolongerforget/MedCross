# MedCross医疗数据跨链共享平台 - 后端开发文档

## 目录

1. [系统架构概述](#系统架构概述)
2. [后端技术栈](#后端技术栈)
3. [跨链网关实现原理](#跨链网关实现原理)
4. [API接口设计](#api接口设计)
5. [与前端交互关系](#与前端交互关系)
6. [数据模型](#数据模型)
7. [部署指南](#部署指南)

## 系统架构概述

MedCross是一个基于区块链技术的医疗数据跨链共享平台，旨在实现医疗数据的安全共享和可信交换。系统采用双链架构，结合以太坊和Hyperledger Fabric的优势，通过跨链网关实现数据的互通和同步。

### 系统组件

1. **后端服务**：基于Go语言的RESTful API服务，处理前端请求并与区块链交互
2. **以太坊区块链**：存储医疗数据的元数据和授权信息，提供公开验证
3. **Hyperledger Fabric区块链**：存储完整的医疗数据和详细的访问控制策略
4. **跨链网关**：协调两个区块链之间的数据同步和状态一致性
5. **前端应用**：基于React的用户界面，提供数据上传、查询和授权管理功能

### 系统架构图

```
+----------------+     +----------------+     +----------------+
|                |     |                |     |                |
|  前端应用      | <-> |  后端服务      | <-> |  数据库        |
|  (React)       |     |  (Go/Gin)      |     |  (PostgreSQL)  |
|                |     |                |     |                |
+----------------+     +----------------+     +----------------+
                              ^   ^
                              |   |
                  +-----------+   +-----------+
                  |                           |
        +-----------------+         +-----------------+
        |                 |         |                 |
        |  以太坊区块链   | <-----> |  Fabric区块链   |
        |                 |         |                 |
        +-----------------+         +-----------------+
                       \               /
                        \             /
                      +------------------+
                      |                  |
                      |   跨链网关       |
                      |                  |
                      +------------------+
```

## 后端技术栈

- **编程语言**：Go (Golang)
- **Web框架**：Gin
- **数据库**：PostgreSQL (通过GORM进行ORM映射)
- **区块链交互**：
  - 以太坊：go-ethereum客户端
  - Fabric：Fabric SDK Go
- **认证**：JWT (JSON Web Token)
- **API文档**：Swagger

## 跨链网关实现原理

跨链网关是MedCross平台的核心组件，负责协调以太坊和Hyperledger Fabric区块链之间的数据同步和交互。

### 跨链网关核心功能

1. **数据格式转换**：将不同区块链的数据结构进行转换，确保数据格式兼容
2. **交易验证**：验证跨链交易的有效性，确保数据一致性
3. **状态同步**：保持两个区块链上数据状态的同步
4. **事件监听**：监听区块链事件并触发相应的处理逻辑
5. **交易队列管理**：管理跨链交易的处理队列，确保交易顺序和完整性

### 跨链数据传输流程

1. **发起跨链交易**：通过`InitiateDataTransfer`方法发起数据传输请求
2. **创建交易记录**：生成唯一的交易ID并记录交易状态
3. **数据获取与转换**：从源区块链获取数据并转换为目标区块链格式
4. **目标链上传**：将转换后的数据上传到目标区块链
5. **更新交易状态**：更新交易状态为已完成或失败
6. **创建跨链引用**：建立两个区块链上数据的引用关系

### 数据一致性验证

跨链网关通过`VerifyDataConsistency`方法验证两个区块链上数据的一致性：

1. 从两个区块链获取同一数据ID的数据
2. 比较数据哈希值和关键字段
3. 如发现不一致，记录差异并触发同步操作

## API接口设计

### 认证相关接口

| 接口路径 | 方法 | 描述 | 请求参数 | 响应 |
|---------|------|------|---------|------|
| `/api/v1/auth/login` | POST | 用户登录 | username, password | token, expiresAt, user |
| `/api/v1/auth/register` | POST | 用户注册 | username, password, email, fullName, orgName, role | success |

### 医疗数据管理接口

| 接口路径 | 方法 | 描述 | 请求参数 | 响应 |
|---------|------|------|---------|------|
| `/api/v1/data/upload` | POST | 上传医疗数据 | dataHash, dataType, description, tags, isConfidential | dataId, fabricTxId |
| `/api/v1/data/list` | GET | 获取用户拥有的数据 | - | data列表 |
| `/api/v1/data/search` | GET | 搜索医疗数据 | dataType, keyword, tag | data列表 |
| `/api/v1/data/:id` | GET | 获取数据详情 | id (路径参数) | data详情 |

### 授权管理接口

| 接口路径 | 方法 | 描述 | 请求参数 | 响应 |
|---------|------|------|---------|------|
| `/api/v1/data/:id/authorize` | POST | 授权数据访问 | authorizedUserId, startTime, endTime | success |
| `/api/v1/data/:id/authorize/:userId` | DELETE | 撤销授权 | id, userId (路径参数) | success |
| `/api/v1/data/:id/authorizations` | GET | 获取授权列表 | id (路径参数) | authorizations列表 |
| `/api/v1/data/:id/access-logs` | GET | 获取访问记录 | id (路径参数) | accessLogs列表 |

### 区块链记录接口

| 接口路径 | 方法 | 描述 | 请求参数 | 响应 |
|---------|------|------|---------|------|
| `/api/v1/blockchain/transactions` | GET | 获取交易记录 | - | transactions列表 |
| `/api/v1/blockchain/transactions/:txHash` | GET | 获取交易详情 | txHash (路径参数) | transaction详情 |

### 跨链操作接口

| 接口路径 | 方法 | 描述 | 请求参数 | 响应 |
|---------|------|------|---------|------|
| `/api/v1/crosschain/sync/:id` | POST | 跨链同步数据 | id (路径参数), sourceChain | success |
| `/api/v1/crosschain/status/:id` | GET | 获取跨链状态 | id (路径参数) | status |

## 与前端交互关系

### 数据上传页面 (DataUploadPage)

- **对应接口**: `/api/v1/data/upload`
- **交互流程**:
  1. 用户在前端选择文件并填写元数据信息
  2. 前端调用上传接口，将数据哈希和元数据发送到后端
  3. 后端将数据上传到Fabric区块链，并在以太坊上创建引用
  4. 返回数据ID和交易哈希给前端显示

### 数据查询页面 (DataQueryPage)

- **对应接口**: `/api/v1/data/search`, `/api/v1/data/list`
- **交互流程**:
  1. 用户输入搜索条件（关键词、数据类型、标签等）
  2. 前端调用搜索接口，将查询参数发送到后端
  3. 后端从区块链获取符合条件的数据列表
  4. 前端展示搜索结果，支持排序和筛选

### 数据详情页面 (DataDetailPage)

- **对应接口**: `/api/v1/data/:id`, `/api/v1/data/:id/access-logs`
- **交互流程**:
  1. 用户点击数据项查看详情
  2. 前端调用详情接口，获取数据完整信息
  3. 后端检查用户访问权限，记录访问日志
  4. 前端展示数据详情、授权历史和区块链记录

### 授权管理页面 (AuthManagementPage)

- **对应接口**: `/api/v1/data/:id/authorize`, `/api/v1/data/:id/authorize/:userId`, `/api/v1/data/:id/authorizations`
- **交互流程**:
  1. 用户选择要授权的数据和接收方
  2. 前端调用授权接口，发送授权参数
  3. 后端在两个区块链上创建授权记录
  4. 前端更新授权列表显示

### 区块链记录页面 (BlockchainRecordsPage)

- **对应接口**: `/api/v1/blockchain/transactions`, `/api/v1/blockchain/transactions/:txHash`
- **交互流程**:
  1. 用户查看区块链交易记录
  2. 前端调用交易记录接口，获取交易列表
  3. 用户可点击交易哈希查看详情
  4. 前端展示交易详情，包括区块信息和操作类型

## 数据模型

### MedicalData (医疗数据)

```go
type MedicalData struct {
	DataID         string         // 数据唯一标识符
	DataHash       string         // 数据哈希值（IPFS或其他存储系统的引用）
	DataType       string         // 数据类型（如：电子病历、影像数据、基因组数据等）
	Description    string         // 数据描述
	Tags           []string       // 数据标签
	Owner          string         // 数据所有者（用户ID）
	Timestamp      int64          // 上传时间戳
	IsConfidential bool           // 是否为机密数据
	EthereumTxID   string         // 以太坊链上对应的交易ID
	FabricTxID     string         // Fabric链上对应的交易ID
}
```

### Authorization (授权信息)

```go
type Authorization struct {
	DataID         string         // 被授权的数据ID
	OwnerID        string         // 数据所有者ID
	AuthorizedUser string         // 被授权的用户ID
	StartTime      int64          // 授权开始时间
	EndTime        int64          // 授权结束时间（0表示永久授权）
	IsActive       bool           // 授权是否有效
}
```

### AccessRecord (访问记录)

```go
type AccessRecord struct {
	DataID     string    // 访问的数据ID
	Accessor   string    // 访问者ID
	Timestamp  int64     // 访问时间戳
	Operation  string    // 操作类型（查看、下载等）
	ChainType  string    // 区块链类型（Ethereum或Fabric）
}
```

### CrossChainReference (跨链引用)

```go
type CrossChainReference struct {
	DataID       string         // 数据ID
	EthereumTxID string         // 以太坊交易ID
	FabricTxID   string         // Fabric交易ID
	SyncStatus   string         // 同步状态（已同步、同步中、同步失败）
	LastSyncTime int64          // 最后同步时间
}
```

## 部署指南

### 环境要求

- Go 1.16+
- PostgreSQL 12+
- 以太坊节点（Geth或Infura）
- Hyperledger Fabric网络（v2.2+）

### 配置文件

配置文件位于`backend/config/config.yaml`，包含以下主要配置项：

```yaml
server:
  port: 8080
  mode: development

database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: medcross

eth:
  rpcUrl: https://goerli.infura.io/v3/your-api-key
  contractAddress: 0x1234567890abcdef1234567890abcdef12345678
  privateKey: your-private-key

fabric:
  configPath: connection-profile.yaml
  channelId: mychannel
  chaincodeName: medicaldata
  userName: admin
  orgName: org1
```

### 启动步骤

1. 克隆代码库
2. 配置环境变量和配置文件
3. 安装依赖：`go mod download`
4. 初始化数据库：`go run database/migration.go`
5. 启动后端服务：`go run backend/main.go`

### 部署到生产环境

1. 构建二进制文件：`go build -o medcross-backend backend/main.go`
2. 配置系统服务或Docker容器
3. 设置环境变量为生产环境
4. 启动服务并配置反向代理（如Nginx）