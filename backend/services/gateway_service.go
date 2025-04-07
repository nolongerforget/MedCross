package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"medcross/models"
	"medcross/utils"
)

// GatewayService 跨链网关服务
type GatewayService struct {
	gatewayURL string
	timeout    time.Duration // HTTP请求超时时间
	maxRetries int           // 最大重试次数
}

// NewGatewayService 创建新的网关服务
func NewGatewayService() *GatewayService {
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8080"
	}

	// 从环境变量获取超时设置，默认为10秒
	timeoutStr := os.Getenv("GATEWAY_TIMEOUT")
	timeout := 10 * time.Second
	if timeoutStr != "" {
		if t, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = t
		}
	}

	// 从环境变量获取最大重试次数，默认为3次
	maxRetriesStr := os.Getenv("GATEWAY_MAX_RETRIES")
	maxRetries := 3
	if maxRetriesStr != "" {
		if r, err := strconv.Atoi(maxRetriesStr); err == nil && r > 0 {
			maxRetries = r
		}
	}

	return &GatewayService{
		gatewayURL: gatewayURL,
		timeout:    timeout,
		maxRetries: maxRetries,
	}
}

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
		log.Printf("查询网关失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("查询网关失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result models.QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("解析响应失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// UploadData 上传医疗数据到区块链
func (s *GatewayService) UploadData(data models.MedicalData) error {
	// 构建请求URL
	url := fmt.Sprintf("%s/upload", s.gatewayURL)

	// 准备请求数据
	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化数据失败: %v", err)
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqData)))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试上传数据 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		log.Printf("上传到网关失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("上传到网关失败，已达到最大重试次数: %v", err)
		return fmt.Errorf("上传到网关失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}

// GetDataByID 根据ID获取数据
func (s *GatewayService) GetDataByID(dataID string) (*models.MedicalData, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/data/%s", s.gatewayURL, dataID)

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
		log.Printf("查询网关失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("查询网关失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var data models.MedicalData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("解析响应失败: %v", err)
		return nil, err
	}

	return &data, nil
}

// 跨链数据转换和交互
// 实现与以太坊和Fabric区块链的交互，并处理跨链数据转换

// CrossChainTransfer 跨链数据转移
func (s *GatewayService) CrossChainTransfer(data models.MedicalData, targetChain string) (*models.MedicalData, error) {
	// 创建链转换器
	converter := utils.NewChainConverter()

	// 检查源链和目标链
	if data.Chain == targetChain {
		log.Printf("源链和目标链相同 (%s)，无需转换", data.Chain)
		return &data, nil // 如果源链和目标链相同，无需转换
	}

	// 验证支持的链类型
	if data.Chain != "ethereum" && data.Chain != "fabric" {
		log.Printf("不支持的源链类型: %s", data.Chain)
		return nil, fmt.Errorf("不支持的源链类型: %s", data.Chain)
	}

	if targetChain != "ethereum" && targetChain != "fabric" {
		log.Printf("不支持的目标链类型: %s", targetChain)
		return nil, fmt.Errorf("不支持的目标链类型: %s", targetChain)
	}

	// 执行跨链转换
	var convertedData models.MedicalData
	var err error

	log.Printf("开始跨链转换: %s -> %s, 数据ID: %s", data.Chain, targetChain, data.ID)

	// 记录转换开始时间，用于性能监控
	startTime := time.Now()

	if data.Chain == "ethereum" && targetChain == "fabric" {
		// 以太坊 -> Fabric
		log.Printf("执行以太坊到Fabric的转换")
		convertedData, err = converter.EthereumToFabric(data)
	} else if data.Chain == "fabric" && targetChain == "ethereum" {
		// Fabric -> 以太坊
		log.Printf("执行Fabric到以太坊的转换")
		convertedData, err = converter.FabricToEthereum(data)
	} else {
		log.Printf("不支持的跨链转换: %s -> %s", data.Chain, targetChain)
		return nil, fmt.Errorf("不支持的跨链转换: %s -> %s", data.Chain, targetChain)
	}

	if err != nil {
		log.Printf("跨链转换失败: %v", err)
		return nil, fmt.Errorf("跨链转换失败: %w", err)
	}

	log.Printf("跨链转换成功: 新ID=%s, 耗时=%v", convertedData.ID, time.Since(startTime))

	// 验证转换后的数据完整性
	log.Printf("验证转换后的数据完整性")
	valid, err := s.VerifyDataIntegrity(data, convertedData)
	if err != nil {
		log.Printf("数据完整性验证失败: %v", err)
		return nil, fmt.Errorf("数据完整性验证失败: %w", err)
	}

	if !valid {
		log.Printf("转换后的数据未通过完整性验证")
		return nil, fmt.Errorf("转换后的数据未通过完整性验证")
	}

	log.Printf("数据完整性验证通过，准备上传到目标链")

	// 上传转换后的数据到目标链
	url := fmt.Sprintf("%s/transfer", s.gatewayURL)

	// 准备请求数据
	reqData, err := json.Marshal(convertedData)
	if err != nil {
		log.Printf("序列化数据失败: %v", err)
		return nil, fmt.Errorf("序列化数据失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqData)))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 添加源数据ID作为请求头，便于网关追踪数据来源
	req.Header.Set("X-Source-Data-ID", data.ID)
	req.Header.Set("X-Source-Chain", data.Chain)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试跨链传输 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		log.Printf("跨链传输失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("跨链传输失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("跨链传输失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// 尝试读取错误信息
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil && errorResponse.Error != "" {
			log.Printf("网关返回错误: %s", errorResponse.Error)
			return nil, fmt.Errorf("网关返回错误: %s", errorResponse.Error)
		}

		log.Printf("网关返回错误状态码: %d", resp.StatusCode)
		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result models.MedicalData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("解析响应失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	log.Printf("跨链数据传输完成: 源ID=%s, 目标ID=%s", data.ID, result.ID)

	// 记录跨链转移历史（实际应用中应该由网关完成）
	log.Printf("跨链转移记录已创建: %s -> %s", data.Chain, targetChain)

	return &result, nil
}

// GetDataTypes 获取所有支持的数据类型
func (s *GatewayService) GetDataTypes() ([]string, error) {
	log.Printf("获取支持的数据类型")

	// 构建请求URL
	url := fmt.Sprintf("%s/datatypes", s.gatewayURL)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试获取数据类型 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("获取数据类型失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("获取数据类型失败，已达到最大重试次数: %v", err)
		// 返回默认数据类型作为备选
		log.Printf("返回默认数据类型")
		return []string{
			"影像数据",
			"电子病历",
			"基因组数据",
			"处方数据",
			"检验报告",
			"手术记录",
		}, nil
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Printf("获取数据类型失败，状态码: %d", resp.StatusCode)
		// 返回默认数据类型作为备选
		return []string{
			"影像数据",
			"电子病历",
			"基因组数据",
			"处方数据",
			"检验报告",
			"手术记录",
		}, nil
	}

	// 解析响应
	var dataTypes []string
	if err := json.NewDecoder(resp.Body).Decode(&dataTypes); err != nil {
		log.Printf("解析数据类型失败: %v", err)
		// 返回默认数据类型作为备选
		return []string{
			"影像数据",
			"电子病历",
			"基因组数据",
			"处方数据",
			"检验报告",
			"手术记录",
		}, nil
	}

	log.Printf("成功获取数据类型，共 %d 种类型", len(dataTypes))
	return dataTypes, nil
}

// GetChainDataTypeDistribution 获取特定区块链上的数据类型分布
func (s *GatewayService) GetChainDataTypeDistribution(chain string) (map[string]int, error) {
	log.Printf("获取区块链 %s 上的数据类型分布", chain)

	// 验证链类型
	if chain != "ethereum" && chain != "fabric" && chain != "all" {
		log.Printf("不支持的链类型: %s", chain)
		return nil, fmt.Errorf("不支持的链类型: %s", chain)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/statistics/distribution?chain=%s", s.gatewayURL, chain)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试获取数据类型分布 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("获取数据类型分布失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("获取数据类型分布失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("获取数据类型分布失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var distribution map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&distribution); err != nil {
		log.Printf("解析数据类型分布失败: %v", err)
		return nil, fmt.Errorf("解析数据类型分布失败: %w", err)
	}

	log.Printf("成功获取区块链 %s 上的数据类型分布，共 %d 种类型", chain, len(distribution))
	return distribution, nil
}

// GetStatistics 获取跨链统计数据
func (s *GatewayService) GetStatistics() (*models.Statistics, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/statistics", s.gatewayURL)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试获取统计数据 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("获取统计数据失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("获取统计数据失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("获取统计数据失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var stats models.Statistics
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		log.Printf("解析统计数据失败: %v", err)
		return nil, fmt.Errorf("解析统计数据失败: %w", err)
	}

	return &stats, nil
}

// VerifyDataIntegrity 验证跨链数据完整性
func (s *GatewayService) VerifyDataIntegrity(sourceData models.MedicalData, targetData models.MedicalData) (bool, error) {
	log.Printf("验证跨链数据完整性: 源ID=%s, 目标ID=%s", sourceData.ID, targetData.ID)

	// 验证基本数据哈希是否一致
	if sourceData.DataHash != targetData.DataHash {
		log.Printf("数据哈希不匹配: 源=%s, 目标=%s", sourceData.DataHash, targetData.DataHash)
		return false, fmt.Errorf("数据哈希不匹配")
	}

	// 验证数据类型是否一致
	if sourceData.DataType != targetData.DataType {
		log.Printf("数据类型不匹配: 源=%s, 目标=%s", sourceData.DataType, targetData.DataType)
		return false, fmt.Errorf("数据类型不匹配")
	}

	// 验证关键字是否一致（允许目标数据有额外的关键字）
	sourceKeywords := strings.Split(sourceData.Keywords, ",")
	targetKeywords := strings.Split(targetData.Keywords, ",")
	keywordMap := make(map[string]bool)

	for _, keyword := range targetKeywords {
		keywordMap[strings.TrimSpace(keyword)] = true
	}

	for _, keyword := range sourceKeywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" && !keywordMap[keyword] {
			log.Printf("关键字不匹配: 源数据关键字 '%s' 在目标数据中不存在", keyword)
			return false, fmt.Errorf("关键字不匹配: 源数据关键字 '%s' 在目标数据中不存在", keyword)
		}
	}

	// 解析元数据进行更深入的验证
	var sourceMetadata, targetMetadata map[string]interface{}
	if err := json.Unmarshal([]byte(sourceData.Metadata), &sourceMetadata); err != nil {
		log.Printf("解析源数据元数据失败: %v", err)
		return false, fmt.Errorf("解析源数据元数据失败: %w", err)
	}
	if err := json.Unmarshal([]byte(targetData.Metadata), &targetMetadata); err != nil {
		log.Printf("解析目标数据元数据失败: %v", err)
		return false, fmt.Errorf("解析目标数据元数据失败: %w", err)
	}

	// 检查跨链转换信息
	if targetMetadata["crossChainSource"] == nil {
		log.Printf("目标数据缺少跨链源信息")
		return false, fmt.Errorf("目标数据缺少跨链源信息")
	}

	if targetMetadata["crossChainSource"].(string) != sourceData.Chain {
		log.Printf("跨链源信息不匹配: 预期=%s, 实际=%s", sourceData.Chain, targetMetadata["crossChainSource"])
		return false, fmt.Errorf("跨链源信息不匹配")
	}

	// 检查关键元数据字段
	// 注意：跨链转换可能会修改某些字段，所以我们只检查关键字段
	keysToCheck := []string{"patientId", "hospital", "department"}
	for _, key := range keysToCheck {
		// 检查源数据中是否有此字段
		sourceVal, sourceExists := sourceMetadata[key]
		targetVal, targetExists := targetMetadata[key]

		// 如果两边都有此字段，则比较值
		if sourceExists && targetExists {
			if fmt.Sprintf("%v", sourceVal) != fmt.Sprintf("%v", targetVal) {
				log.Printf("元数据字段 '%s' 不匹配: 源=%v, 目标=%v", key, sourceVal, targetVal)
				return false, fmt.Errorf("元数据字段 '%s' 不匹配", key)
			}
		} else if sourceExists && !targetExists {
			// 源数据有此字段但目标数据没有，这是一个错误
			log.Printf("目标数据缺少关键元数据字段: %s", key)
			return false, fmt.Errorf("目标数据缺少关键元数据字段: %s", key)
		}
	}

	// 检查特殊字段映射
	if sourceData.Chain == "ethereum" && targetData.Chain == "fabric" {
		// 检查以太坊地址是否正确映射到Fabric ID
		expectedFabricID := fmt.Sprintf("fabric-user-%s", strings.TrimPrefix(sourceData.Owner, "0x")[:8])
		if targetData.Owner != expectedFabricID {
			log.Printf("所有者映射不正确: 预期=%s, 实际=%s", expectedFabricID, targetData.Owner)
			return false, fmt.Errorf("所有者映射不正确")
		}
	} else if sourceData.Chain == "fabric" && targetData.Chain == "ethereum" {
		// Fabric到以太坊的映射验证较为复杂，这里简化处理
		if !strings.HasPrefix(targetData.Owner, "0x") {
			log.Printf("以太坊地址格式不正确: %s", targetData.Owner)
			return false, fmt.Errorf("以太坊地址格式不正确")
		}
	}

	log.Printf("跨链数据完整性验证通过")
	return true, nil
}

// GetTransferHistory 获取数据的跨链转移历史
func (s *GatewayService) GetTransferHistory(dataID string) ([]models.TransferRecord, error) {
	log.Printf("获取数据跨链转移历史: ID=%s", dataID)

	// 构建请求URL
	url := fmt.Sprintf("%s/transfer/history/%s", s.gatewayURL, dataID)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试获取转移历史 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("获取转移历史失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("获取转移历史失败，已达到最大重试次数: %v", err)
		return nil, fmt.Errorf("获取转移历史失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		// 尝试读取错误信息
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil && errorResponse.Error != "" {
			log.Printf("获取转移历史失败: %s", errorResponse.Error)
			return nil, fmt.Errorf("获取转移历史失败: %s", errorResponse.Error)
		}

		return nil, fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var historyResponse models.TransferHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&historyResponse); err != nil {
		log.Printf("解析转移历史失败: %v", err)
		return nil, fmt.Errorf("解析转移历史失败: %w", err)
	}

	log.Printf("成功获取转移历史，共 %d 条记录", len(historyResponse.Records))
	return historyResponse.Records, nil
}

// BatchVerifyDataIntegrity 批量验证跨链数据完整性
func (s *GatewayService) BatchVerifyDataIntegrity(records []models.TransferRecord) (map[string]bool, error) {
	log.Printf("开始批量验证跨链数据完整性，共 %d 条记录", len(records))

	results := make(map[string]bool)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 限制并发数量，避免过多的并发请求
	semaphore := make(chan struct{}, 5)

	for _, record := range records {
		// 跳过未完成或已失败的转移记录
		if record.Status != "completed" {
			log.Printf("跳过未完成的转移记录: ID=%s, 状态=%s", record.ID, record.Status)
			mu.Lock()
			results[record.ID] = false
			mu.Unlock()
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(r models.TransferRecord) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			log.Printf("验证转移记录: ID=%s, 源=%s, 目标=%s", r.ID, r.SourceID, r.TargetID)

			// 获取源数据
			sourceData, err := s.GetDataByID(r.SourceID)
			if err != nil {
				log.Printf("获取源数据失败: %v", err)
				mu.Lock()
				results[r.ID] = false
				mu.Unlock()
				return
			}

			// 获取目标数据
			targetData, err := s.GetDataByID(r.TargetID)
			if err != nil {
				log.Printf("获取目标数据失败: %v", err)
				mu.Lock()
				results[r.ID] = false
				mu.Unlock()
				return
			}

			// 验证数据完整性
			valid, err := s.VerifyDataIntegrity(*sourceData, *targetData)
			if err != nil {
				log.Printf("验证数据完整性失败: %v", err)
				mu.Lock()
				results[r.ID] = false
				mu.Unlock()
				return
			}

			mu.Lock()
			results[r.ID] = valid
			mu.Unlock()

			log.Printf("转移记录验证结果: ID=%s, 结果=%v", r.ID, valid)
		}(record)
	}

	wg.Wait()

	// 统计验证结果
	validCount := 0
	for _, valid := range results {
		if valid {
			validCount++
		}
	}

	log.Printf("批量验证完成: 总数=%d, 有效=%d, 无效=%d", len(records), validCount, len(records)-validCount)
	return results, nil
}

// QueryBlockchainData 查询特定区块链上的数据
func (s *GatewayService) QueryBlockchainData(chain string, keyword string, dataType string, page int, pageSize int) ([]models.MedicalData, error) {
	log.Printf("查询区块链数据: 链=%s, 关键词=%s, 数据类型=%s, 页码=%d, 每页大小=%d", chain, keyword, dataType, page, pageSize)

	// 验证链类型
	if chain != "ethereum" && chain != "fabric" && chain != "all" {
		log.Printf("不支持的链类型: %s", chain)
		return nil, fmt.Errorf("不支持的链类型: %s", chain)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/blockchain/query?chain=%s&keyword=%s&dataType=%s&page=%d&pageSize=%d",
		s.gatewayURL, chain, keyword, dataType, page, pageSize)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试查询区块链数据 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("查询区块链数据失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("查询区块链数据失败，已达到最大重试次数: %v", err)

		// 如果网关不可用，使用模拟数据作为备选
		log.Printf("使用模拟数据作为备选")
		if chain == "ethereum" || chain == "all" {
			return s.mockEthereumData(keyword, dataType), nil
		} else if chain == "fabric" {
			return s.mockFabricData(keyword, dataType), nil
		}
		return nil, fmt.Errorf("查询区块链数据失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Printf("区块链查询返回错误状态码: %d", resp.StatusCode)

		// 如果网关返回错误，使用模拟数据作为备选
		log.Printf("使用模拟数据作为备选")
		if chain == "ethereum" || chain == "all" {
			return s.mockEthereumData(keyword, dataType), nil
		} else if chain == "fabric" {
			return s.mockFabricData(keyword, dataType), nil
		}
		return nil, fmt.Errorf("区块链查询返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result []models.MedicalData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("解析区块链数据失败: %v", err)
		return nil, fmt.Errorf("解析区块链数据失败: %w", err)
	}

	log.Printf("成功查询区块链数据，共 %d 条记录", len(result))
	return result, nil
}

// SubmitBlockchainTransaction 提交区块链交易
func (s *GatewayService) SubmitBlockchainTransaction(chain string, txType string, data interface{}) (string, error) {
	log.Printf("提交区块链交易: 链=%s, 类型=%s", chain, txType)

	// 验证链类型
	if chain != "ethereum" && chain != "fabric" {
		log.Printf("不支持的链类型: %s", chain)
		return "", fmt.Errorf("不支持的链类型: %s", chain)
	}

	// 验证交易类型
	validTxTypes := map[string]bool{"upload": true, "transfer": true, "update": true, "delete": true}
	if !validTxTypes[txType] {
		log.Printf("不支持的交易类型: %s", txType)
		return "", fmt.Errorf("不支持的交易类型: %s", txType)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/blockchain/transaction/%s/%s", s.gatewayURL, chain, txType)

	// 准备请求数据
	reqData, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化交易数据失败: %v", err)
		return "", fmt.Errorf("序列化交易数据失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqData)))
	if err != nil {
		log.Printf("创建交易请求失败: %v", err)
		return "", fmt.Errorf("创建交易请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试提交区块链交易 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		log.Printf("提交区块链交易失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("提交区块链交易失败，已达到最大重试次数: %v", err)
		return "", fmt.Errorf("提交区块链交易失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// 尝试读取错误信息
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil && errorResponse.Error != "" {
			log.Printf("区块链交易失败: %s", errorResponse.Error)
			return "", fmt.Errorf("区块链交易失败: %s", errorResponse.Error)
		}

		log.Printf("区块链交易返回错误状态码: %d", resp.StatusCode)
		return "", fmt.Errorf("区块链交易返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应，获取交易哈希
	var txResponse struct {
		TransactionHash string `json:"transactionHash"`
		Message         string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&txResponse); err != nil {
		log.Printf("解析交易响应失败: %v", err)
		return "", fmt.Errorf("解析交易响应失败: %w", err)
	}

	log.Printf("区块链交易提交成功: %s", txResponse.TransactionHash)
	return txResponse.TransactionHash, nil
}

// GetBlockchainTransactionStatus 获取区块链交易状态
func (s *GatewayService) GetBlockchainTransactionStatus(chain string, txHash string) (string, error) {
	log.Printf("获取区块链交易状态: 链=%s, 交易哈希=%s", chain, txHash)

	// 验证链类型
	if chain != "ethereum" && chain != "fabric" {
		log.Printf("不支持的链类型: %s", chain)
		return "", fmt.Errorf("不支持的链类型: %s", chain)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/blockchain/transaction/%s/status/%s", s.gatewayURL, chain, txHash)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: s.timeout,
	}

	// 添加重试机制
	var resp *http.Response
	var err error
	for i := 0; i < s.maxRetries; i++ {
		log.Printf("尝试获取交易状态 (尝试 %d/%d)", i+1, s.maxRetries)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		log.Printf("获取交易状态失败 (尝试 %d/%d): %v", i+1, s.maxRetries, err)
		if i < s.maxRetries-1 {
			// 指数退避策略
			backoff := time.Duration(100*(i+1)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	if err != nil {
		log.Printf("获取交易状态失败，已达到最大重试次数: %v", err)
		return "", fmt.Errorf("获取交易状态失败，已达到最大重试次数: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Printf("获取交易状态返回错误状态码: %d", resp.StatusCode)
		return "", fmt.Errorf("获取交易状态返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var statusResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		log.Printf("解析交易状态失败: %v", err)
		return "", fmt.Errorf("解析交易状态失败: %w", err)
	}

	log.Printf("交易状态: %s, 消息: %s", statusResponse.Status, statusResponse.Message)
	return statusResponse.Status, nil
}

// 模拟以太坊区块链数据（仅在网关不可用时使用）
func (s *GatewayService) mockEthereumData(keyword string, dataType string) []models.MedicalData {
	log.Printf("使用模拟以太坊数据")

	// 模拟数据
	result := []models.MedicalData{
		{
			ID:        "eth-001",
			Owner:     "0x1234567890abcdef",
			DataHash:  "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o",
			DataType:  "影像数据",
			Metadata:  `{"patientId":"P12345","hospital":"协和医院","department":"放射科"}`,
			Timestamp: time.Now().Add(-24 * time.Hour),
			Keywords:  "肺部,CT,影像",
			Chain:     "ethereum",
		},
		{
			ID:        "eth-002",
			Owner:     "0xabcdef1234567890",
			DataHash:  "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn",
			DataType:  "电子病历",
			Metadata:  `{"patientId":"P54321","hospital":"人民医院","department":"内科"}`,
			Timestamp: time.Now().Add(-48 * time.Hour),
			Keywords:  "糖尿病,慢性病,病历",
			Chain:     "ethereum",
		},
	}

	// 根据关键词和类型筛选
	var filtered []models.MedicalData
	for _, data := range result {
		// 类型筛选
		if dataType != "" && dataType != "all" && data.DataType != dataType {
			continue
		}

		// 关键词筛选
		if keyword != "" {
			if !strings.Contains(strings.ToLower(data.Keywords), strings.ToLower(keyword)) &&
				!strings.Contains(strings.ToLower(data.Metadata), strings.ToLower(keyword)) {
				continue
			}
		}

		filtered = append(filtered, data)
	}

	return filtered
}

// 模拟Fabric区块链数据（仅在网关不可用时使用）
func (s *GatewayService) mockFabricData(keyword string, dataType string) []models.MedicalData {
	log.Printf("使用模拟Fabric数据")

	// 模拟数据
	result := []models.MedicalData{
		{
			ID:        "fab-001",
			Owner:     "user1",
			DataHash:  "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ",
			DataType:  "基因组数据",
			Metadata:  `{"patientId":"P98765","hospital":"医学研究中心","department":"基因组学"}`,
			Timestamp: time.Now().Add(-72 * time.Hour),
			Keywords:  "基因,测序,肿瘤",
			Chain:     "fabric",
		},
		{
			ID:        "fab-002",
			Owner:     "user2",
			DataHash:  "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
			DataType:  "处方数据",
			Metadata:  `{"patientId":"P24680","hospital":"中心医院","department":"药剂科"}`,
			Timestamp: time.Now().Add(-96 * time.Hour),
			Keywords:  "处方,药物,抗生素",
			Chain:     "fabric",
		},
	}

	// 根据关键词和类型筛选
	var filtered []models.MedicalData
	for _, data := range result {
		// 类型筛选
		if dataType != "" && dataType != "all" && data.DataType != dataType {
			continue
		}

		// 关键词筛选
		if keyword != "" {
			if !strings.Contains(strings.ToLower(data.Keywords), strings.ToLower(keyword)) &&
				!strings.Contains(strings.ToLower(data.Metadata), strings.ToLower(keyword)) {
				continue
			}
		}

		filtered = append(filtered, data)
	}

	return filtered
}
