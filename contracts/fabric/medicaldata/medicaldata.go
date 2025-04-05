package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MedicalDataContract 定义医疗数据链码的结构
type MedicalDataContract struct {
	contractapi.Contract
}

// MedicalData 定义医疗数据的结构
type MedicalData struct {
	DataID        string   `json:"dataId"`        // 数据唯一标识符
	DataHash      string   `json:"dataHash"`      // 数据哈希值
	DataType      string   `json:"dataType"`      // 数据类型
	Description   string   `json:"description"`   // 数据描述
	Tags          []string `json:"tags"`          // 数据标签
	Owner         string   `json:"owner"`         // 数据所有者（用户ID）
	Timestamp     int64    `json:"timestamp"`     // 上传时间戳
	IsConfidential bool     `json:"isConfidential"` // 是否为机密数据
	EthereumTxID  string   `json:"ethereumTxId"`  // 以太坊链上对应的交易ID
}

// Authorization 定义授权信息的结构
type Authorization struct {
	DataID         string `json:"dataId"`         // 被授权的数据ID
	AuthorizedUser string `json:"authorizedUser"` // 被授权的用户ID
	StartTime      int64  `json:"startTime"`      // 授权开始时间
	EndTime        int64  `json:"endTime"`        // 授权结束时间（0表示永久授权）
	IsActive       bool   `json:"isActive"`       // 授权是否有效
}

// AccessRecord 定义访问记录的结构
type AccessRecord struct {
	DataID    string `json:"dataId"`    // 访问的数据ID
	Accessor  string `json:"accessor"`  // 访问者ID
	Timestamp int64  `json:"timestamp"` // 访问时间戳
	Operation string `json:"operation"` // 操作类型（查看、下载等）
}

// InitLedger 初始化账本
func (s *MedicalDataContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 初始化账本，可以为空或添加一些测试数据
	return nil
}

// UploadData 上传新的医疗数据
func (s *MedicalDataContract) UploadData(ctx contractapi.TransactionContextInterface, dataID, dataHash, dataType, description string, tags []string, owner string, isConfidential bool, ethereumTxID string) error {
	// 检查数据是否已存在
	exists, err := s.DataExists(ctx, dataID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("数据ID %s 已存在", dataID)
	}

	// 创建新的医疗数据记录
	medicalData := MedicalData{
		DataID:        dataID,
		DataHash:      dataHash,
		DataType:      dataType,
		Description:   description,
		Tags:          tags,
		Owner:         owner,
		Timestamp:     time.Now().Unix(),
		IsConfidential: isConfidential,
		EthereumTxID:  ethereumTxID,
	}

	// 将数据转换为JSON并存储到账本
	medicalDataJSON, err := json.Marshal(medicalData)
	if err != nil {
		return err
	}

	// 将数据写入账本
	err = ctx.GetStub().PutState(dataID, medicalDataJSON)
	if err != nil {
		return fmt.Errorf("存储数据失败: %v", err)
	}

	// 记录数据上传事件
	return ctx.GetStub().SetEvent("DataUploaded", medicalDataJSON)
}

// GrantAccess 授权其他用户访问数据
func (s *MedicalDataContract) GrantAccess(ctx contractapi.TransactionContextInterface, dataID, authorizedUser string, startTime, endTime int64) error {
	// 检查数据是否存在
	exists, err := s.DataExists(ctx, dataID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("数据ID %s 不存在", dataID)
	}

	// 获取数据信息
	medicalData, err := s.GetMedicalData(ctx, dataID)
	if err != nil {
		return err
	}

	// 检查调用者是否是数据所有者
	creator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取调用者身份失败: %v", err)
	}

	// 在实际应用中，需要将creator转换为可读的用户ID进行比较
	// 这里简化处理，假设creator就是用户ID
	if medicalData.Owner != creator {
		return fmt.Errorf("只有数据所有者才能授权访问")
	}

	// 创建授权记录
	authorization := Authorization{
		DataID:         dataID,
		AuthorizedUser: authorizedUser,
		StartTime:      startTime,
		EndTime:        endTime,
		IsActive:       true,
	}

	// 生成授权记录的复合键
	authKey, err := ctx.GetStub().CreateCompositeKey("authorization", []string{dataID, authorizedUser})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 将授权记录转换为JSON并存储
	authJSON, err := json.Marshal(authorization)
	if err != nil {
		return err
	}

	// 将授权记录写入账本
	err = ctx.GetStub().PutState(authKey, authJSON)
	if err != nil {
		return fmt.Errorf("存储授权记录失败: %v", err)
	}

	// 记录授权事件
	return ctx.GetStub().SetEvent("AuthorizationGranted", authJSON)
}

// RevokeAccess 撤销用户的数据访问权限
func (s *MedicalDataContract) RevokeAccess(ctx contractapi.TransactionContextInterface, dataID, authorizedUser string) error {
	// 检查数据是否存在
	exists, err := s.DataExists(ctx, dataID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("数据ID %s 不存在", dataID)
	}

	// 获取数据信息
	medicalData, err := s.GetMedicalData(ctx, dataID)
	if err != nil {
		return err
	}

	// 检查调用者是否是数据所有者
	creator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取调用者身份失败: %v", err)
	}

	// 在实际应用中，需要将creator转换为可读的用户ID进行比较
	if medicalData.Owner != creator {
		return fmt.Errorf("只有数据所有者才能撤销授权")
	}

	// 生成授权记录的复合键
	authKey, err := ctx.GetStub().CreateCompositeKey("authorization", []string{dataID, authorizedUser})
	if err != nil {
		return fmt.Errorf("创建复合键失败: %v", err)
	}

	// 获取授权记录
	authJSON, err := ctx.GetStub().GetState(authKey)
	if err != nil {
		return fmt.Errorf("获取授权记录失败: %v", err)
	}
	if authJSON == nil {
		return fmt.Errorf("授权记录不存在")
	}

	// 解析授权记录
	var authorization Authorization
	err = json.Unmarshal(authJSON, &authorization)
	if err != nil {
		return fmt.Errorf("解析授权记录失败: %v", err)
	}

	// 撤销授权
	authorization.IsActive = false

	// 将更新后的授权记录转换为JSON并存储
	updatedAuthJSON, err := json.Marshal(authorization)
	if err != nil {
		return err
	}

	// 将更新后的授权记录写入账本
	err = ctx.GetStub().PutState(authKey, updatedAuthJSON)
	if err != nil {
		return fmt.Errorf("更新授权记录失败: %v", err)
	}

	// 记录撤销授权事件
	return ctx.GetStub().SetEvent("AuthorizationRevoked", updatedAuthJSON)
}

// CheckAccess 检查用户是否有权访问数据
func (s *MedicalDataContract) CheckAccess(ctx contractapi.TransactionContextInterface, dataID, userID string) (bool, error) {
	// 检查数据是否存在
	exists, err := s.DataExists(ctx, dataID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("数据ID %s 不存在", dataID)
	}

	// 获取数据信息
	medicalData, err := s.GetMedicalData(ctx, dataID)
	if err != nil {
		return false, err
	}

	// 数据所有者始终有访问权限
	if medicalData.Owner == userID {
		return true, nil
	}

	// 生成授权记录的复合键
	authKey, err := ctx.GetStub().CreateCompositeKey("authorization", []string{dataID, userID})
	if err != nil {
		return false, fmt.Errorf("创建复合键失败: %v", err)
	}

	// 获取授权记录
	authJSON, err := ctx.GetStub().GetState(authKey)
	if err != nil {
		return false, fmt.Errorf("获取授权记录失败: %v", err)
	}
	if authJSON == nil {
		return false, nil // 没有授权记录
	}

	// 解析授权记录
	var authorization Authorization
	err = json.Unmarshal(authJSON, &authorization)
	if err != nil {
		return false, fmt.Errorf("解析授权记录失败: %v", err)
	}

	// 检查授权是否有效
	if !authorization.IsActive {
		return false, nil
	}

	// 检查授权是否在有效期内
	currentTime := time.Now().Unix()
	if authorization.StartTime <= currentTime && (authorization.EndTime == 0 || authorization.EndTime >= currentTime) {
		return true, nil
	}

	return false, nil
}

// LogAccess 记录数据访问事件
func (s *MedicalDataContract) LogAccess(ctx contractapi.TransactionContextInterface, dataID, operation string) error {
	// 获取调用者身份
	userID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取调用者身份失败: %v", err)
	}

	// 检查用户是否有权访问数据
	hasAccess, err := s.CheckAccess(ctx, dataID, userID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("访问被拒绝")
	}

	// 创建访问记录
	accessRecord := AccessRecord{
		DataID:    dataID,
		Accessor:  userID,
		Timestamp: time.Now().Unix(),
		Operation: operation,
	}

	// 生成访问记录的键
	accessKey := fmt.Sprintf("access_%s_%s_%d", dataID, userID, accessRecord.Timestamp)

	// 将访问记录转换为JSON并存储
	accessJSON, err := json.Marshal(accessRecord)
	if err != nil {
		return err
	}

	// 将访问记录写入账本
	err = ctx.GetStub().PutState(accessKey, accessJSON)
	if err != nil {
		return fmt.Errorf("存储访问记录失败: %v", err)
	}

	// 记录数据访问事件
	return ctx.GetStub().SetEvent("DataAccessed", accessJSON)
}

// GetMedicalData 获取医疗数据详情
func (s *MedicalDataContract) GetMedicalData(ctx contractapi.TransactionContextInterface, dataID string) (*MedicalData, error) {
	// 获取数据
	medicalDataJSON, err := ctx.GetStub().GetState(dataID)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败: %v", err)
	}
	if medicalDataJSON == nil {
		return nil, fmt.Errorf("数据ID %s 不存在", dataID)
	}

	// 解析数据
	var medicalData MedicalData
	err = json.Unmarshal(medicalDataJSON, &medicalData)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return &medicalData, nil
}

// GetDataByOwner 获取用户拥有的所有数据
func (s *MedicalDataContract) GetDataByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*MedicalData, error) {
	// 创建复合键查询
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("owner~dataId", []string{owner})
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer resultIterator.Close()

	// 收集查询结果
	var medicalDataList []*MedicalData
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代查询结果失败: %v", err)
		}

		// 从复合键中提取数据ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, fmt.Errorf("分解复合键失败: %v", err)
		}
		dataID := compositeKeyParts[1]

		// 获取数据详情
		medicalData, err := s.GetMedicalData(ctx, dataID)
		if err != nil {
			return nil, err
		}

		medicalDataList = append(medicalDataList, medicalData)
	}

	return medicalDataList, nil
}

// GetAuthorizedData 获取用户被授权访问的所有数据
func (s *MedicalDataContract) GetAuthorizedData(ctx contractapi.TransactionContextInterface, userID string) ([]*MedicalData, error) {
	// 创建复合键查询
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("authorization", []string{userID})
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer resultIterator.Close()

	// 收集查询结果
	var medicalDataList []*MedicalData
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代查询结果失败: %v", err)
		}

		// 解析授权记录
		var authorization Authorization
		err = json.Unmarshal(queryResponse.Value, &authorization)
		if err != nil {
			return nil, fmt.Errorf("解析授权记录失败: %v", err)
		}

		// 检查授权是否有效
		if !authorization.IsActive {
			continue
		}

		// 检查授权是否在有效期内
		currentTime := time.Now().Unix()
		if authorization.StartTime <= currentTime && (authorization.EndTime == 0 || authorization.EndTime >= currentTime) {
			// 获取数据详情
			medicalData, err := s.GetMedicalData(ctx, authorization.DataID)
			if err != nil {
				return nil, err
			}

			medicalDataList = append(medicalDataList, medicalData)
		}
	}

	return medicalDataList, nil
}

// UpdateEthereumTxID 更新数据的以太坊交易ID
func (s *MedicalDataContract) UpdateEthereumTxID(ctx contractapi.TransactionContextInterface, dataID, ethereumTxID string) error {
	// 获取数据
	medicalData, err := s.GetMedicalData(ctx, dataID)
	if err != nil {
		return err
	}

	// 检查调用者是否是数据所有者
	creator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("获取调用者身份失败: %v", err)
	}

	if medicalData.Owner != creator {
		return fmt.Errorf("只有数据所有者才能更新以太坊交易ID")
	}

	// 更新以太坊交易ID
	medicalData.EthereumTxID = ethereumTxID

	// 将更新后的数据转换为JSON并存储
	medicalDataJSON, err := json.Marshal(medicalData)
	if err != nil {
		return err
	}

	// 将更新后的数据写入账本
	return ctx.GetStub().PutState(dataID, medicalDataJSON)
}

// DataExists 检查数据是否存在
func (s *MedicalDataContract) DataExists(ctx contractapi.TransactionContextInterface, dataID string) (bool, error) {
	medicalDataJSON, err := ctx.GetStub().GetState(dataID)
	if err != nil {
		return false, fmt.Errorf("获取数据失败: %v", err)
	}

	return medicalDataJSON != nil, nil
}

// 主函数
func main() {
	contract := new(MedicalDataContract)

	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		fmt.Printf("创建医疗数据链码失败: %v\n", err)
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("启动医疗数据链码失败: %v\n", err)
	}
}