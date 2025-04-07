package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 跨链网关服务 - 负责协调以太坊和Fabric链上的数据查询

// 医疗数据结构 - 统一的数据格式
type MedicalData struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	DataHash  string    `json:"dataHash"`
	DataType  string    `json:"dataType"`
	Metadata  string    `json:"metadata"`
	Timestamp time.Time `json:"timestamp"`
	Keywords  string    `json:"keywords"`
	Chain     string    `json:"chain"` // 标识数据来源的区块链: "ethereum" 或 "fabric"
}

// 查询结果结构
type QueryResult struct {
	TotalCount int           `json:"totalCount"`
	Data       []MedicalData `json:"data"`
	Errors     []string      `json:"errors,omitempty"`
}

// 模拟以太坊区块链接口
func queryEthereumData(keyword string, dataType string) ([]MedicalData, error) {
	// 这里应该是实际调用以太坊节点的代码
	// 为了演示，我们返回一些模拟数据
	log.Printf("查询以太坊数据: 关键词=%s, 类型=%s", keyword, dataType)

	// 模拟数据
	result := []MedicalData{
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
	var filtered []MedicalData
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

	return filtered, nil
}

// 模拟Fabric区块链接口
func queryFabricData(keyword string, dataType string) ([]MedicalData, error) {
	// 这里应该是实际调用Fabric链码的代码
	// 为了演示，我们返回一些模拟数据
	log.Printf("查询Fabric数据: 关键词=%s, 类型=%s", keyword, dataType)

	// 模拟数据
	result := []MedicalData{
		{
			ID:        "fab-001",
			Owner:     "user1",
			DataHash:  "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ",
			DataType:  "基因组数据",
			Metadata:  `{"patientId":"P98765","hospital":"医学研究中心","project":"癌症基因研究"}`,
			Timestamp: time.Now().Add(-12 * time.Hour),
			Keywords:  "基因,癌症,研究",
			Chain:     "fabric",
		},
		{
			ID:        "fab-002",
			Owner:     "user2",
			DataHash:  "QmT8CUmNPMYGe8P9G2XKZHUuWaq9ZqCTGGYVqx57FuLSdT",
			DataType:  "影像数据",
			Metadata:  `{"patientId":"P24680","hospital":"第三医院","department":"神经外科"}`,
			Timestamp: time.Now().Add(-36 * time.Hour),
			Keywords:  "脑部,MRI,影像",
			Chain:     "fabric",
		},
	}

	// 根据关键词和类型筛选
	var filtered []MedicalData
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

	return filtered, nil
}

// 跨链查询处理函数
func crossChainQuery(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	dataType := c.Query("dataType")
	chainSource := c.Query("chain") // 可选值: "all", "ethereum", "fabric"

	if chainSource == "" {
		chainSource = "all" // 默认查询所有链
	}

	var ethData, fabricData []MedicalData
	var ethErr, fabricErr error
	var errors []string

	// 根据指定的链源进行查询
	if chainSource == "all" || chainSource == "ethereum" {
		ethData, ethErr = queryEthereumData(keyword, dataType)
		if ethErr != nil {
			errors = append(errors, fmt.Sprintf("以太坊查询错误: %v", ethErr))
		}
	}

	if chainSource == "all" || chainSource == "fabric" {
		fabricData, fabricErr = queryFabricData(keyword, dataType)
		if fabricErr != nil {
			errors = append(errors, fmt.Sprintf("Fabric查询错误: %v", fabricErr))
		}
	}

	// 合并结果
	allData := append(ethData, fabricData...)

	// 构建响应
	result := QueryResult{
		TotalCount: len(allData),
		Data:       allData,
		Errors:     errors,
	}

	c.JSON(http.StatusOK, result)
}

// 上传数据到指定区块链
func uploadData(c *gin.Context) {
	// 解析请求体
	var data struct {
		DataHash  string `json:"dataHash"`
		DataType  string `json:"dataType"`
		Metadata  string `json:"metadata"`
		Keywords  string `json:"keywords"`
		Chain     string `json:"chain"` // 目标区块链: "ethereum" 或 "fabric"
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证必填字段
	if data.DataHash == "" || data.DataType == "" || data.Chain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必填字段"})
		return
	}

	// 验证目标链
	if data.Chain != "ethereum" && data.Chain != "fabric" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标区块链，必须是 'ethereum' 或 'fabric'"})
		return
	}

	// 这里应该是实际调用区块链的代码
	// 为了演示，我们只记录日志并返回成功
	log.Printf("上传数据到%s链: %+v", data.Chain, data)

	// 生成模拟ID
	var id string
	if data.Chain == "ethereum" {
		id = fmt.Sprintf("eth-%d", time.Now().Unix())
	} else {
		id = fmt.Sprintf("fab-%d", time.Now().Unix())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
		"message": fmt.Sprintf("数据已成功上传到%s链", data.Chain),
	})
}

// 获取数据类型列表
func getDataTypes(c *gin.Context) {
	// 在实际应用中，这些类型可能来自区块链或数据库
	dataTypes := []string{
		"影像数据",
		"电子病历",
		"基因组数据",
		"处方数据",
		"检验报告",
	}

	c.JSON(http.StatusOK, gin.H{
		"dataTypes": dataTypes,
	})
}

// 获取单个数据详情
func getDataDetail(c *gin.Context) {
	id := c.Param("id")

	// 解析ID前缀以确定数据来源
	var data MedicalData
	var err error

	if strings.HasPrefix(id, "eth-") {
		// 从以太坊获取数据
		// 这里应该是实际调用以太坊的代码
		// 为了演示，我们返回模拟数据
		data = MedicalData{
			ID:        id,
			Owner:     "0x1234567890abcdef",
			DataHash:  "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o",
			DataType:  "影像数据",
			Metadata:  `{"patientId":"P12345","hospital":"协和医院","department":"放射科","description":"这是一份详细的CT扫描数据，显示患者肺部有轻微炎症"}`,
			Timestamp: time.Now().Add(-24 * time.Hour),
			Keywords:  "肺部,CT,影像",
			Chain:     "ethereum",
		}
	} else if strings.HasPrefix(id, "fab-") {
		// 从Fabric获取数据
		// 这里应该是实际调用Fabric的代码
		// 为了演示，我们返回模拟数据
		data = MedicalData{
			ID:        id,
			Owner:     "user1",
			DataHash:  "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ",
			DataType:  "基因组数据",
			Metadata:  `{"patientId":"P98765","hospital":"医学研究中心","project":"癌症基因研究","description":"这是一份癌症患者的基因测序数据，用于精准医疗研究"}`,
			Timestamp: time.Now().Add(-12 * time.Hour),
			Keywords:  "基因,癌症,研究",
			Chain:     "fabric",
		}
	} else {
		err = fmt.Errorf("无效的数据ID格式")
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("获取数据失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, data)
}

func main() {
	// 加载环境变量
	godotenv.Load()

	// 设置日志
	log.SetOutput(os.Stdout)

	// 创建Gin路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 定义API路由
	r.GET("/api/query", crossChainQuery)         // 跨链查询
	r.POST("/api/upload", uploadData)            // 上传数据
	r.GET("/api/data-types", getDataTypes)       // 获取数据类型列表
	r.GET("/api/data/:id", getDataDetail)        // 获取数据详情

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默认端口
	}

	// 启动服务器
	log.Printf("跨链网关服务启动在 :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}