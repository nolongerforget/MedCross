# MedCross 环境配置说明文档

## 目录

- [概述](#概述)
- [环境要求](#环境要求)
- [数据库部署](#数据库部署)
- [后端部署](#后端部署)
- [前端部署](#前端部署)
- [区块链环境配置](#区块链环境配置)
  - [以太坊私链部署](#以太坊私链部署)
  - [Hyperledger Fabric部署](#hyperledger-fabric部署)
- [系统集成](#系统集成)
- [验证部署](#验证部署)
- [生产环境配置](#生产环境配置)
- [常见问题排查](#常见问题排查)
- [安全建议](#安全建议)
- [参考资源](#参考资源)

## 概述

MedCross是一个医疗数据跨链共享平台，旨在通过区块链技术实现医疗数据的安全共享和授权管理。本文档提供了完整的环境配置和部署指南，包括所有组件的安装、配置和验证步骤。

系统架构包括：
- 前端：基于React的用户界面
- 后端：基于Go+Gin的API服务
- 数据库：PostgreSQL关系型数据库
- 区块链：以太坊私链和Hyperledger Fabric私链

## 环境要求

### 通用要求

- Git
- Docker 和 Docker Compose
- Node.js 16.x 或更高版本
- Go 1.18 或更高版本
- PostgreSQL 12.0 或更高版本

### 硬件推荐配置

- CPU: 4核或更高
- 内存: 至少8GB RAM
- 存储: 至少50GB可用空间
- 网络: 稳定的互联网连接

### 区块链环境要求

- 以太坊环境
  - Ganache (开发环境)
  - Geth 或 Infura (生产环境)
  - Truffle 或 Hardhat (智能合约开发和部署)
  - Web3.js 或 ethers.js (与以太坊交互)

- Hyperledger Fabric环境
  - Hyperledger Fabric 2.2.x 或更高版本
  - Fabric SDK for Go

## 数据库部署

### 安装PostgreSQL

#### Windows

1. 从[PostgreSQL官网](https://www.postgresql.org/download/windows/)下载安装程序
2. 运行安装程序，按照向导完成安装
3. 安装过程中设置管理员密码
4. 默认端口为5432，如需更改请记录新端口号

#### Linux (Ubuntu/Debian)

```bash
# 安装PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### macOS

```bash
# 使用Homebrew安装
brew install postgresql

# 启动服务
brew services start postgresql
```

### 创建数据库和用户

1. 登录PostgreSQL

```bash
# Linux/macOS
sudo -u postgres psql

# Windows (使用psql命令行工具)
psql -U postgres
```

2. 创建数据库和用户

```sql
CREATE DATABASE medcross;
CREATE USER medcross WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE medcross TO medcross;
```

3. 连接到新创建的数据库

```sql
\c medcross
```

### 初始化数据库结构

#### 方法1：使用SQL脚本初始化

1. 导航到项目的database目录

```bash
cd /path/to/MedCross/database
```

2. 使用psql执行初始化脚本

```bash
psql -U medcross -d medcross -f init.sql
```

初始化脚本将创建以下表：
- users: 用户信息
- medical_data: 医疗数据元数据
- tags: 数据标签
- data_tags: 数据与标签的关联
- authorizations: 数据授权记录
- access_records: 数据访问记录
- cross_chain_references: 跨链引用信息

#### 方法2：使用应用程序自动迁移

1. 确保在`config/config.yaml`中正确配置数据库连接信息

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: medcross
  password: your_secure_password
  dbname: medcross
  sslmode: disable
```

2. 启动后端应用程序，它将自动执行数据库迁移

```bash
cd /path/to/MedCross/backend
go run main.go
```

## 后端部署

### 安装Go环境

1. 从[Go官网](https://golang.org/dl/)下载适合您操作系统的安装包
2. 按照安装向导完成安装
3. 验证安装

```bash
go version
```

### 获取项目代码

```bash
git clone <repository-url>
cd MedCross/backend
```

### 安装依赖

```bash
go mod tidy
```

### 配置后端

编辑 `config/config.yaml` 文件，配置数据库连接、区块链连接和API设置：

```yaml
server:
  port: 8080
  mode: debug  # 开发环境使用debug模式

database:
  driver: postgres
  host: localhost
  port: 5432
  user: medcross
  password: your_secure_password
  dbname: medcross
  sslmode: disable
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

ethereum:
  rpc_url: http://localhost:8545
  contract_address: your_contract_address
  private_key: your_private_key

fabric:
  config_path: ./fabric-config
  channel_id: mychannel
  chaincode_name: medicaldata
  user_name: Admin
  org_name: Org1

jwt:
  secret: medcross_secret_key
  expires_in: 86400
```

### 启动后端服务

```bash
go run main.go
```

后端API将在 `http://localhost:8080` 运行。

## 前端部署

### 安装Node.js和npm

1. 从[Node.js官网](https://nodejs.org/)下载并安装Node.js
2. 验证安装

```bash
node -v
npm -v
```

### 获取前端代码

```bash
git clone <repository-url>
cd MedCross/medcrossfrontend
```

### 安装依赖

```bash
npm install
```

### 配置环境变量

创建 `.env.local` 文件，配置API端点和区块链连接信息：

```
VITE_API_BASE_URL=http://localhost:8080
VITE_ETHEREUM_NETWORK=http://localhost:7545
VITE_CONTRACT_ADDRESS=<your-deployed-contract-address>
```

### 启动开发服务器

```bash
npm run dev
```

前端应用将在 `http://localhost:5173` 运行。

## 区块链环境配置

### 以太坊私链部署

#### 安装Ganache (开发环境)

1. 安装Ganache

```bash
npm install -g ganache-cli
# 或下载Ganache GUI应用: https://www.trufflesuite.com/ganache
```

2. 启动本地以太坊网络

```bash
ganache-cli
# 或启动Ganache GUI应用
```

#### 安装Truffle或Hardhat

```bash
npm install -g truffle
# 或
npm install -g hardhat
```

#### 部署智能合约

1. 导航到以太坊合约目录

```bash
cd MedCross/contracts/ethereum
```

2. 编译和部署合约

```bash
# 使用Truffle
truffle compile
truffle migrate --network development

# 或使用Hardhat
npx hardhat compile
npx hardhat run scripts/deploy.js --network localhost
```

3. 记录合约地址

部署完成后，记录智能合约地址，并更新前端和后端配置。

### Hyperledger Fabric部署

#### 设置Fabric测试网络

1. 克隆Fabric示例代码

```bash
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples
```

2. 安装Fabric二进制文件和Docker镜像

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.0 1.4.9
```

3. 启动测试网络

```bash
cd test-network
./network.sh up createChannel -c medcrosschannel -ca
```

#### 部署链码

1. 复制链码到Fabric示例目录

```bash
cp -r /path/to/MedCross/contracts/fabric/medicaldata fabric-samples/chaincode/
```

2. 部署链码

```bash
./network.sh deployCC -c medcrosschannel -ccn medicaldata -ccp ../chaincode/medicaldata -ccl go
```

3. 配置连接信息

将生成的连接配置文件复制到后端配置目录：

```bash
mkdir -p /path/to/MedCross/backend/fabric-config
cp fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml /path/to/MedCross/backend/fabric-config/
```

## 系统集成

完成各组件部署后，需要进行系统集成，确保各部分能够正常协同工作。

### 更新配置文件

1. 更新后端配置

编辑 `backend/config/config.yaml`，确保以下配置正确：

```yaml
ethereum:
  rpc_url: http://localhost:7545  # Ganache默认端口
  contract_address: "0x..."  # 部署的合约地址

fabric:
  config_path: "./fabric-config/connection-org1.yaml"
  channel_id: "medcrosschannel"
  chaincode_name: "medicaldata"
```

2. 更新前端环境变量

编辑 `.env.local` 文件：

```
VITE_CONTRACT_ADDRESS=0x...  # 部署的合约地址
```

### 重启服务

1. 重启后端服务

```bash
cd /path/to/MedCross/backend
go run main.go
```

2. 重启前端服务

```bash
cd /path/to/MedCross/medcrossfrontend
npm run dev
```

## 验证部署

### 前端验证

1. 访问前端应用URL：http://localhost:5173
2. 验证登录功能
3. 测试数据上传和查询功能
4. 验证授权管理功能

### 后端验证

1. 测试健康检查端点：`GET http://localhost:8080/api/health`
2. 验证用户认证：`POST http://localhost:8080/api/auth/login`
3. 测试数据API：`GET http://localhost:8080/api/data`

### 区块链验证

1. **以太坊验证**
   - 使用Ganache界面检查合约状态
   - 通过前端应用验证数据上传和授权功能

2. **Fabric验证**
   - 使用Fabric CLI查询链码：
     ```bash
     cd fabric-samples/test-network
     export PATH=${PWD}/../bin:$PATH
     export FABRIC_CFG_PATH=$PWD/../config/
     export CORE_PEER_TLS_ENABLED=true
     export CORE_PEER_LOCALMSPID="Org1MSP"
     export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
     export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
     export CORE_PEER_ADDRESS=localhost:7051
     
     peer chaincode query -C medcrosschannel -n medicaldata -c '{"Args":["queryAllMedicalData"]}'  
     ```

## 生产环境配置

### 数据库生产配置

1. **启用SSL连接**

   修改`config/config.yaml`中的`sslmode`为`require`：

   ```yaml
   database:
     sslmode: require
   ```

2. **配置连接池**

   ```yaml
   database:
     max_idle_conns: 10     # 空闲连接池中连接的最大数量
     max_open_conns: 100    # 打开数据库连接的最大数量
     conn_max_lifetime: 3600 # 连接可复用的最大时间（秒）
   ```

3. **定期备份**

   设置定时任务，每天自动备份数据库：

   ```bash
   # 添加到crontab
   0 2 * * * pg_dump -U medcross -d medcross -F c -f /backup/medcross_$(date +\%Y\%m\%d).dump
   ```

### 前端生产部署

1. **构建前端应用**

   ```bash
   cd MedCross/medcrossfrontend
   npm run build
   ```

2. **部署到Web服务器**

   将 `dist` 目录中的文件部署到Nginx、Apache或其他Web服务器：

   **Nginx配置示例：**

   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       root /path/to/MedCross/medcrossfrontend/dist;
       index index.html;
       
       location / {
           try_files $uri $uri/ /index.html;
       }
       
       location /api/ {
           proxy_pass http://backend-server:8080/;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

### 后端生产部署

1. **构建后端应用**

   ```bash
   cd MedCross/backend
   go build -o medcross-api
   ```

2. **配置生产环境**

   创建生产环境配置文件 `config/config.prod.yaml`：

   ```yaml
   server:
     port: 8080
     mode: release  # 生产环境使用release模式

   database:
     driver: postgres
     host: your-db-host
     port: 5432
     user: medcross
     password: strong_password
     dbname: medcross
     sslmode: require

   ethereum:
     rpc_url: https://mainnet.infura.io/v3/your-project-id
     # 或使用私有以太坊网络
     contract_address: your_contract_address
     private_key: your_private_key

   fabric:
     config_path: /etc/medcross/connection-profile.yaml
     channel_id: medcrosschannel
     chaincode_name: medicaldata
     org_name: Org1
     user_name: Admin
   ```

3. **使用Docker部署**

   创建 `Dockerfile`：

   ```dockerfile
   FROM golang:1.18-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go mod download
   RUN go build -o medcross-api

   FROM alpine:latest
   WORKDIR /app
   COPY --from=builder /app/medcross-api .
   COPY --from=builder /app/config ./config
   EXPOSE 8080
   CMD ["./medcross-api"]
   ```

   构建和运行Docker容器：

   ```bash
   docker build -t medcross-api .
   docker run -d -p 8080:8080 -v /path/to/config:/app/config --name medcross-api medcross-api
   ```

### 区块链生产部署

#### 以太坊网络

1. **选择以太坊网络**

   - 使用公共以太坊网络（主网或测试网）
   - 或部署私有以太坊网络（使用Geth或Quorum）

2. **部署智能合约到生产网络**

   ```bash
   # 使用Truffle
   truffle migrate --network mainnet
   
   # 或使用Hardhat
   npx hardhat run scripts/deploy.js --network mainnet
   ```

#### Hyperledger Fabric网络

1. **设置生产级Fabric网络**

   按照[Hyperledger Fabric文档](https://hyperledger-fabric.readthedocs.io/en/release-2.2/deployment_guide_overview.html)设置多组织Fabric网络。

2. **部署链码到生产网络**

   ```bash
   # 使用Fabric工具部署链码
   peer lifecycle chaincode package medicaldata.tar.gz --path /path/to/MedCross/contracts/fabric/medicaldata --lang golang --label medicaldata_1.0
   
   # 安装链码包
   peer lifecycle chaincode install medicaldata.tar.gz
   
   # 批准链码定义
   peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --channelID medcrosschannel --name medicaldata --version 1.0 --package-id <package-id> --sequence 1
   
   # 提交链码定义
   peer lifecycle chaincode commit -o orderer.example.com:7050 --channelID medcrosschannel --name medicaldata --version 1.0 --sequence 1
   ```

## 常见问题排查

### 前端问题

1. **无法连接到后端API**
   - 检查API基础URL配置
   - 验证后端服务是否正在运行
   - 检查CORS配置

2. **以太坊连接问题**
   - 确认以太坊网络URL正确
   - 检查合约地址配置
   - 验证MetaMask或Web3提供者配置

### 后端问题

1. **数据库连接失败**
   - 验证数据库凭据和连接字符串
   - 检查数据库服务是否运行
   - 确认网络连接和防火墙设置

2. **区块链交互错误**
   - 检查以太坊网络配置和私钥
   - 验证Fabric连接配置和证书
   - 确认链码已正确部署

### 区块链问题

1. **以太坊智能合约问题**
   - 检查合约部署交易是否成功
   - 验证合约ABI是否正确
   - 确认账户有足够的ETH支付Gas费用

2. **Fabric链码问题**
   - 检查链码安装和实例化日志
   - 验证链码背书策略配置
   - 确认组织MSP配置正确

## 安全建议

1. **数据库安全**
   - 使用强密码保护数据库账户
   - 限制数据库服务器的网络访问
   - 定期更新PostgreSQL到最新版本
   - 审核数据库访问日志
   - 实施数据加密（特别是敏感医疗数据）

2. **API安全**
   - 使用HTTPS加密所有API通信
   - 实施适当的认证和授权机制
   - 使用安全的JWT配置
   - 限制API请求速率

3. **区块链安全**
   - 安全存储私钥
   - 定期审核智能合约
   - 使用多重签名机制
   - 实施访问控制列表

4. **服务器安全**
   - 定期更新操作系统和软件包
   - 配置防火墙限制访问
   - 使用入侵检测系统
   - 定期备份数据

## 参考资源

- [PostgreSQL官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL安全最佳实践](https://www.postgresql.org/docs/current/security.html)
- [Go语言官方文档](https://golang.org/doc/)
- [Gin Web框架文档](https://gin-gonic.com/docs/)
- [React官方文档](https://reactjs.org/docs/getting-started.html)
- [Vite构建工具文档](https://vitejs.dev/guide/)
- [以太坊开发者文档](https://ethereum.org/developers/)
- [Truffle框架文档](https://www.trufflesuite.com/docs/truffle/overview)
- [Hyperledger Fabric文档](https://hyperledger-fabric.readthedocs.io/)
- [Docker文档](https://docs.docker.com/)