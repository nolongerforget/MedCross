package services

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"medcross/models"
)

// DataService 数据服务
type DataService struct {
	// 在实际应用中，这里应该有数据库连接和文件存储服务
	// 为了演示，我们使用内存存储
	data map[string]*models.MedicalData
}

// NewDataService 创建新的数据服务
func NewDataService() *DataService {
	return &DataService{
		data: make(map[string]*models.MedicalData),
	}
}

// SaveData 保存医疗数据
func (s *DataService) SaveData(data models.MedicalData) error {
	// 检查数据ID是否已存在
	if _, exists := s.data[data.ID]; exists {
		return errors.New("数据ID已存在")
	}

	// 保存数据
	s.data[data.ID] = &data

	return nil
}

// GetDataByID 根据ID获取数据
func (s *DataService) GetDataByID(dataID string) (*models.MedicalData, error) {
	data, exists := s.data[dataID]
	if !exists {
		return nil, errors.New("数据不存在")
	}

	return data, nil
}

// StoreFile 存储文件并返回哈希值
func (s *DataService) StoreFile(fileData []byte, fileName string) (string, error) {
	// 在实际应用中，这里应该将文件存储到IPFS或其他存储系统
	// 为了演示，我们生成一个模拟的哈希值
	rand.Seed(time.Now().UnixNano())
	hash := "Qm"
	for i := 0; i < 44; i++ {
		hash += string("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[rand.Intn(62)])
	}

	log.Printf("存储文件: %s, 大小: %d 字节, 哈希: %s", fileName, len(fileData), hash)

	return hash, nil
}

// GetFileInfo 获取文件信息
func (s *DataService) GetFileInfo(dataHash string) (map[string]interface{}, error) {
	// 在实际应用中，这里应该从IPFS或其他存储系统获取文件信息
	// 为了演示，我们返回一些模拟数据
	fileInfo := map[string]interface{}{
		"fileUrl":      "https://example.com/files/" + dataHash,
		"thumbnailUrl": "https://example.com/thumbnails/" + dataHash + ".jpg",
		"fileSize":     "15.2 MB",
		"format":       "DICOM",
	}

	return fileInfo, nil
}

// GetStatistics 获取统计数据
func (s *DataService) GetStatistics() (*models.Statistics, error) {
	// 在实际应用中，这里应该从数据库获取统计数据
	// 为了演示，我们返回一些模拟数据
	totalRecords := len(s.data)
	ethereumRecords := 0
	fabricRecords := 0
	dataTypeDistribution := make(map[string]int)

	// 统计各类数据
	for _, data := range s.data {
		if data.Chain == "ethereum" {
			ethereumRecords++
		} else if data.Chain == "fabric" {
			fabricRecords++
		}

		dataTypeDistribution[data.DataType]++
	}

	// 如果没有数据，使用模拟数据
	if totalRecords == 0 {
		totalRecords = 156
		ethereumRecords = 98
		fabricRecords = 58
		dataTypeDistribution = map[string]int{
			"影像数据":  42,
			"电子病历":  65,
			"基因组数据": 18,
			"处方数据":  25,
			"检验报告":  6,
		}
	}

	// 生成上传趋势数据
	uploadTrend := []models.DailyUpload{
		{Date: "2024-06-01", Count: 12},
		{Date: "2024-06-02", Count: 8},
		{Date: "2024-06-03", Count: 15},
		{Date: "2024-06-04", Count: 10},
		{Date: "2024-06-05", Count: 14},
	}

	// 构建统计数据
	stats := &models.Statistics{
		TotalRecords:         totalRecords,
		EthereumRecords:      ethereumRecords,
		FabricRecords:        fabricRecords,
		DataTypeDistribution: dataTypeDistribution,
		UploadTrend:          uploadTrend,
	}

	return stats, nil
}

// MapToJSON 将map转换为JSON字符串
func (s *DataService) MapToJSON(m map[string]string) string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		log.Printf("转换JSON失败: %v", err)
		return "{}"
	}

	return string(jsonBytes)
}

// JSONToMap 将JSON字符串转换为map
func (s *DataService) JSONToMap(jsonStr string) map[string]interface{} {
	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Printf("解析JSON失败: %v", err)
		return map[string]interface{}{}
	}

	return result
}

// SearchDataByKeyword 根据关键词搜索医疗数据
func (s *DataService) SearchDataByKeyword(keyword string, dataType string, chain string, page int, pageSize int) (*models.QueryResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 搜索结果
	var results []models.MedicalData

	// 遍历所有数据
	for _, data := range s.data {
		// 应用筛选条件
		if dataType != "" && data.DataType != dataType {
			continue
		}
		if chain != "" && data.Chain != chain {
			continue
		}

		// 如果没有关键词，则直接添加到结果中
		if keyword == "" {
			results = append(results, *data)
			continue
		}

		// 检查关键词是否匹配
		metadata := s.JSONToMap(data.Metadata)
		description, hasDesc := metadata["description"].(string)

		// 在关键词字段中搜索
		if strings.Contains(strings.ToLower(data.Keywords), strings.ToLower(keyword)) {
			results = append(results, *data)
			continue
		}

		// 在描述中搜索
		if hasDesc && strings.Contains(strings.ToLower(description), strings.ToLower(keyword)) {
			results = append(results, *data)
			continue
		}

		// 在元数据中搜索
		for k, v := range metadata {
			if strValue, ok := v.(string); ok {
				if strings.Contains(strings.ToLower(k), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(strValue), strings.ToLower(keyword)) {
					results = append(results, *data)
					break
				}
			}
		}
	}

	// 计算分页
	totalCount := len(results)
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= totalCount {
		// 页码超出范围，返回空结果
		return &models.QueryResult{
			TotalCount: totalCount,
			Data:       []models.MedicalData{},
		}, nil
	}

	if endIndex > totalCount {
		endIndex = totalCount
	}

	// 返回分页后的结果
	return &models.QueryResult{
		TotalCount: totalCount,
		Data:       results[startIndex:endIndex],
	}, nil
}
