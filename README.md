# MedCross - 医疗数据跨链共享平台

## 项目概述

MedCross是一个专注于医疗数据跨链查询的平台，通过区块链技术实现医疗数据的安全存储和跨链查询。该平台集成了以太坊和Hyperledger Fabric两条区块链，实现了跨链数据共享的功能。用户可以选择将数据上传到任意一条链上，并能够跨链查询所有数据。

## 项目结构

```
MedCross/
├── medcrossfrontend/       # 前端应用（React）
├── contracts/              # 智能合约
│   ├── ethereum/           # 以太坊智能合约（Solidity）
│   └── fabric/             # Hyperledger Fabric链码（Go）
├── backend/                # 后端API服务（Go + Gin）
└── crosschain-gateway/     # 跨链网关服务（Go）
```

## 核心功能

### 1. 数据上传

- 支持医疗数据上传到以太坊或Hyperledger Fabric链
- 用户可选择目标区块链
- 支持多种医疗数据类型
- 提供元数据标记功能

### 2. 跨链数据查询

- 统一查询接口，支持跨链数据检索
- 关键词搜索功能
- 按数据类型、上传时间等多维度筛选
- 数据格式转换和验证

## 组件说明

### 前端 (medcrossfrontend)

基于React的用户界面，提供：
- 用户注册/登录（同时创建区块链账户）
- 数据上传界面（可选择目标区块链）
- 跨链数据查询界面
- 数据详情展示

### 智能合约 (contracts)

#### 以太坊智能合约 (ethereum)

使用Solidity语言编写的智能合约，主要负责：
- 医疗数据元数据存储
- 数据索引管理
- 跨链数据查询接口

#### Hyperledger Fabric链码 (fabric)

使用Go语言编写的链码，主要负责：
- 医疗数据详细记录
- 数据索引管理
- 跨链数据查询接口

### 后端API (backend)

基于Go语言和Gin框架的后端服务，提供：
- RESTful API接口
- 区块链交互逻辑
- 用户账户管理（与区块链账户关联）

### 跨链网关 (crosschain-gateway)

基于Go语言的跨链服务，负责：
- 跨链数据查询协调
- 数据格式转换
- 数据验证和一致性检查

## 技术栈

- 前端：React, TailwindCSS
- 后端：Go, Gin框架
- 区块链：以太坊, Hyperledger Fabric
- 跨链技术：自定义跨链网关
- 数据库：PostgreSQL（用户数据）

## 用户流程

1. 用户注册/登录（同时创建区块链账户）
2. 上传医疗数据（选择目标区块链）
3. 查询医疗数据（跨链查询）
4. 查看数据详情
