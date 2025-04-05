package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"d.com/MedCross/backend/models"
)

// ChainType 区块链类型枚举
type ChainType string

const (
	Ethereum ChainType = "ethereum"
	Fabric   ChainType = "fabric"
)

// DataTransformer 数据转换接口
// 负责在不同区块链数据格式之间进行转换
type DataTransformer interface {
	TransformToEthereum(data interface{}) (interface{}, error)
	TransformToFabric(data interface{}) (interface{}, error)
}

// EventListener 事件监听接口
// 负责监听区块链上的事件并触发相应的处理逻辑
type EventListener interface {
	StartListening(chainType ChainType) error
	StopListening(chainType ChainType) error
	RegisterHandler(eventType string, handler func(event interface{}) error)
}

// GatewayService 跨链网关服务
// 负责协调以太坊和Fabric区块链之间的数据转换和传输
type GatewayService struct {
	ethereum       *EthereumService  // 以太坊服务实例
	fabric         *FabricService    // Fabric服务实例
	transformer    DataTransformer   // 数据转换器
	eventListener  EventListener     // 事件监听器
	transactionMap sync.Map          // 交易映射表，用于跟踪跨链交易状态
	pendingQueue   chan *CrossChainTransaction // 待处理交易队列
}

// CrossChainTransaction 跨链交易结构
// 用于在网关内部跟踪跨链交易的状态和信息
type CrossChainTransaction struct {
	ID            string                     // 交易唯一标识符
	DataID        string                     // 数据ID
	SourceChain   ChainType                  // 源区块链类型
	TargetChain   ChainType                  // 目标区块链类型
	Status        string                     // 交易状态（等待中、处理中、已完成、失败）
	CreatedAt     time.Time                  // 创建时间
	CompletedAt   time.Time                  // 完成时间
	Reference     *models.CrossChainReference // 跨链引用信息
	ErrorMessage  string                     // 错误信息
}

// NewGatewayService 创建跨链网关服务实例
// 初始化网关服务，加载以太坊和Fabric服务，设置数据转换器和事件监听器
// 返回:
//   - *GatewayService: 网关服务实例
func NewGatewayService(ethereum *EthereumService, fabric *FabricService) *GatewayService {
	service := &GatewayService{
		ethereum:     ethereum,
		fabric:       fabric,
		transformer:  NewDefaultDataTransformer(),
		pendingQueue: make(chan *CrossChainTransaction, 100),
	}

	// 启动交易处理协程
	go service.processTransactions()

	return service
}

// DefaultDataTransformer 默认数据转换器实现
type DefaultDataTransformer struct{}

// NewDefaultDataTransformer 创建默认数据转换器
func NewDefaultDataTransformer() *DefaultDataTransformer {
	return &DefaultDataTransformer{}
}

// TransformToEthereum 将数据转换为以太坊格式
func (t *DefaultDataTransformer) TransformToEthereum(data interface{}) (interface{}, error) {
	// 在实际应用中，这里应该实现具体的转换逻辑
	// 例如将Fabric的数据结构转换为以太坊智能合约需要的格式
	return data, nil
}

// TransformToFabric 将数据转换为Fabric格式
func (t *DefaultDataTransformer) TransformToFabric(data interface{}) (interface{}, error) {
	// 在实际应用中，这里应该实现具体的转换逻辑
	// 例如将以太坊的数据结构转换为Fabric链码需要的格式
	return data, nil
}

// InitiateDataTransfer 发起数据跨链传输
// 创建跨链交易并将其加入处理队列
// 参数:
//   - dataID: 数据ID
//   - sourceChain: 源区块链类型
// 返回:
//   - string: 跨链交易ID
//   - error: 错误信息
func (g *GatewayService) InitiateDataTransfer(dataID string, sourceChain ChainType) (string, error) {
	// 创建跨链交易
	transaction := &CrossChainTransaction{
		ID:          fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		DataID:      dataID,
		SourceChain: sourceChain,
		Status:      "等待中",
		CreatedAt:   time.Now(),
	}

	// 设置目标链
	if sourceChain == Ethereum {
		transaction.TargetChain = Fabric
	} else {
		transaction.TargetChain = Ethereum
	}

	// 存储交易信息
	g.transactionMap.Store(transaction.ID, transaction)

	// 将交易加入处理队列
	g.pendingQueue <- transaction

	return transaction.ID, nil
}

// processTransactions 处理跨链交易队列
// 持续监听交易队列并处理待处理的交易
func (g *GatewayService) processTransactions() {
	for tx := range g.pendingQueue {
		// 更新交易状态
		tx.Status = "处理中"
		g.transactionMap.Store(tx.ID, tx)

		// 执行跨链数据传输
		go func(transaction *CrossChainTransaction) {
			var reference *models.CrossChainReference
			var err error

			// 根据源链类型执行不同的传输逻辑
			if transaction.SourceChain == Ethereum {
				// 从以太坊传输到Fabric
				reference, err = g.transferFromEthereumToFabric(transaction.DataID)
			} else {
				// 从Fabric传输到以太坊
				reference, err = g.transferFromFabricToEthereum(transaction.DataID)
			}

			// 更新交易状态
			transaction.CompletedAt = time.Now()
			transaction.Reference = reference

			if err != nil {
				transaction.Status = "失败"
				transaction.ErrorMessage = err.Error()
				log.Printf("跨链交易失败: %s, 错误: %v", transaction.ID, err)
			} else {
				transaction.Status = "已完成"
				log.Printf("跨链交易完成: %s", transaction.ID)
			}

			// 更新存储的交易信息
			g.transactionMap.Store(transaction.ID, transaction)
		}(tx)
	}
}

// transferFromEthereumToFabric 从以太坊传输数据到Fabric
// 参数:
//   - dataID: 数据ID
// 返回:
//   - *models.CrossChainReference: 跨链引用信息
//   - error: 错误信息
func (g *GatewayService) transferFromEthereumToFabric(dataID string) (*models.CrossChainReference, error) {
	// 从以太坊获取数据
	data, err := g.ethereum.GetDataById(dataID)
	if err != nil {
		return nil, fmt.Errorf("从以太坊获取数据失败: %v", err)
	}

	// 转换数据格式（如果需要）
	transformedData, err := g.transformer.TransformToFabric(data)
	if err != nil {
		return nil, fmt.Errorf("数据格式转换失败: %v", err)
	}

	// 上传到Fabric
	fabricTxID, err := g.fabric.UploadData(*data)
	if err != nil {
		return nil, fmt.Errorf("上传到Fabric失败: %v", err)
	}

	// 创建跨链引用
	reference := &models.CrossChainReference{
		DataID:       dataID,
		EthereumTxID: data.EthereumTxID,
		FabricTxID:   fabricTxID,
		SyncStatus:   "已同步",
		LastSyncTime: time.Now().Unix(),
	}

	return reference, nil
}

// transferFromFabricToEthereum 从Fabric传输数据到以太坊
// 参数:
//   - dataID: 数据ID
// 返回:
//   - *models.CrossChainReference: 跨链引用信息
//   - error: 错误信息
func (g *GatewayService) transferFromFabricToEthereum(dataID string) (*models.CrossChainReference, error) {
	// 从Fabric获取数据
	data, err := g.fabric.GetDataByID(dataID)
	if err != nil {
		return nil, fmt.Errorf("从Fabric获取数据失败: %v", err)
	}

	// 转换数据格式（如果需要）
	transformedData, err := g.transformer.TransformToEthereum(data)
	if err != nil {
		return nil, fmt.Errorf("数据格式转换失败: %v", err)
	}

	// 上传到以太坊
	ethereumTxID, err := g.ethereum.UploadData(data, data.FabricTxID)
	if err != nil {
		return nil, fmt.Errorf("上传到以太坊失败: %v", err)
	}

	// 创建跨链引用
	reference := &models.CrossChainReference{
		DataID:       dataID,
		EthereumTxID: ethereumTxID,
		FabricTxID:   data.FabricTxID,
		SyncStatus:   "已同步",
		LastSyncTime: time.Now().Unix(),
	}

	return reference, nil
}

// GetTransactionStatus 获取跨链交易状态
// 参数:
//   - transactionID: 交易ID
// 返回:
//   - *CrossChainTransaction: 交易信息
//   - error: 错误信息
func (g *GatewayService) GetTransactionStatus(transactionID string) (*CrossChainTransaction, error) {
	// 从交易映射表中获取交易信息
	value, exists := g.transactionMap.Load(transactionID)
	if !exists {
		return nil, fmt.Errorf("交易不存在: %s", transactionID)
	}

	transaction, ok := value.(*CrossChainTransaction)
	if !ok {
		return nil, fmt.Errorf("交易数据类型错误")
	}

	return transaction, nil
}

// GetAllTransactions 获取所有跨链交易
// 返回:
//   - []*CrossChainTransaction: 交易列表
func (g *GatewayService) GetAllTransactions() []*CrossChainTransaction {
	var transactions []*CrossChainTransaction

	g.transactionMap.Range(func(key, value interface{}) bool {
		if tx, ok := value.(*CrossChainTransaction); ok {
			transactions = append(transactions, tx)
		}
		return true
	})

	return transactions
}

// SyncAccessControl 同步访问控制信息
// 将一个区块链上的授权信息同步到另一个区块链
// 参数:
//   - dataID: 数据ID
//   - sourceChain: 源区块链类型
// 返回:
//   - error: 错误信息
func (g *GatewayService) SyncAccessControl(dataID string, sourceChain ChainType) error {
	var authorizations []models.Authorization
	var err error

	// 根据源链类型获取授权信息
	if sourceChain == Ethereum {
		// 从以太坊获取授权信息
		authorizations, err = g.ethereum.GetDataAuthorizations(dataID, "")
		if err != nil {
			return fmt.Errorf("从以太坊获取授权信息失败: %v", err)
		}

		// 同步到Fabric
		for _, auth := range authorizations {
			_, err := g.fabric.GrantAccess(
				auth.DataID,
				auth.OwnerID,
				auth.AuthorizedUser,
				auth.StartTime,
				auth.EndTime,
			)
			if err != nil {
				return fmt.Errorf("同步授权到Fabric失败: %v", err)
			}
		}
	} else {
		// 从Fabric获取授权信息
		authorizations, err = g.fabric.GetAuthorizations(dataID)
		if err != nil {
			return fmt.Errorf("从Fabric获取授权信息失败: %v", err)
		}

		// 同步到以太坊
		for _, auth := range authorizations {
			_, err := g.ethereum.GrantAccess(
				auth.DataID,
				auth.OwnerID,
				auth.AuthorizedUser,
				auth.StartTime,
				auth.EndTime,
			)
			if err != nil {
				return fmt.Errorf("同步授权到以太坊失败: %v", err)
			}
		}
	}

	return nil
}

// VerifyDataConsistency 验证跨链数据一致性
// 检查指定数据在两个区块链上的一致性
// 参数:
//   - dataID: 数据ID
// 返回:
//   - bool: 数据是否一致
//   - map[string]interface{}: 不一致的字段及其值
//   - error: 错误信息
func (g *GatewayService) VerifyDataConsistency(dataID string) (bool, map[string]interface{}, error) {
	// 获取以太坊上的数据
	ethData, err := g.ethereum.GetDataById(dataID)
	if err != nil {
		return false, nil, fmt.Errorf("从以太坊获取数据失败: %v", err)
	}

	// 获取Fabric上的数据
	fabricData, err := g.fabric.GetDataByID(dataID)
	if err != nil {
		return false, nil, fmt.Errorf("从Fabric获取数据失败: %v", err)
	}

	// 比较数据一致性
	differences := make(map[string]interface{})

	// 比较关键字段
	if ethData.DataHash != fabricData.DataHash {
		differences["DataHash"] = map[string]string{
			"ethereum": ethData.DataHash,
			"fabric":   fabricData.DataHash,
		}
	}

	if ethData.DataType != fabricData.DataType {
		differences["DataType"] = map[string]string{
			"ethereum": ethData.DataType,
			"fabric":   fabricData.DataType,
		}
	}

	if ethData.Description != fabricData.Description {
		differences["Description"] = map[string]string{
			"ethereum": ethData.Description,
			"fabric":   fabricData.Description,
		}
	}

	if ethData.Owner != fabricData.Owner {
		differences["Owner"] = map[string]string{
			"ethereum": ethData.Owner,
			"fabric":   fabricData.Owner,
		}
	}

	// 比较标签（可能顺序不同）
	ethTags := make(map[string]bool)
	fabricTags := make(map[string]bool)

	for _, tag := range ethData.Tags {
		ethTags[tag] = true
	}

	for _, tag := range fabricData.Tags {
		fabricTags[tag] = true
	}

	if !mapsEqual(ethTags, fabricTags) {
		differences["Tags"] = map[string][]string{
			"ethereum": ethData.Tags,
			"fabric":   fabricData.Tags,
		}
	}

	return len(differences) == 0, differences, nil
}

// mapsEqual 比较两个map是否相等
func mapsEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if _, exists := b[k]; !exists {
			return false
		}
	}

	return true
}

// SerializeTransaction 序列化交易信息为JSON
// 参数:
//   - transaction: 跨链交易
// 返回:
//   - []byte: JSON字节数组
//   - error: 错误信息
func (g *GatewayService) SerializeTransaction(transaction *CrossChainTransaction) ([]byte, error) {
	return json.Marshal(transaction)
}

// DeserializeTransaction 从JSON反序列化交易信息
// 参数:
//   - data: JSON字节数组
// 返回:
//   - *CrossChainTransaction: 跨链交易
//   - error: 错误信息
func (g *GatewayService) DeserializeTransaction(data []byte) (*CrossChainTransaction, error) {
	var transaction CrossChainTransaction
	err := json.Unmarshal(data, &transaction)
	return &transaction, err
}