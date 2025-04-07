package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MedicalData 智能合约实现
type MedicalData struct {
	contractapi.Contract
}

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

// GetDataByOwner 获取特定所有者的所有数据
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
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		// 从复合键中提取ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(response.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}

		if len(compositeKeyParts) > 1 {
			id := compositeKeyParts[1]
			// 获取数据记录
			record, err := s.GetData(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get data: %v", err)
			}
			records = append(records, record)
		}
	}

	return records, nil
}

// GetDataByType 获取特定类型的所有数据
func (s *MedicalData) GetDataByType(ctx contractapi.TransactionContextInterface, dataType string) ([]*MedicalRecord, error) {
	// 创建复合键迭代器
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("type~id", []string{dataType})
	if err != nil {
		return nil, fmt.Errorf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	var records []*MedicalRecord

	// 遍历结果
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		// 从复合键中提取ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(response.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}

		if len(compositeKeyParts) > 1 {
			id := compositeKeyParts[1]
			// 获取数据记录
			record, err := s.GetData(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get data: %v", err)
			}
			records = append(records, record)
		}
	}

	return records, nil
}

// QueryDataByKeywords 根据关键词查询数据
func (s *MedicalData) QueryDataByKeywords(ctx contractapi.TransactionContextInterface, keyword string) ([]*MedicalRecord, error) {
	// 构建富查询
	queryString := fmt.Sprintf(`{"selector":{"keywords":{"$regex":".*%s.*"}}}", keyword)

	// 执行查询
	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var records []*MedicalRecord

	// 遍历结果
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		var record MedicalRecord
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %v", err)
		}

		records = append(records, &record)
	}

	return records, nil
}

// DataExists 检查数据是否存在
func (s *MedicalData) DataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	recordJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return recordJSON != nil, nil
}

// GetAllData 获取所有医疗数据
func (s *MedicalData) GetAllData(ctx contractapi.TransactionContextInterface) ([]*MedicalRecord, error) {
	// 获取所有数据的迭代器
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get all data: %v", err)
	}
	defer iterator.Close()

	var records []*MedicalRecord

	// 遍历所有数据
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		// 跳过复合键
		if response.Key[0:1] != "~" {
			var record MedicalRecord
			err = json.Unmarshal(response.Value, &record)
			if err != nil {
				// 跳过非医疗数据记录
				continue
			}

			records = append(records, &record)
		}
	}

	return records, nil
}

func main() {
	contract := new(MedicalData)

	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		fmt.Printf("Error creating MedicalData chaincode: %v\n", err)
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting MedicalData chaincode: %v\n", err)
	}
}