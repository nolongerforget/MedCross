package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"medcross/models"
	"medcross/services"
)

// DataController 处理医疗数据相关请求
type DataController struct {
	dataService    *services.DataService
	gatewayService *services.GatewayService
}

// NewDataController 创建新的数据控制器
func NewDataController(dataService *services.DataService, gatewayService *services.GatewayService) *DataController {
	return &DataController{
		dataService:    dataService,
		gatewayService: gatewayService,
	}
}

// QueryData 处理数据查询请求
func (dc *DataController) QueryData(c *gin.Context) {
	var query models.MedicalDataQuery

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的查询参数"})
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 调用跨链网关服务进行查询
	result, err := dc.gatewayService.QueryData(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UploadData 处理数据上传请求
func (dc *DataController) UploadData(c *gin.Context) {
	var uploadData models.MedicalDataUpload

	// 绑定请求数据
	if err := c.ShouldBindJSON(&uploadData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 生成唯一ID
	dataID := uuid.New().String()

	// 处理文件上传
	dataHash, err := dc.dataService.StoreFile(uploadData.File, uploadData.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储文件失败"})
		return
	}

	// 准备元数据
	metadata := map[string]string{
		"fileName":    uploadData.FileName,
		"description": uploadData.Description,
		"uploadedBy":  userID.(string),
		"fileSize":    strconv.Itoa(len(uploadData.File)),
	}

	// 创建医疗数据记录
	medicalData := models.MedicalData{
		ID:        dataID,
		Owner:     userID.(string),
		DataHash:  dataHash,
		DataType:  uploadData.DataType,
		Metadata:  dc.dataService.MapToJSON(metadata),
		Timestamp: time.Now(),
		Keywords:  uploadData.Keywords,
		Chain:     uploadData.TargetChain,
	}

	// 上传到区块链
	err = dc.gatewayService.UploadData(medicalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到区块链失败"})
		return
	}

	// 保存到本地数据库
	err = dc.dataService.SaveData(medicalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存数据失败"})
		return
	}

	c.JSON(http.StatusCreated, models.UploadResponse{
		ID:       dataID,
		Message:  "数据上传成功",
		DataHash: dataHash,
		Chain:    uploadData.TargetChain,
	})
}

// GetDataTypes 获取数据类型列表
func (dc *DataController) GetDataTypes(c *gin.Context) {
	// 返回预定义的数据类型列表
	dataTypes := []string{
		"影像数据",
		"电子病历",
		"基因组数据",
		"处方数据",
		"检验报告",
	}

	c.JSON(http.StatusOK, dataTypes)
}

// GetDataDetail 获取数据详情
func (dc *DataController) GetDataDetail(c *gin.Context) {
	// 获取数据ID
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少数据ID"})
		return
	}

	// 从本地数据库获取数据
	data, err := dc.dataService.GetDataByID(dataID)
	if err != nil {
		// 如果本地数据库没有，尝试从区块链获取
		data, err = dc.gatewayService.GetDataByID(dataID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "数据不存在"})
			return
		}
	}

	// 获取文件内容信息
	fileInfo, err := dc.dataService.GetFileInfo(data.DataHash)
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		// 继续处理，但不包含文件内容信息
	}

	// 解析元数据
	metadata := dc.dataService.JSONToMap(data.Metadata)

	// 构建响应
	response := gin.H{
		"id":        data.ID,
		"owner":     data.Owner,
		"dataHash":  data.DataHash,
		"dataType":  data.DataType,
		"metadata":  metadata,
		"timestamp": data.Timestamp,
		"keywords":  data.Keywords,
		"chain":     data.Chain,
	}

	// 如果有文件信息，添加到响应中
	if fileInfo != nil {
		response["content"] = fileInfo
	}

	c.JSON(http.StatusOK, response)
}

// GetStatistics 获取统计数据
func (dc *DataController) GetStatistics(c *gin.Context) {
	// 获取统计数据
	stats, err := dc.dataService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// SearchData 处理数据搜索请求
func (dc *DataController) SearchData(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	dataType := c.Query("dataType")
	chain := c.Query("chain")

	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 调用服务进行搜索
	result, err := dc.dataService.SearchDataByKeyword(keyword, dataType, chain, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索数据失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}
