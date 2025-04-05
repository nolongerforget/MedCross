# MedCross 部署手册

## 目录

- [环境要求](#环境要求)
- [开发环境部署](#开发环境部署)
  - [前端部署](#前端部署)
  - [后端部署](#后端部署)
  - [以太坊智能合约部署](#以太坊智能合约部署)
  - [Hyperledger Fabric链码部署](#hyperledger-fabric链码部署)
- [生产环境部署](#生产环境部署)
  - [前端生产部署](#前端生产部署)
  - [后端生产部署](#后端生产部署)
  - [区块链网络生产部署](#区块链网络生产部署)
- [配置说明](#配置说明)
- [验证部署](#验证部署)
- [常见问题排查](#常见问题排查)

## 环境要求

### 通用要求

- Git
- Docker 和 Docker Compose
- Node.js 16.x 或更高版本
- Go 1.18 或更高版本
- PostgreSQL 12.0 或更高版本

### 区块链环境

- 以太坊环境
  - Ganache (开发环境)
  - Geth 或 Infura (生产环境)
  - Truffle 或 Hardhat (智能合约开发和部署)
  - Web3.js 或 ethers.js (与以太坊交互)

- Hyperledger Fabric环境
  - Hyperledger Fabric 2.2.x 或更高版本
  - Fabric SDK for Go

## 开发环境部署

### 数据库部署

1. **安装PostgreSQL**

   根据您的操作系统，按照以下步骤安装PostgreSQL：

   **Windows:**
   - 从[PostgreSQL官网](https://www.postgresql.org/download/windows/)下载安装程序
   - 运行安装程序，按照向导完成安装
   - 安装过程中设置管理员密码

   **Linux (Ubuntu/Debian):**
   ```bash
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   sudo systemctl start postgresql
   sudo systemctl enable postgresql
   ```

2. **创建数据库和用户**

   ```bash
   # 登录PostgreSQL
   # Windows
   psql -U postgres
   # Linux/macOS
   sudo -u postgres psql
   ```

   ```sql
   CREATE DATABASE medcross;
   CREATE USER medcross WITH ENCRYPTED PASSWORD 'your_secure_password';
   GRANT ALL PRIVILEGES ON DATABASE medcross TO medcross;
   ```

3. **初始化数据库结构**

   方法1：使用SQL脚本初始化
   ```bash
   cd MedCross/database
   psql -U medcross -d medcross -f init.sql
   ```

   方法2：使用应用程序自动迁移（启动后端应用时自动执行）

### 前端部署

1. **克隆代码库**

   ```bash
   git clone <repository-url>
   cd MedCross/medcrossfrontend
   ```

2. **安装依赖**

   ```bash
   npm install
   ```

3. **配置环境变量**

   创建 `.env.local` 文件，配置API端点和区块链连接信息：

   ```
   VITE_API_BASE_URL=http://localhost:8080
   VITE_ETHEREUM_NETWORK=http://localhost:7545
   VITE_CONTRACT_ADDRESS=<your-deployed-contract-address>
   ```

4. **启动开发服务器**

   ```bash
   npm run dev
   ```

   前端应用将在 `http://localhost:5173` 运行。

### 后端部署

1. **克隆代码库**

   ```bash
   git clone <repository-url>
   cd MedCross/backend
   ```

2. **安装Go依赖**

   ```bash
   go mod tidy
   ```

3. **配置后端**

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
     password: your_password
     dbname: medcross
     sslmode: disable
     max_idle_conns: 10
     max_open_conns: 100
     conn_max_lifetime: 3600

   ethereum:
     network: http://localhost:7545
     contract_address: your_contract_address
     private_key: your_private_key

   fabric:
     connection_profile: path/to/connection-profile.yaml
     channel_name: medcrosschannel
     chaincode_name: medicaldata
     org_name: Org1
     user_name: Admin
   ```

4. **启动后端服务**

   ```bash
   go run main.go
   ```

   后端API将在 `http://localhost:8080` 运行。

### 以太坊智能合约部署

1. **安装Truffle或Hardhat**

   ```bash
   npm install -g truffle
   # 或
   npm install -g hardhat
   ```

2. **启动本地以太坊网络**

   使用Ganache启动本地以太坊网络：

   ```bash
   npx ganache-cli
   # 或使用Ganache GUI应用
   ```

3. **部署智能合约**

   ```bash
   cd MedCross/contracts/ethereum
   
   # 使用Truffle
   truffle compile
   truffle migrate --network development
   
   # 或使用Hardhat
   npx hardhat compile
   npx hardhat run scripts/deploy.js --network localhost
   ```

4. **记录合约地址**

   部署完成后，记录智能合约地址，并更新前端和后端配置。

### Hyperledger Fabric链码部署

1. **设置Fabric测试网络**

   ```bash
   # 克隆Fabric示例代码
   git clone https://github.com/hyperledger/fabric-samples.git
   cd fabric-samples/test-network
   
   # 启动测试网络
   ./network.sh up createChannel -c medcrosschannel -ca
   ```

2. **部署链码**

   ```bash
   # 复制链码到Fabric示例目录
   cp -r /path/to/MedCross/contracts/fabric/medicaldata fabric-samples/chaincode/
   
   # 部署链码
   ./network.sh deployCC -c medcrosschannel -ccn medicaldata -ccp ../chaincode/medicaldata -ccl go
   ```

3. **配置连接信息**

   将生成的连接配置文件复制到后端配置目录：

   ```bash
   cp fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml /path/to/MedCross/backend/config/
   ```

## 生产环境部署

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
     network: https://mainnet.infura.io/v3/your-project-id
     # 或使用私有以太坊网络
     contract_address: your_contract_address
     private_key: your_private_key

   fabric:
     connection_profile: /etc/medcross/connection-profile.yaml
     channel_name: medcrosschannel
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

### 区块链网络生产部署

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

## 配置说明

### 前端配置参数

| 参数 | 说明 | 示例值 |
|------|------|--------|
| VITE_API_BASE_URL | 后端API基础URL | http://localhost:8080 |
| VITE_ETHEREUM_NETWORK | 以太坊网络URL | http://localhost:7545 |
| VITE_CONTRACT_ADDRESS | 以太坊合约地址 | 0x1234... |

### 后端配置参数

| 参数 | 说明 | 示例值 |
|------|------|--------|
| server.port | 服务器端口 | 8080 |
| server.mode | 运行模式 | debug/release |
| database.driver | 数据库驱动 | postgres |
| ethereum.network | 以太坊网络URL | http://localhost:7545 |
| fabric.channel_name | Fabric通道名称 | medcrosschannel |

## 验证部署

### 前端验证

1. 访问前端应用URL（开发环境：http://localhost:5173，生产环境：您的域名）
2. 验证登录功能
3. 测试数据上传和查询功能
4. 验证授权管理功能

### 后端验证

1. 测试健康检查端点：`GET /api/health`
2. 验证用户认证：`POST /api/auth/login`
3. 测试数据API：`GET /api/data`

### 区块链验证

1. **以太坊验证**
   - 使用Etherscan或区块浏览器检查合约状态
   - 通过前端应用验证数据上传和授权功能

2. **Fabric验证**
   - 使用Fabric CLI查询链码：
     ```bash
     peer chaincode query -C medcrosschannel -n medicaldata -c '{"Args":["queryAllMedicalData"]}'  
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