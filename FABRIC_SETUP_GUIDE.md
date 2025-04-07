# Hyperledger Fabric 开发环境搭建与私链部署手册

本手册详细介绍如何在Ubuntu系统上搭建Hyperledger Fabric开发环境并部署私有链网络，适用于MedCross项目的开发和测试。

## 目录

- [1. 前置依赖安装](#1-前置依赖安装)
- [2. Fabric环境安装](#2-fabric环境安装)
- [3. 测试网络部署](#3-测试网络部署)
- [4. 通道配置](#4-通道配置)
- [5. 链码部署](#5-链码部署)
- [6. 与MedCross应用集成](#6-与medcross应用集成)
- [7. 常见问题排查](#7-常见问题排查)
- [8. 生产环境部署建议](#8-生产环境部署建议)

## 1. 前置依赖安装

### 1.1 更新系统包

```bash
sudo apt-get update
sudo apt-get upgrade -y
```

### 1.2 安装Docker

```bash
# 安装必要的系统工具
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release

# 添加Docker的官方GPG密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 设置Docker稳定版仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 更新apt包索引
sudo apt-get update

# 安装最新版本的Docker Engine和containerd
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# 将当前用户添加到docker组
sudo usermod -aG docker $USER

# 验证Docker安装
docker --version
```

> 注意：添加用户到docker组后，需要注销并重新登录才能生效。

### 1.3 安装Docker Compose

```bash
# 下载Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.17.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 添加可执行权限
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

### 1.4 安装Go语言环境

```bash
# 下载Go语言安装包（以Go 1.18为例，可根据需要选择其他版本）
wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz

# 解压到/usr/local目录
sudo tar -C /usr/local -xzf go1.18.10.linux-amd64.tar.gz

# 设置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export GOPATH=$HOME/go' >> ~/.profile
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile

# 使环境变量生效
source ~/.profile

# 验证Go安装
go version
```

### 1.5 安装Node.js和npm（用于Fabric应用开发）

```bash
# 使用nvm安装Node.js
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash

# 使nvm命令可用
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# 安装Node.js LTS版本
nvm install --lts

# 验证安装
node --version
npm --version
```

### 1.6 安装其他工具

```bash
# 安装git
sudo apt-get install -y git

# 安装jq（用于处理JSON）
sudo apt-get install -y jq

# 安装Python（某些Fabric工具可能需要）
sudo apt-get install -y python3 python3-pip
```

## 2. Fabric环境安装

### 2.1 下载Fabric示例和二进制文件

```bash
# 创建工作目录
mkdir -p ~/fabric-workspace
cd ~/fabric-workspace

# 克隆Fabric示例代码
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples

# 下载Fabric二进制文件和Docker镜像
# 这里使用Fabric 2.2.0版本和CA 1.4.9版本
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.0 1.4.9
```

### 2.2 验证安装

```bash
# 检查Fabric二进制文件
ls -la ~/fabric-workspace/fabric-samples/bin

# 检查Docker镜像
docker images | grep hyperledger
```

正确安装后，应该能看到以下二进制文件：
- configtxgen
- configtxlator
- cryptogen
- discover
- idemixgen
- orderer
- peer
- fabric-ca-client

以及以下Docker镜像：
- hyperledger/fabric-ca
- hyperledger/fabric-tools
- hyperledger/fabric-peer
- hyperledger/fabric-orderer
- hyperledger/fabric-ccenv
- hyperledger/fabric-baseos

## 3. 测试网络部署

### 3.1 启动测试网络

```bash
cd ~/fabric-workspace/fabric-samples/test-network

# 确保没有运行中的网络
./network.sh down

# 启动网络并创建通道
./network.sh up createChannel -c medcrosschannel -ca
```

这个命令会：
1. 启动一个包含两个组织（Org1和Org2）的Fabric网络
2. 每个组织有一个peer节点
3. 启动一个排序服务节点
4. 创建一个名为medcrosschannel的通道
5. 两个组织的peer节点都会加入该通道

### 3.2 验证网络状态

```bash
# 查看运行中的Docker容器
docker ps
```

应该能看到以下容器运行：
- peer0.org1.example.com
- peer0.org2.example.com
- orderer.example.com
- ca_org1
- ca_org2
- ca_orderer

## 4. 通道配置

### 4.1 设置环境变量

为了与Org1的peer节点交互，需要设置以下环境变量：

```bash
cd ~/fabric-workspace/fabric-samples/test-network

# 设置Org1的环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

### 4.2 查看通道信息

```bash
# 获取通道信息
peer channel getinfo -c medcrosschannel
```

## 5. 链码部署

### 5.1 准备MedCross链码

```bash
# 复制MedCross链码到Fabric示例目录
cp -r /path/to/MedCross/contracts/fabric/medicaldata ~/fabric-workspace/fabric-samples/chaincode/
```

> 注意：请将`/path/to/MedCross`替换为实际的MedCross项目路径。

### 5.2 部署链码

```bash
cd ~/fabric-workspace/fabric-samples/test-network

# 部署链码
./network.sh deployCC -c medcrosschannel -ccn medicaldata -ccp ../chaincode/medicaldata -ccl go
```

这个命令会：
1. 打包链码
2. 在Org1和Org2的peer节点上安装链码
3. 批准链码定义
4. 提交链码定义到通道
5. 初始化链码

### 5.3 测试链码

```bash
# 调用链码初始化函数
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C medcrosschannel -n medicaldata --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

# 查询链码
peer chaincode query -C medcrosschannel -n medicaldata -c '{"Args":["GetAllMedicalData"]}'
```

## 6. 与MedCross应用集成

### 6.1 配置连接信息

```bash
# 创建配置目录
mkdir -p /path/to/MedCross/backend/fabric-config

# 复制连接配置文件
cp ~/fabric-workspace/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml /path/to/MedCross/backend/fabric-config/
```

> 注意：请将`/path/to/MedCross`替换为实际的MedCross项目路径。

### 6.2 配置Fabric网关服务

编辑MedCross后端的Fabric服务配置文件，确保正确设置了连接信息：

```yaml
fabric:
  connection_profile: "./fabric-config/connection-org1.yaml"
  channel_name: "medcrosschannel"
  chaincode_name: "medicaldata"
  msp_id: "Org1MSP"
  wallet_path: "./fabric-config/wallet"
  user_id: "appUser"
```

### 6.3 注册应用用户

```bash
# 设置环境变量
export PATH=${PWD}/../bin:$PATH
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

# 注册应用用户
fabric-ca-client register --caname ca-org1 --id.name appUser --id.secret appUserPW --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem

# 生成用户证书和密钥
fabric-ca-client enroll -u https://appUser:appUserPW@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/appUser@org1.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem

# 复制管理员证书到用户的MSP目录
mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/users/appUser@org1.example.com/msp/admincerts
cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem ${PWD}/organizations/peerOrganizations/org1.example.com/users/appUser@org1.example.com/msp/admincerts/
```

## 7. 常见问题排查

### 7.1 Docker相关问题

**问题：无法启动Docker容器**

```bash
# 检查Docker服务状态
sudo systemctl status docker

# 如果服务未运行，启动服务
sudo systemctl start docker

# 检查Docker日志
sudo journalctl -u docker
```

**问题：Docker权限问题**

```bash
# 确保当前用户在docker组中
groups

# 如果不在，添加用户到docker组并重新登录
sudo usermod -aG docker $USER
newgrp docker
```

### 7.2 Fabric网络问题

**问题：网络启动失败**

```bash
# 清理网络并重新启动
cd ~/fabric-workspace/fabric-samples/test-network
./network.sh down
docker system prune -a
./network.sh up createChannel -c medcrosschannel -ca
```

**问题：链码部署失败**

```bash
# 检查链码日志
docker logs $(docker ps -q --filter name=dev-peer0.org1.example.com-medicaldata)

# 检查peer节点日志
docker logs peer0.org1.example.com
```

### 7.3 Go模块问题

**问题：Go模块下载失败**

```bash
# 设置GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 清理Go模块缓存
go clean -modcache
```

## 8. 生产环境部署建议

### 8.1 网络规划

生产环境中，建议：

1. 至少部署3个组织，每个组织至少2个peer节点
2. 排序服务使用Raft共识，至少部署5个排序节点
3. 每个组织部署独立的CA服务器
4. 使用TLS加密所有通信
5. 配置适当的背书策略

### 8.2 硬件要求

每个节点的最低硬件配置：

- CPU: 4核
- 内存: 8GB
- 磁盘: 100GB SSD
- 网络: 1Gbps

### 8.3 安全建议

1. 使用防火墙限制网络访问
2. 定期更新系统和Docker
3. 使用私有Docker Registry
4. 实施证书轮换策略
5. 监控系统和网络活动
6. 定期备份账本和证书

### 8.4 生产部署步骤

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

> 注意：生产环境部署需要根据实际网络拓扑和组织结构进行调整。

---

本手册提供了在Ubuntu系统上搭建Hyperledger Fabric开发环境和部署私有链的完整指南。如有问题，请参考[Hyperledger Fabric官方文档](https://hyperledger-fabric.readthedocs.io/)或联系MedCross技术支持团队。