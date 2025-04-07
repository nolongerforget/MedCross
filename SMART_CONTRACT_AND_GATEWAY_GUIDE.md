# MedCross 智能合约与跨链网关开发细则

## 1. 概述

### 1.1 文档目的

本文档详细说明MedCross平台中智能合约和跨链网关的设计、实现和部署细则，为开发团队提供技术指导。

### 1.2 系统架构

MedCross平台采用双链架构，集成以太坊和Hyperledger Fabric两条区块链，通过跨链网关实现数据的互通和查询。系统核心组件包括：

- **以太坊智能合约**：基于Solidity开发，负责在以太坊网络上管理医疗数据
- **Fabric链码**：基于Go语言开发，负责在Fabric网络上管理医疗数据
- **跨链网关**：协调两条链上的数据交互，实现跨链数据查询和转移

### 1.3 技术栈

- **以太坊智能合约**：Solidity 0.8.0+
- **Fabric链码**：Go 1.18+, Hyperledger Fabric Contract API
- **跨链网关**：Go 1.18+, Gin框架

## 2. 以太坊智能合约

### 2.1 合约设计

#### 2.1.1 数据结构

```solidity
struct Data {
    uint256 id;
    address owner;
    string dataHash;      // IPFS或其他存储系统的哈希值
    string dataType;      // 数据类型（如：影像数据、电子病历等）
    string metadata;      // JSON格式的元数据
    uint256 timestamp;    // 上传时间戳
    string keywords;      // 关键词，用于搜索，以逗号分隔
}
```

#### 2.1.2 存储结构

```solidity
// 存储所有医疗数据
Data[] private allData;

// 用户拥有的数据映射
mapping(address => uint256[]) private userDataIds;

// 数据类型到数据ID的映射
mapping(string => uint256[]) private typeToDataIds;
```

#### 2.1.3 事件定义

```solidity
event DataUploaded(uint256 indexed id, address indexed owner, string dataType, uint256 timestamp);
```

### 2.2 核心功能实现

#### 2.2.1 数据上传

```solidity
function uploadData(
    string memory dataHash,
    string memory dataType,
    string memory metadata,
    string memory keywords
) public returns (uint256) {
    uint256 newId = allData.length;
    
    Data memory newData = Data({
        id: newId,
        owner: msg.sender,
        dataHash: dataHash,
        dataType: dataType,
        metadata: metadata,
        timestamp: block.timestamp,
        keywords: keywords
    });
    
    allData.push(newData);
    userDataIds[msg.sender].push(newId);
    typeToDataIds[dataType].push(newId);
    
    emit DataUploaded(newId, msg.sender, dataType, block.timestamp);
    
    return newId;
}
```

#### 2.2.2 数据查询

```solidity
function getData(uint256 id) public view returns (
    uint256,
    address,
    string memory,
    string memory,
    string memory,
    uint256,
    string memory
) {
    require(id < allData.length, "Data does not exist");
    Data memory data = allData[id];
    
    return (
        data.id,
        data.owner,
        data.dataHash,
        data.dataType,
        data.metadata,
        data.timestamp,
        data.keywords
    );
}
```

#### 2.2.3 用户数据查询

```solidity
function getUserDataIds(address user) public view returns (uint256[] memory) {
    return userDataIds[user];
}
```

### 2.3 安全考虑

#### 2.3.1 访问控制

- 实现基于角色的访问控制（RBAC）
- 添加数据所有权验证
- 实现数据共享权限管理

```solidity
// 示例：添加访问控制修饰符
modifier onlyOwner(uint256 dataId) {
    require(dataId < allData.length, "Data does not exist");
    require(allData[dataId].owner == msg.sender, "Not the owner");
    _;
}
```

#### 2.3.2 数据隐私

- 敏感数据应存储在链下，链上只保存哈希值
- 使用加密技术保护元数据
- 考虑使用零知识证明等隐私保护技术

#### 2.3.3 防御攻击

- 防止重入攻击
- 防止整数溢出
- 防止拒绝服务攻击

### 2.4 优化策略

#### 2.4.1 Gas优化

- 减少存储操作
- 优化数据结构
- 批量处理操作

#### 2.4.2 查询优化

- 实现高效的索引结构
- 优化关键词搜索算法
- 实现分页查询

## 3. Hyperledger Fabric链码

### 3.1 链码设计

#### 3.1.1 数据结构

```go
// MedicalRecord 结构定义医疗数据记录
type MedicalRecord struct {
    ID        string    `json:"id"`
    Owner     string    `json:"owner"`
    DataHash  string    `json:"dataHash"`  // IPFS或其他存储系统的哈希值
    DataType  string    `json:"dataType"`  // 数据类型（如：影像数据、电子病历等）
    Metadata  string    `json:"metadata"`  // JSON格式的元数据
    Timestamp time.Time `json:"timestamp"` // 上传时间戳
    Keywords  string    `json:"keywords"`  // 关键词，用于搜索，以逗号分隔
}
```

#### 3.1.2 索引设计

使用Fabric的复合键创建索引：

- 按所有者查询：`owner~id`
- 按数据类型查询：`type~id`
- 按关键词查询：`keyword~id`

### 3.2 核心功能实现

#### 3.2.1 数据上传

```go
// UploadData 上传新的医疗数据
func (s *MedicalData) UploadData(ctx contractapi.TransactionContextInterface, id string, owner string, dataHash string, dataType string, metadata string, keywords string) error {
    // 检查数据是否已存在
    exists, err := s.DataExists(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to check if data exists: %v", err)
    }
    if exists {
        return fmt.Errorf("data already exists: %s", id)
    }

    // 创建新的医疗数据记录
    record := MedicalRecord{
        ID:        id,
        Owner:     owner,
        DataHash:  dataHash,
        DataType:  dataType,
        Metadata:  metadata,
        Timestamp: time.Now(),
        Keywords:  keywords,
    }

    // 将数据转换为JSON并存储
    recordJSON, err := json.Marshal(record)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %v", err)
    }

    // 将数据写入账本
    err = ctx.GetStub().PutState(id, recordJSON)
    if err != nil {
        return fmt.Errorf("failed to put data in world state: %v", err)
    }

    // 创建复合键用于按所有者查询
    ownerCompositeKey, err := ctx.GetStub().CreateCompositeKey("owner~id", []string{owner, id})
    if err != nil {
        return fmt.Errorf("failed to create composite key: %v", err)
    }

    // 存储复合键
    err = ctx.GetStub().PutState(ownerCompositeKey, []byte{0})
    if err != nil {
        return fmt.Errorf("failed to put owner composite key: %v", err)
    }

    // 创建复合键用于按类型查询
    typeCompositeKey, err := ctx.GetStub().CreateCompositeKey("type~id", []string{dataType, id})
    if err != nil {
        return fmt.Errorf("failed to create composite key: %v", err)
    }

    // 存储复合键
    err = ctx.GetStub().PutState(typeCompositeKey, []byte{0})
    if err != nil {
        return fmt.Errorf("failed to put type composite key: %v", err)
    }

    return nil
}
```

#### 3.2.2 数据查询

```go
// GetData 根据ID获取医疗数据
func (s *MedicalData) GetData(ctx contractapi.TransactionContextInterface, id string) (*MedicalRecord, error) {
    // 从账本中获取数据
    recordJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read data from world state: %v", err)
    }
    if recordJSON == nil {
        return nil, fmt.Errorf("data does not exist: %s", id)
    }

    // 将JSON转换为结构体
    var record MedicalRecord
    err = json.Unmarshal(recordJSON, &record)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal data: %v", err)
    }

    return &record, nil
}
```

#### 3.2.3 按所有者查询

```go
// GetDataByOwner 获取特定所有者的所有医疗数据
func (s *MedicalData) GetDataByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*MedicalRecord, error) {
    // 创建复合键迭代器
    iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("owner~id", []string{owner})
    if err != nil {
        return nil, fmt.Errorf("failed to get iterator: %v", err)
    }
    defer iterator.Close()

    var records []*MedicalRecord

    // 遍历结果
    for iterator.HasNext() {
        queryResponse, err := iterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to iterate: %v", err)
        }

        // 从复合键中提取ID
        _, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
        if err != nil {
            return nil, fmt.Errorf("failed to split composite key: %v", err)
        }

        if len(compositeKeyParts) < 2 {
            continue
        }

        id := compositeKeyParts[1]

        // 获取数据记录
        record, err := s.GetData(ctx, id)
        if err != nil {
            // 记录错误但继续处理
            log.Printf("Error getting data %s: %v", id, err)
            continue
        }

        records = append(records, record)
    }

    return records, nil
}
```

### 3.3 安全考虑

#### 3.3.1 访问控制

- 利用Fabric的MSP（成员服务提供商）进行身份验证
- 使用基于属性的访问控制（ABAC）
- 实现细粒度的数据访问策略

#### 3.3.2 数据隐私

- 使用私有数据集合存储敏感信息
- 实现基于通道的数据隔离
- 使用加密技术保护元数据

#### 3.3.3 防御攻击

- 输入验证和清理
- 防止SQL注入和其他注入攻击
- 实现适当的错误处理和日志记录

### 3.4 优化策略

#### 3.4.1 性能优化

- 优化查询设计
- 实现高效的索引结构
- 使用分页查询处理大量数据

#### 3.4.2 存储优化

- 仅在链上存储必要数据
- 使用链外存储（如IPFS）存储大型文件
- 实现数据归档策略

## 4. 跨链网关

### 4.1 网关架构

#### 4.1.1 组件结构

```
crosschain-gateway/
├── main.go                 # 主程序入口
├── config/                 # 配置文件
├── api/                    # API定义
├── handlers/               # 请求处理器
├── blockchain/             # 区块链接口
│   ├── ethereum/           # 以太坊接口
│   └── fabric/             # Fabric接口
├── converter/              # 数据转换器
└── utils/                  # 工具函数
```

#### 4.1.2 服务接口

```go
// GatewayService 跨链网关服务
type GatewayService struct {
    gatewayURL string
    timeout    time.Duration // HTTP请求超时时间
    maxRetries int           // 最大重试次数
}
```

### 4.2 数据转换机制

#### 4.2.1 转换器设计

```go
// ChainConverter 区块链数据转换工具
type ChainConverter struct{}

// EthereumToFabric 将以太坊格式的医疗数据转换为Fabric格式
func (c *ChainConverter) EthereumToFabric(ethData models.MedicalData) (models.MedicalData, error) {
    // 复制基本数据
    fabricData := ethData

    // 修改链标识
    fabricData.Chain = "fabric"

    // 转换所有者格式 (从以太坊地址转为Fabric身份)
    fabricData.Owner = convertEthAddressToFabricID(ethData.Owner)

    // 转换元数据格式
    metadata, err := c.convertMetadata(ethData.Metadata, "ethereum", "fabric")
    if err != nil {
        return models.MedicalData{}, err
    }
    fabricData.Metadata = metadata

    // 生成新的Fabric兼容ID (保留原始ID的引用)
    fabricData.ID = fmt.Sprintf("fab-%s", strings.TrimPrefix(ethData.ID, "eth-"))

    return fabricData, nil
}

// FabricToEthereum 将Fabric格式的医疗数据转换为以太坊格式
func (c *ChainConverter) FabricToEthereum(fabricData models.MedicalData) (models.MedicalData, error) {
    // 复制基本数据
    ethData := fabricData

    // 修改链标识
    ethData.Chain = "ethereum"

    // 转换所有者格式 (从Fabric身份转为以太坊地址)
    ethData.Owner = convertFabricIDToEthAddress(fabricData.Owner)

    // 转换元数据格式
    metadata, err := c.convertMetadata(fabricData.Metadata, "fabric", "ethereum")
    if err != nil {
        return models.MedicalData{}, err
    }
    ethData.Metadata = metadata

    // 生成新的以太坊兼容ID (保留原始ID的引用)
    ethData.ID = fmt.Sprintf("eth-%s", strings.TrimPrefix(fabricData.ID, "fab-"))

    return ethData, nil
}
```

#### 4.2.2 元数据转换

```go
// 转换元数据格式
func (c *ChainConverter) convertMetadata(metadataJSON, sourceChain, targetChain string) (string, error) {
    // 解析原始元数据
    var metadata map[string]interface{}
    if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
        log.Printf("解析元数据失败: %v", err)
        return "", err
    }

    // 添加跨链转换信息
    metadata["crossChainSource"] = sourceChain
    metadata["crossChainTimestamp"] = time.Now().Format(time.RFC3339)

    // 根据目标链调整格式
    if targetChain == "ethereum" {
        // 以太坊特定的元数据调整
        if _, ok := metadata["patientId"]; ok {
            metadata["patientAddress"] = "0x" + fmt.Sprintf("%x", metadata["patientId"])
        }
    } else if targetChain == "fabric" {
        // Fabric特定的元数据调整
        if _, ok := metadata["patientAddress"]; ok {
            // 移除0x前缀
            address := metadata["patientAddress"].(string)
            if strings.HasPrefix(address, "0x") {
                metadata["patientId"] = strings.TrimPrefix(address, "0x")
            }
        }
    }

    // 序列化调整后的元数据
    result, err := json.Marshal(metadata)
    if err != nil {
        return "", err
    }

    return string(result), nil
}
```

### 4.3 跨链查询流程

#### 4.3.1 查询接口

```go
// QueryData 查询医疗数据
func (s *GatewayService) QueryData(query models.MedicalDataQuery) (*models.QueryResult, error) {
    // 构建查询URL
    url := fmt.Sprintf("%s/query?keyword=%s&dataType=%s&chain=%s",
        s.gatewayURL,
        query.Keyword,
        query.DataType,
        query.Chain)

    // 添加日期范围
    if query.StartDate != "" {
        url += "&startDate=" + query.StartDate
    }
    if query.EndDate != "" {
        url += "&endDate=" + query.EndDate
    }

    // 添加排序和分页
    url += fmt.Sprintf("&sortBy=%s&page=%d&pageSize=%d",
        query.SortBy,
        query.Page,
        query.PageSize)

    // 创建带超时的HTTP客户端
    client := &http.Client{
        Timeout: s.timeout,
    }

    // 添加重试机制
    var resp *http.Response
    var err error
    for i := 0; i < s.maxRetries; i++ {
        log.Printf("尝试查询网关 (尝试 %d/%d): %s", i+1, s.maxRetries, url)
        resp, err = client.Get(url)
        if err == nil {
            break
        }
        log.Printf("查询网关失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
        if i < s.maxRetries-1 {
            // 指数退避策略
            backoff := time.Duration(100*(i+1)) * time.Millisecond
            time.Sleep(backoff)
        }
    }

    if err != nil {
        return nil, fmt.Errorf("查询网关失败: %v", err)
    }
    defer resp.Body.Close()

    // 解析响应
    var result models.QueryResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("解析响应失败: %v", err)
    }

    return &result, nil
}
```

#### 4.3.2 并行查询实现

```go
// QueryAllChains 并行查询所有区块链
func (s *GatewayService) QueryAllChains(query models.MedicalDataQuery) (*models.QueryResult, error) {
    // 创建通道接收查询结果
    ethResultCh := make(chan struct {
        result *models.QueryResult
        err    error
    })
    fabricResultCh := make(chan struct {
        result *models.QueryResult
        err    error
    })

    // 并行查询以太坊
    go func() {
        ethQuery := query
        ethQuery.Chain = "ethereum"
        result, err := s.QueryData(ethQuery)
        ethResultCh <- struct {
            result *models.QueryResult
            err    error
        }{result, err}
    }()

    // 并行查询Fabric
    go func() {
        fabricQuery := query
        fabricQuery.Chain = "fabric"
        result, err := s.QueryData(fabricQuery)
        fabricResultCh <- struct {
            result *models.QueryResult
            err    error
        }{result, err}
    }()

    // 接收以太坊结果
    ethResp := <-ethResultCh
    // 接收Fabric结果
    fabricResp := <-fabricResultCh

    // 合并结果
    mergedResult := &models.QueryResult{
        TotalCount: 0,
        Data:       []models.MedicalData{},
        Errors:     []string{},
    }

    // 处理以太坊结果
    if ethResp.err != nil {
        mergedResult.Errors = append(mergedResult.Errors, fmt.Sprintf("以太坊查询错误: %v", ethResp.err))
    } else if ethResp.result != nil {
        mergedResult.Data = append(mergedResult.Data, ethResp.result.Data...)
        mergedResult.TotalCount += ethResp.result.TotalCount
        mergedResult.Errors = append(mergedResult.Errors, ethResp.result.Errors...)
    }

    // 处理Fabric结果
    if fabricResp.err != nil {
        mergedResult.Errors = append(mergedResult.Errors, fmt.Sprintf("Fabric查询错误: %v", fabricResp.err))
    } else if fabricResp.result != nil {
        mergedResult.Data = append(mergedResult.Data, fabricResp.result.Data...)
        mergedResult.TotalCount += fabricResp.result.TotalCount
        mergedResult.Errors = append(mergedResult.Errors, fabricResp.result.Errors...)
    }

    return mergedResult, nil
}
```

### 4.4 数据转移流程

#### 4.4.1 转移接口

```go
// TransferData 跨链转移数据
func (s *GatewayService) TransferData(transfer models.TransferRequest) (*models.TransferResult, error) {
    // 构建转移请求URL
    url := fmt.Sprintf("%s/transfer", s.gatewayURL)

    // 序列化请求体
    requestBody, err := json.Marshal(transfer)
    if err != nil {
        return nil, fmt.Errorf("序列化请求失败: %v", err)
    }

    // 创建带超时的HTTP客户端
    client := &http.Client{
        Timeout: s.timeout * 2, // 转移操作可能需要更长时间
    }

    // 发送POST请求
    resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("转移请求失败: %v", err)
    }
    defer resp.Body.Close()

    // 解析响应
    var result models.TransferResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("解析响应失败: %v", err)
    }

    return &result, nil
}
```

#### 4.4.2 转移记录管理

```go
// SaveTransferRecord 保存转移记录
func (s *GatewayService) SaveTransferRecord(record models.TransferRecord) error {
    // 构建API URL
    url := fmt.Sprintf("%s/records", s.gatewayURL)

    // 序列化请求体
    requestBody, err := json.Marshal(record)
    if err != nil {
        return fmt.Errorf("序列化记录失败: %v", err)
    }

    // 创建HTTP客户端
    client := &http.Client{
        Timeout: s.timeout,
    }

    // 发送POST请求
    resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return fmt.Errorf("保存记录失败: %v", err)
    }
    defer resp.Body.Close()

    // 检查响应状态
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("保存记录失败，状态码: %d", resp.StatusCode)
    }

    return nil
}
```

### 4.5 一致性验证

#### 4.5.1 哈希验证

```go
// VerifyDataConsistency 验证跨链数据一致性
func (s *GatewayService) VerifyDataConsistency(sourceData, targetData models.MedicalData) (bool, error) {
    // 验证基本字段
    if sourceData.DataHash != targetData.DataHash {
        return false, fmt.Errorf("数据哈希不匹配")
    }

    // 验证元数据一致性
    sourceMetadata := make(map[string]interface{})
    targetMetadata := make(map[string]interface{})

    if err := json.Unmarshal([]byte(sourceData.Metadata), &sourceMetadata); err != nil {
        return false, fmt.Errorf("解析源元数据失败: %v", err)
    }

    if err := json.Unmarshal([]byte(targetData.Metadata), &targetMetadata); err != nil {
        return false, fmt.Errorf("解析目标元数据失败: %v", err)
    }

    // 忽略跨链转换添加的字段
    delete(targetMetadata, "crossChainSource")
    delete(targetMetadata, "crossChainTimestamp")

    // 比较关键字段
    for key, sourceValue := range sourceMetadata {
        // 跳过链特定字段
        if key == "patientAddress" || key == "patientId" {
            continue
        }
        
        // 检查目标元数据中是否存在该字段
        targetValue, exists := targetMetadata[key]
        if !exists {
            return false, fmt.Errorf("目标元数据缺少字段: %s", key)
        }
        
        // 比较字段值（简单比较，可能需要更复杂的比较逻辑）
        if fmt.Sprintf("%v", sourceValue) != fmt.Sprintf("%v", targetValue) {
            return false, fmt.Errorf("字段 %s 的值不匹配: %v != %v", key, sourceValue, targetValue)
        }
    }

    return true, nil
}
```

#### 4.5.2 数据完整性验证

```go
// VerifyDataIntegrity 验证数据完整性
func (s *GatewayService) VerifyDataIntegrity(data models.MedicalData) (bool, error) {
    // 验证数据哈希
    // 这里应该实现实际的哈希验证逻辑，例如从IPFS获取数据并计算哈希
    // 为了示例，我们假设已经实现了这个功能
    calculatedHash, err := calculateDataHash(data.DataHash)
    if err != nil {
        return false, fmt.Errorf("计算哈希失败: %v", err)
    }
    
    if calculatedHash != data.DataHash {
        return false, fmt.Errorf("数据哈希不匹配: %s != %s", calculatedHash, data.DataHash)
    }
    
    return true, nil
}
```

### 4.6 错误处理

#### 4.6.1 错误类型定义

```go
// 定义错误类型常量
const (
    // 客户端错误
    ErrInvalidRequest   = "INVALID_REQUEST"   // 无效请求
    ErrUnauthorized     = "UNAUTHORIZED"      // 未授权
    ErrNotFound         = "NOT_FOUND"         // 资源不存在
    
    // 服务器错误
    ErrInternalServer   = "INTERNAL_SERVER"   // 内部服务器错误
    ErrDatabaseError    = "DATABASE_ERROR"    // 数据库错误
    
    // 区块链错误
    ErrEthereumError    = "ETHEREUM_ERROR"    // 以太坊错误
    ErrFabricError      = "FABRIC_ERROR"      // Fabric错误
    ErrCrossChainError  = "CROSS_CHAIN_ERROR" // 跨链错误
)

// GatewayError 网关错误结构
type GatewayError struct {
    Type    string `json:"type"`    // 错误类型
    Message string `json:"message"` // 错误消息
    Details string `json:"details"` // 错误详情
}
```

#### 4.6.2 错误处理中间件

```go
// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            // 默认为内部服务器错误
            statusCode := http.StatusInternalServerError
            errorType := ErrInternalServer
            errorMessage := "内部服务器错误"
            errorDetails := err.Error()
            
            // 根据错误类型设置状态码和错误信息
            switch e := err.(type) {
            case *GatewayError:
                // 自定义错误类型
                switch e.Type {
                case ErrInvalidRequest:
                    statusCode = http.StatusBadRequest
                    errorMessage = "无效请求"
                case ErrUnauthorized:
                    statusCode = http.StatusUnauthorized
                    errorMessage = "未授权"
                case ErrNotFound:
                    statusCode = http.StatusNotFound
                    errorMessage = "资源不存在"
                }
                errorDetails = e.Details
            }
            
            // 返回JSON格式的错误响应
            c.JSON(statusCode, gin.H{
                "error": gin.H{
                    "type":    errorType,
                    "message": errorMessage,
                    "details": errorDetails,
                },
            })
            
            // 终止后续处理
            c.Abort()
        }
    }
}
```

## 5. 测试策略

### 5.1 智能合约测试

#### 5.1.1 以太坊合约测试

以太坊智能合约测试使用Truffle框架和Ganache本地区块链进行：

```javascript
// 测试数据上传功能
it("should upload medical data correctly", async () => {
    const medicalData = await MedicalData.deployed();
    
    // 测试数据
    const dataHash = "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o";
    const dataType = "影像数据";
    const metadata = "{\"patientId\":\"P12345\",\"hospital\":\"协和医院\",\"department\":\"放射科\"}";
    const keywords = "肺部,CT,影像";
    
    // 上传数据
    const result = await medicalData.uploadData(
        dataHash,
        dataType,
        metadata,
        keywords,
        { from: accounts[0] }
    );
    
    // 验证事件
    assert.equal(result.logs.length, 1, "应该触发一个事件");
    assert.equal(result.logs[0].event, "DataUploaded", "应该是DataUploaded事件");
    
    // 获取数据并验证
    const dataId = result.logs[0].args.id.toNumber();
    const data = await medicalData.getData(dataId);
    
    assert.equal(data[0].toNumber(), dataId, "ID应该匹配");
    assert.equal(data[1], accounts[0], "所有者应该匹配");
    assert.equal(data[2], dataHash, "数据哈希应该匹配");
    assert.equal(data[3], dataType, "数据类型应该匹配");
    assert.equal(data[4], metadata, "元数据应该匹配");
    assert.equal(data[6], keywords, "关键词应该匹配");
});
```

#### 5.1.2 Fabric链码测试

Fabric链码测试使用Fabric链码测试框架：

```go
func TestUploadData(t *testing.T) {
    // 创建模拟链码环境
    ctx, chaincodeStub := createMockStub(t)
    
    // 创建链码实例
    chaincode := new(MedicalData)
    
    // 测试数据
    id := "test-001"
    owner := "user1"
    dataHash := "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"
    dataType := "影像数据"
    metadata := "{\"patientId\":\"P12345\",\"hospital\":\"协和医院\",\"department\":\"放射科\"}"
    keywords := "肺部,CT,影像"
    
    // 模拟调用上传数据函数
    err := chaincode.UploadData(ctx, id, owner, dataHash, dataType, metadata, keywords)
    
    // 验证结果
    assert.NoError(t, err, "上传数据应该成功")
    
    // 验证数据是否正确存储
    recordJSON, err := chaincodeStub.GetState(id)
    assert.NoError(t, err, "获取数据应该成功")
    assert.NotNil(t, recordJSON, "数据应该存在")
    
    // 解析数据并验证
    var record MedicalRecord
    err = json.Unmarshal(recordJSON, &record)
    assert.NoError(t, err, "解析数据应该成功")
    
    assert.Equal(t, id, record.ID, "ID应该匹配")
    assert.Equal(t, owner, record.Owner, "所有者应该匹配")
    assert.Equal(t, dataHash, record.DataHash, "数据哈希应该匹配")
    assert.Equal(t, dataType, record.DataType, "数据类型应该匹配")
    assert.Equal(t, metadata, record.Metadata, "元数据应该匹配")
    assert.Equal(t, keywords, record.Keywords, "关键词应该匹配")
}
```

### 5.2 跨链网关测试

#### 5.2.1 单元测试

```go
func TestEthereumToFabric(t *testing.T) {
    // 创建转换器实例
    converter := NewChainConverter()
    
    // 创建测试数据
    ethData := models.MedicalData{
        ID:        "eth-001",
        Owner:     "0x1234567890abcdef",
        DataHash:  "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o",
        DataType:  "影像数据",
        Metadata:  `{"patientAddress":"0xabcdef1234567890","hospital":"协和医院","department":"放射科"}`,
        Timestamp: time.Now(),
        Keywords:  "肺部,CT,影像",
        Chain:     "ethereum",
    }
    
    // 执行转换
    fabricData, err := converter.EthereumToFabric(ethData)
    
    // 验证结果
    assert.NoError(t, err, "转换应该成功")
    assert.Equal(t, "fabric", fabricData.Chain, "链标识应该是fabric")
    assert.Equal(t, "fab-001", fabricData.ID, "ID应该正确转换")
    assert.Equal(t, ethData.DataHash, fabricData.DataHash, "数据哈希应该保持不变")
    
    // 验证元数据转换
    var metadata map[string]interface{}
    err = json.Unmarshal([]byte(fabricData.Metadata), &metadata)
    assert.NoError(t, err, "解析元数据应该成功")
    
    // 验证patientId字段（从patientAddress转换而来）
    assert.Contains(t, metadata, "patientId", "应该包含patientId字段")
    assert.Equal(t, "abcdef1234567890", metadata["patientId"], "patientId应该正确转换")
    
    // 验证跨链信息
    assert.Contains(t, metadata, "crossChainSource", "应该包含crossChainSource字段")
    assert.Equal(t, "ethereum", metadata["crossChainSource"], "crossChainSource应该是ethereum")
}
```

#### 5.2.2 集成测试

```go
func TestQueryAllChains(t *testing.T) {
    // 跳过实际网络测试，除非明确指定
    if testing.Short() {
        t.Skip("跳过需要网络连接的测试")
    }
    
    // 创建网关服务实例
    service := NewGatewayService()
    
    // 创建查询参数
    query := models.MedicalDataQuery{
        Keyword:   "肺部",
        DataType:  "影像数据",
        Chain:     "all",
        Page:      1,
        PageSize:  10,
        SortBy:    "timestamp",
    }
    
    // 执行查询
    result, err := service.QueryAllChains(query)
    
    // 验证结果
    assert.NoError(t, err, "查询应该成功")
    assert.NotNil(t, result, "结果不应该为空")
    
    // 验证是否包含两条链的数据
    hasEthereum := false
    hasFabric := false
    
    for _, data := range result.Data {
        if data.Chain == "ethereum" {
            hasEthereum = true
        } else if data.Chain == "fabric" {
            hasFabric = true
        }
    }
    
    assert.True(t, hasEthereum || hasFabric, "应该至少包含一条链的数据")
}
```

### 5.3 性能测试

#### 5.3.1 负载测试

使用Apache JMeter或Locust进行负载测试，测试网关在高并发情况下的性能：

- 测试场景：模拟100个并发用户，每秒发送10个请求
- 测试指标：响应时间、吞吐量、错误率
- 测试持续时间：10分钟

#### 5.3.2 区块链性能测试

- 以太坊交易吞吐量测试
- Fabric交易吞吐量测试
- 跨链数据转移性能测试

## 6. 部署指南

### 6.1 以太坊智能合约部署

#### 6.1.1 开发环境部署

```bash
# 安装Truffle和Ganache
npm install -g truffle ganache-cli

# 启动本地以太坊节点
ganache-cli

# 编译合约
cd ethereum
truffle compile

# 部署合约
truffle migrate --network development
```

#### 6.1.2 测试网部署

```bash
# 配置.env文件
INFURA_API_KEY=your_infura_api_key
MNEMONIC=your_wallet_mnemonic

# 部署到Rinkeby测试网
truffle migrate --network rinkeby
```

#### 6.1.3 主网部署

```bash
# 部署到以太坊主网（谨慎操作）
truffle migrate --network mainnet
```

### 6.2 Fabric链码部署

#### 6.2.1 开发环境部署

```bash
# 启动Fabric测试网络
cd fabric-samples/test-network
./network.sh up createChannel -c medchannel

# 打包链码
cd ../../fabric
GO111MODULE=on go mod vendor
cd ..
tar cfz medicaldata.tar.gz fabric

# 安装链码
./network.sh deployCC -c medchannel -ccn medicaldata -ccp ../medicaldata.tar.gz -ccl go
```

#### 6.2.2 生产环境部署

1. 准备链码包
```bash
GO111MODULE=on go mod vendor
tar cfz medicaldata.tar.gz fabric
```

2. 在Fabric网络上安装链码
```bash
peer lifecycle chaincode package medicaldata.tar.gz --path ./fabric --lang golang --label medicaldata_1.0
peer lifecycle chaincode install medicaldata.tar.gz
```

3. 批准链码定义
```bash
peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --channelID medchannel --name medicaldata --version 1.0 --package-id $PACKAGE_ID --sequence 1
```

4. 提交链码定义
```bash
peer lifecycle chaincode commit -o orderer.example.com:7050 --channelID medchannel --name medicaldata --version 1.0 --sequence 1
```

### 6.3 跨链网关部署

#### 6.3.1 Docker部署

1. 创建Dockerfile
```dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./crosschain-gateway

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway /app/
COPY --from=builder /app/.env /app/

EXPOSE 8080
CMD ["/app/gateway"]
```

2. 构建和运行Docker镜像
```bash
docker build -t medcross/gateway:latest .
docker run -p 8080:8080 medcross/gateway:latest
```

#### 6.3.2 Kubernetes部署

1. 创建Kubernetes部署配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: medcross-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: medcross-gateway
  template:
    metadata:
      labels:
        app: medcross-gateway
    spec:
      containers:
      - name: gateway
        image: medcross/gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: ETHEREUM_URL
          valueFrom:
            configMapKeyRef:
              name: blockchain-config
              key: ethereum_url
        - name: FABRIC_CONFIG
          valueFrom:
            configMapKeyRef:
              name: blockchain-config
              key: fabric_config
---
apiVersion: v1
kind: Service
metadata:
  name: medcross-gateway-service
spec:
  selector:
    app: medcross-gateway
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

2. 部署到Kubernetes集群
```bash
kubectl apply -f kubernetes/gateway.yaml
```

## 7. 维护与监控

### 7.1 日志管理

- 使用ELK（Elasticsearch, Logstash, Kibana）堆栈收集和分析日志
- 实现结构化日志记录，包含时间戳、请求ID、错误级别等信息
- 设置日志轮转和归档策略

### 7.2 监控指标

- 系统健康状态：CPU、内存、磁盘使用率
- API性能指标：请求量、响应时间、错误率
- 区块链指标：交易确认时间、Gas使用量、链码调用次数
- 跨链指标：转移成功率、转移时间、数据一致性

### 7.3 告警机制

- 设置关键指标阈值告警
- 实现多渠道告警通知（邮件、短信、企业微信等）
- 建立告警升级机制和处理流程

## 8. 安全最佳实践

### 8.1 智能合约安全

- 使用OpenZeppelin等安全库
- 进行形式化验证和安全审计
- 实现紧急暂停机制
- 遵循最小权限原则

### 8.2 跨链网关安全

- 实现API认证和授权
- 使用HTTPS加密传输
- 实现请求限流和防DDoS措施
- 定期进行安全漏洞扫描

### 8.3 密钥管理

- 使用硬件安全模块（HSM）存储私钥
- 实现密钥轮换机制
- 采用多重签名机制
- 建立密钥备份和恢复流程

## 9. 常见问题与解决方案

### 9.1 智能合约问题

| 问题 | 解决方案 |
|------|----------|
| Gas费用过高 | 优化合约代码，减少存储操作，使用批量处理 |
| 交易确认慢 | 调整Gas价格，使用交易池监控，实现异步确认机制 |
| 合约升级 | 使用代理模式设计合约，实现可升级性 |

### 9.2 跨链问题

| 问题 | 解决方案 |
|------|----------|
| 数据不一致 | 实现数据校验和修复机制，使用事务补偿模式 |
| 跨链延迟高 | 优化网关性能，实现异步处理，使用缓存 |
| 链间通信失败 | 实现重试机制，使用消息队列，建立失败恢复流程 |

### 9.3 性能问题

| 问题 | 解决方案 |
|------|----------|
| 查询响应慢 | 优化索引，实现缓存，使用读写分离 |
| 高并发处理 | 实现负载均衡，水平扩展，使用连接池 |
| 资源消耗高 | 优化算法，实现资源限制，监控资源使用 |

## 10. 未来发展路线

### 10.1 技术升级

- 支持更多区块链平台（如Polkadot、Cosmos等）
- 实现跨链智能合约互操作性
- 集成零知识证明等隐私保护技术
- 支持Layer 2扩展解决方案

### 10.2 功能扩展

- 实现跨链数据分析和挖掘
- 支持医疗数据的AI处理和应用
- 建立医疗数据交易市场
- 实现基于区块链的医疗数据授权和审计