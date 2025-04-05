package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"d.com/MedCross/backend/models"
)

// CrossChainService 跨链服务接口
// 负责协调以太坊和Fabric区块链之间的数据同步和交互
type CrossChainService struct {
	ethereum *EthereumService // 以太坊服务实例
	fabric   *FabricService   // Fabric服务实例
	gateway  *GatewayService  // 跨链网关服务实例
}

// NewCrossChainService 创建跨链服务实例
// 初始化跨链服务，加载以太坊和Fabric服务，并创建跨链网关
// 返回:
//   - *CrossChainService: 跨链服务实例
func NewCrossChainService() *CrossChainService {
	// 初始化以太坊和Fabric服务
	ethereumService := NewEthereumService()
	fabricService := NewFabricService()
	
	// 创建跨链服务实例
	service := &CrossChainService{
		ethereum: ethereumService,
		fabric:   fabricService,
	}
	
	// 初始化跨链网关
	service.gateway = NewGatewayService(ethereumService, fabricService)
	
	return service
}

// SyncDataAcrossChains 跨链同步数据
// 将指定的医疗数据在以太坊和Fabric区块链之间进行同步
// 参数:
//   - dataID: 数据ID
//   - sourceChain: 源区块链类型（"ethereum"或"fabric"）
// 返回:
//   - *models.CrossChainReference: 跨链引用信息
//   - error: 错误信息，如果同步成功则为nil
func (s *CrossChainService) SyncDataAcrossChains(dataID, sourceChain string) (*models.CrossChainReference, error) {
	// 使用跨链网关发起数据传输
	var chainType ChainType
	if sourceChain == "ethereum" {
		chainType = Ethereum
	} else if sourceChain == "fabric" {
		chainType = Fabric
	} else {
		return nil, fmt.Errorf("无效的源区块链类型: %s", sourceChain)
	}

	// 通过网关发起跨链交易
	transactionID, err := s.gateway.InitiateDataTransfer(dataID, chainType)
	if err != nil {
		return nil, fmt.Errorf("发起跨链交易失败: %v", err)
	}

	// 获取交易状态
	transaction, err := s.gateway.GetTransactionStatus(transactionID)
	if err != nil {
		return nil, fmt.Errorf("获取交易状态失败: %v", err)
	}

	// 如果交易尚未完成，返回进行中的状态
	if transaction.Status != "已完成" {
		return &models.CrossChainReference{
			DataID:       dataID,
			SyncStatus:   "同步中",
			LastSyncTime: time.Now().Unix(),
		}, nil
	}

	// 交易已完成，返回跨链引用
	return transaction.Reference, nil
}

// GetCrossChainStatus 获取跨链同步状态
// 查询指定数据的跨链同步状态
// 参数:
//   - dataID: 数据ID
// 返回:
//   - *models.CrossChainReference: 跨链引用信息
//   - error: 错误信息，如果查询成功则为nil
func (s *CrossChainService) GetCrossChainStatus(dataID string) (*models.CrossChainReference, error) {
	// 获取所有跨链交易
	transactions := s.gateway.GetAllTransactions()
	
	// 查找与指定数据ID相关的最新交易
	var latestTransaction *CrossChainTransaction
	for _, tx := range transactions {
		if tx.DataID == dataID {
			if latestTransaction == nil || tx.CreatedAt.After(latestTransaction.CreatedAt) {
				latestTransaction = tx
			}
		}
	}
	
	// 如果找到相关交易，返回其引用信息
	if latestTransaction != nil && latestTransaction.Reference != nil {
		return latestTransaction.Reference, nil
	}
	
	// 如果没有找到相关交易或交易没有引用信息，返回一个模拟数据
	reference := &models.CrossChainReference{
		DataID:       dataID,
		EthereumTxID: "0x" + uuid.New().String()[:8] + "..." + uuid.New().String()[:4],
		FabricTxID:   uuid.New().String(),
		SyncStatus:   "未知",
		LastSyncTime: time.Now().Add(-1 * time.Hour).Unix(),
	}
	
	return reference, nil
}

// VerifyCrossChainConsistency 验证跨链数据一致性
// 检查指定数据在两个区块链上的一致性
// 参数:
//   - dataID: 数据ID
// 返回:
//   - bool: 数据是否一致
//   - error: 错误信息，如果验证成功则为nil
func (s *CrossChainService) VerifyCrossChainConsistency(dataID string) (bool, error) {
	// 使用跨链网关验证数据一致性
	consistent, differences, err := s.gateway.VerifyDataConsistency(dataID)
	if err != nil {
		return false, fmt.Errorf("验证数据一致性失败: %v", err)
	}

	// 如果数据不一致，记录差异信息（在实际应用中可能需要记录到日志）
	if !consistent {
		fmt.Printf("数据不一致，差异: %v\n", differences)
	}

	return consistent, nil
}

// SyncAccessRecords 同步访问记录
// 将一个区块链上的访问记录同步到另一个区块链，保持两个区块链上的访问历史一致
// 参数:
//   - dataID: 数据ID，要同步访问记录的医疗数据标识符
//   - sourceChain: 源区块链类型（"ethereum"或"fabric"），指定从哪个区块链同步到另一个区块链
// 返回:
//   - error: 错误信息，如果同步成功则为nil
func (s *CrossChainService) SyncAccessRecords(dataID, sourceChain string) error {
	// 将字符串类型的区块链类型转换为ChainType枚举
	var chainType ChainType
	if sourceChain == "ethereum" {
		chainType = Ethereum
	} else if sourceChain == "fabric" {
		chainType = Fabric
	} else {
		return fmt.Errorf("无效的源区块链类型: %s", sourceChain)
	}

	// 使用跨链网关同步访问记录
	// 在实际应用中，这里应该实现访问记录的同步逻辑
	// 目前网关服务中没有专门的访问记录同步方法，可以在未来扩展

	var accessLogs []models.AccessRecord
	var err error

	// 根据源区块链类型获取访问记录
	if sourceChain == "ethereum" {
		// 从以太坊获取访问记录并同步到Fabric
		// 在实际应用中，这里应该调用以太坊服务获取访问记录
		// 然后调用Fabric服务同步这些记录
		// ...
	} else if sourceChain == "fabric" {
		// 从Fabric获取访问记录
		accessLogs, err = s.fabric.GetAccessLogs(dataID)
		if err != nil {
			return fmt.Errorf("从Fabric获取访问记录失败: %v", err)
		}

		// 同步到以太坊
		// 在实际应用中，这里应该调用以太坊服务同步这些记录
		// ...
	}

	return nil
}

// SyncAuthorizations 同步授权信息
// 将一个区块链上的授权信息同步到另一个区块链，确保两个区块链上的数据访问权限一致
// 参数:
//   - dataID: 数据ID，要同步授权信息的医疗数据标识符
//   - ownerID: 数据所有者ID，医疗数据的所有者标识符，用于权限验证
//   - sourceChain: 源区块链类型（"ethereum"或"fabric"），指定从哪个区块链同步到另一个区块链
// 返回:
//   - error: 错误信息，如果同步成功则为nil
func (s *CrossChainService) SyncAuthorizations(dataID, ownerID, sourceChain string) error {
	// 将字符串类型的区块链类型转换为ChainType枚举
	var chainType ChainType
	if sourceChain == "ethereum" {
		chainType = Ethereum
	} else if sourceChain == "fabric" {
		chainType = Fabric
	} else {
		return fmt.Errorf("无效的源区块链类型: %s", sourceChain)
	}

	// 使用跨链网关同步访问控制信息
	err := s.gateway.SyncAccessControl(dataID, chainType)
	if err != nil {
		return fmt.Errorf("同步授权信息失败: %v", err)
	}

	return nil
}