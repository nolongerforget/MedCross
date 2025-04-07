package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"medcross/models"
)

// ChainConverter 区块链数据转换工具
// 用于处理以太坊和Fabric之间的数据格式转换
type ChainConverter struct{}

// NewChainConverter 创建新的链转换器
func NewChainConverter() *ChainConverter {
	return &ChainConverter{}
}

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
	updatedJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("序列化元数据失败: %v", err)
		return "", err
	}

	return string(updatedJSON), nil
}

// 将以太坊地址转换为Fabric ID
func convertEthAddressToFabricID(ethAddress string) string {
	// 移除0x前缀
	if strings.HasPrefix(ethAddress, "0x") {
		ethAddress = strings.TrimPrefix(ethAddress, "0x")
	}

	// 为简单起见，我们使用一个固定的前缀
	return fmt.Sprintf("fabric-user-%s", ethAddress[:8])
}

// 将Fabric ID转换为以太坊地址
func convertFabricIDToEthAddress(fabricID string) string {
	// 为简单起见，我们生成一个模拟的以太坊地址
	// 在实际应用中，应该使用更复杂的映射或转换逻辑
	if strings.HasPrefix(fabricID, "fabric-user-") {
		// 提取用户ID部分
		userPart := strings.TrimPrefix(fabricID, "fabric-user-")
		return fmt.Sprintf("0x%s000000000000000000000000", userPart)
	}

	// 生成一个默认地址
	return "0x" + fmt.Sprintf("%040x", time.Now().UnixNano())
}
