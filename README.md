# MedCross - 医疗数据跨链共享平台

## 项目概述

MedCross是一个医疗数据跨链共享平台，旨在通过区块链技术实现医疗数据的安全共享和授权访问。该平台集成了以太坊和Hyperledger Fabric两条区块链，实现了跨链数据共享的功能。

## 项目结构

```
MedCross/
├── medcrossfrontend/       # 前端应用（React）
├── contracts/              # 智能合约
│   ├── ethereum/           # 以太坊智能合约（Solidity）
│   └── fabric/             # Hyperledger Fabric链码（Go）
└── backend/                # 后端API服务（Go + Gin）
```

## 组件说明

### 前端 (medcrossfrontend)

基于React的用户界面，提供数据上传、查询、授权管理等功能。

### 智能合约 (contracts)

#### 以太坊智能合约 (ethereum)

使用Solidity语言编写的智能合约，主要负责：
- 医疗数据元数据存储
- 数据访问权限控制
- 跨链数据索引

#### Hyperledger Fabric链码 (fabric)

使用Go语言编写的链码，主要负责：
- 医疗数据详细记录
- 数据共享授权管理
- 数据访问审计

### 后端API (backend)

基于Go语言和Gin框架的后端服务，提供：
- RESTful API接口
- 区块链交互逻辑
- 跨链数据同步
- 用户认证和授权

## 功能特点

- 医疗数据安全上传和存储
- 基于区块链的数据访问授权
- 跨链数据共享和查询
- 完整的数据操作审计追踪
- 用户友好的界面
