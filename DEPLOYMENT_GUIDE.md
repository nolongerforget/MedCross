# MedCross 医疗数据跨链共享平台部署手册

## 目录

- [1. 系统概述](#1-系统概述)
- [2. 环境要求](#2-环境要求)
- [3. 区块链网络部署](#3-区块链网络部署)
- [4. 后端服务部署](#4-后端服务部署)
- [5. 跨链网关部署](#5-跨链网关部署)
- [6. 前端应用部署](#6-前端应用部署)
- [7. 系统配置](#7-系统配置)
- [8. 系统验证](#8-系统验证)
- [9. 常见问题排查](#9-常见问题排查)
- [10. 生产环境部署建议](#10-生产环境部署建议)

## 1. 系统概述

MedCross是一个专注于医疗数据跨链查询的平台，通过区块链技术实现医疗数据的安全存储和跨链查询。该平台集成了以太坊和Hyperledger Fabric两条区块链，实现了跨链数据共享的功能。

系统架构包括以下组件：

- **前端应用**：基于React的用户界面
- **后端API服务**：基于Go和Gin框架的RESTful API
- **跨链网关**：协调以太坊和Fabric链上的数据查询
- **区块链网络**：以太坊和Hyperledger Fabric

## 2. 环境要求

### 2.1 硬件要求

- **最低配置**：
  - CPU: 4核
  - 内存: 8GB RAM
  - 存储: 100GB SSD

- **推荐配置**：
  - CPU: 8核
  - 内存: 16GB RAM
  - 存储: 500GB SSD

### 2.2 软件要求

- **操作系统**：
  - Ubuntu 20.04 LTS 或更高版本（推荐）
  - Windows 10/11 专业版（开发环境）

- **基础软件**：
  - Docker 20.10.x 或更高版本
  - Docker Compose 2.x
  - Git 2.x
  - Go 1.18 或更高版本
  - Node.js 18.x 或更高版本
  - npm 8.x 或更高版本

## 3. 区块链网络部署

### 3.1 以太坊私有链部署

#### 3.1.1 安装Geth

```bash
# Ubuntu系统
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install -y ethereum

# 验证安装
geth version
```

#### 3.1.2 创建创世区块配置

创建`genesis.json`文件：

```json
{
  "config": {
    "chainId": 15,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "8000000",
  "extradata": "0x0000000000000000000000000000000000000000000000000000000000000000<YOUR_ETHEREUM_ADDRESS>0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "alloc": {
    "<YOUR_ETHEREUM_ADDRESS>": { "balance": "100000000000000000000" }
  }
}
```

> 注意：替换`<YOUR_ETHEREUM_ADDRESS>`为您的以太坊地址（不含0x前缀）

#### 3.1.3 初始化和启动私有链

```bash
# 创建数据目录
mkdir -p ~/medcross/ethereum/data

# 初始化创世区块
geth --datadir ~/medcross/ethereum/data init genesis.json

# 启动私有链
geth --datadir ~/medcross/ethereum/data --networkid 15 --http --http.addr 0.0.0.0 --http.port 8545 --http.corsdomain "*" --http.api "eth,net,web3,personal,miner,admin" --allow-insecure-unlock --mine --miner.threads 1 --miner.etherbase "<YOUR_ETHEREUM_ADDRESS>"
```

### 3.2 Hyperledger Fabric网络部署

请参考项目中的`FABRIC_SETUP_GUIDE.md`文件，该文件提供了详细的Fabric环境搭建与私链部署指南。

## 4. 后端服务部署

### 4.1 获取源代码

```bash
# 克隆代码库
git clone <repository_url> medcross
cd medcross
```

### 4.2 配置环境变量

在`backend`目录下创建`.env`文件：

```
PORT=8000
JWT_SECRET=your-secret-key-change-in-production
CORS_ALLOW_ORIGINS=*
ETHEREUM_NODE_URL=http://localhost:8545
FABRIC_CONFIG_PATH=./fabric-config
GATEWAY_URL=http://localhost:8080
```

### 4.3 编译和运行

```bash
cd backend

# 安装依赖
go mod tidy

# 编译
go build -o medcross-backend

# 运行
./medcross-backend
```

### 4.4 使用Docker部署（可选）

在`backend`目录下创建`Dockerfile`：

```dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o medcross-backend

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/medcross-backend /app/
COPY --from=builder /app/config /app/config
COPY --from=builder /app/fabric-config /app/fabric-config

EXPOSE 8000
CMD ["/app/medcross-backend"]
```

构建和运行Docker容器：

```bash
# 构建镜像
docker build -t medcross-backend:latest .

# 运行容器
docker run -d -p 8000:8000 --name medcross-backend \
  -e JWT_SECRET=your-secret-key \
  -e ETHEREUM_NODE_URL=http://host.docker.internal:8545 \
  -e GATEWAY_URL=http://host.docker.internal:8080 \
  medcross-backend:latest
```

## 5. 跨链网关部署

### 5.1 配置环境变量

在`crosschain-gateway`目录下创建`.env`文件：

```
PORT=8080
ETHEREUM_NODE_URL=http://localhost:8545
FABRIC_CONFIG_PATH=./fabric-config
CORS_ALLOW_ORIGINS=*
```

### 5.2 编译和运行

```bash
cd crosschain-gateway

# 安装依赖
go mod tidy

# 编译
go build -o medcross-gateway

# 运行
./medcross-gateway
```

### 5.3 使用Docker部署（可选）

在`crosschain-gateway`目录下创建`Dockerfile`：

```dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o medcross-gateway

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/medcross-gateway /app/
COPY --from=builder /app/fabric-config /app/fabric-config

EXPOSE 8080
CMD ["/app/medcross-gateway"]
```

构建和运行Docker容器：

```bash
# 构建镜像
docker build -t medcross-gateway:latest .

# 运行容器
docker run -d -p 8080:8080 --name medcross-gateway \
  -e ETHEREUM_NODE_URL=http://host.docker.internal:8545 \
  medcross-gateway:latest
```

## 6. 前端应用部署

### 6.1 安装依赖

```bash
cd medcrossfrontend
npm install
```

### 6.2 配置环境变量

创建`.env`文件：

```
VITE_API_URL=http://localhost:8000/api
```

### 6.3 开发环境运行

```bash
npm run dev
```

### 6.4 生产环境构建

```bash
npm run build
```

构建后的文件将位于`dist`目录中，可以使用任何静态文件服务器托管这些文件。

### 6.5 使用Docker部署（可选）

在`medcrossfrontend`目录下创建`Dockerfile`：

```dockerfile
# 构建阶段
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 生产阶段
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

创建`nginx.conf`文件：

```
server {
    listen 80;
    server_name localhost;

    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://backend:8000/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

构建和运行Docker容器：

```bash
# 构建镜像
docker build -t medcross-frontend:latest .

# 运行容器
docker run -d -p 80:80 --name medcross-frontend medcross-frontend:latest
```

## 7. 系统配置

### 7.1 后端配置文件说明

`backend/config/config.yaml`文件包含以下配置项：

```yaml
# 服务器配置
server:
  port: 8000
  jwt_secret: "your-secret-key-change-in-production"
  cors_allow_origins: "*"

# 区块链配置
blockchain:
  # 以太坊配置
  ethereum:
    node_url: "http://localhost:8545"
  
  # Fabric配置
  fabric:
    config_path: "./fabric-config"

# 跨链网关配置
gateway:
  url: "http://localhost:8080"
```

### 7.2 跨链网关配置

跨链网关的配置主要通过环境变量进行设置，关键配置项包括：

- `PORT`: 网关服务端口
- `ETHEREUM_NODE_URL`: 以太坊节点URL
- `FABRIC_CONFIG_PATH`: Fabric配置文件路径

### 7.3 前端配置

前端应用的API地址配置在`.env`文件中：

```
VITE_API_URL=http://localhost:8000/api
```

## 8. 系统验证

### 8.1 验证后端服务

```bash
# 检查后端服务是否正常运行
curl http://localhost:8000/api/health
```

预期输出：

```json
{"status":"ok","message":"服务正常运行"}
```

### 8.2 验证跨链网关

```bash
# 检查跨链网关是否正常运行
curl http://localhost:8080/health
```

预期输出：

```json
{"status":"ok","ethereum":"connected","fabric":"connected"}
```

### 8.3 验证前端应用

在浏览器中访问前端应用（http://localhost:80 或开发环境中的 http://localhost:5173），确认以下功能正常：

1. 用户注册和登录
2. 数据上传（选择目标区块链）
3. 跨链数据查询
4. 数据详情查看

## 9. 常见问题排查

### 9.1 后端服务无法启动

- 检查环境变量是否正确配置
- 检查端口是否被占用：`lsof -i :8000`
- 检查日志输出，查找错误信息

### 9.2 跨链网关连接问题

- 确认以太坊节点是否正常运行
- 检查Fabric网络配置是否正确
- 验证网关服务日志中的连接信息

### 9.3 前端应用无法连接后端

- 确认API URL配置是否正确
- 检查CORS配置是否允许前端域名
- 使用浏览器开发者工具检查网络请求错误

## 10. 生产环境部署建议

### 10.1 安全性建议

- 使用HTTPS保护所有服务通信
- 更改所有默认密钥和密码
- 限制API访问，实施适当的认证和授权
- 定期更新所有依赖包和系统组件

### 10.2 高可用性部署

- 使用负载均衡器分发流量
- 部署多个后端和网关服务实例
- 实施自动扩展策略
- 使用容器编排工具（如Kubernetes）管理服务

### 10.3 监控和日志

- 实施集中式日志收集（ELK Stack或类似方案）
- 设置系统监控（Prometheus + Grafana）
- 配置关键指标的告警机制
- 定期备份数据和配置

### 10.4 区块链节点管理

- 使用专用服务器或云服务托管区块链节点
- 实施节点监控和自动恢复机制
- 定期备份区块链数据
- 考虑使用托管的区块链服务（如Infura或AWS Managed Blockchain）