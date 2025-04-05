package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"d.com/MedCross/backend/models"
	"d.com/MedCross/backend/services"
)

// 数据上传请求结构
type UploadDataRequest struct {
	DataHash      string   `json:"dataHash" binding:"required"`
	DataType      string   `json:"dataType" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	Tags          []string `json:"tags"`
	IsConfidential bool     `json:"isConfidential"`
}

// 数据授权请求结构
type AuthorizeDataRequest struct {
	AuthorizedUserID string `json:"authorizedUserId" binding:"required"`
	StartTime        int64  `json:"startTime"`
	EndTime          int64  `json:"endTime"`
}

// UploadData 处理医疗数据上传
func UploadData(c *gin.Context) {
	var req UploadDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 生成唯一的数据ID
	dataID := services.GenerateUniqueID()

	// 创建医疗数据模型
	medicalData := models.MedicalData{
		DataID:        dataID,
		DataHash:      req.DataHash,
		DataType:      req.DataType,
		Description:   req.Description,
		Tags:          req.Tags,
		Owner:         userID.(string),
		Timestamp:     time.Now().Unix(),
		IsConfidential: req.IsConfidential,
	}

	// 调用区块链服务上传数据
	ethereum := services.NewEthereumService()
	fabric := services.NewFabricService()

	// 首先上传到Fabric
	fabricTxID, err := fabric.UploadData(medicalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到Fabric失败: " + err.Error()})
		return
	}

	// 然后上传到以太坊
	ethereum.UploadData(medicalData, fabricTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到以太坊失败: " + err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "数据上传成功",
		"dataId":  dataID,
		"fabricTxId": fabricTxID,
	})
}

// ListData 获取用户拥有的所有数据
func ListData(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用以太坊服务获取用户数据
	ethereum := services.NewEthereumService()
	dataList, err := ethereum.GetUserData(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败: " + err.Error()})
		return
	}

	// 返回数据列表
	c.JSON(http.StatusOK, gin.H{
		"data": dataList,
	})
}

// SearchData 搜索医疗数据
func SearchData(c *gin.Context) {
	// 获取查询参数
	dataType := c.Query("dataType")
	keyword := c.Query("keyword")
	tag := c.Query("tag")

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用服务搜索数据
	ethereum := services.NewEthereumService()
	fabric := services.NewFabricService()

	// 从以太坊获取基本数据列表
	dataList, err := ethereum.SearchData(userID.(string), dataType, keyword, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索数据失败: " + err.Error()})
		return
	}

	// 对于每个数据，从Fabric获取详细信息
	for i, data := range dataList {
		detail, err := fabric.GetDataDetail(data.DataID)
		if err == nil && detail != nil {
			// 合并详细信息
			dataList[i].AccessCount = detail.AccessCount
			dataList[i].LastAccessTime = detail.LastAccessTime
		}
	}

	// 返回搜索结果
	c.JSON(http.StatusOK, gin.H{
		"data": dataList,
	})
}

// GetDataById 获取特定数据的详细信息
func GetDataById(c *gin.Context) {
	// 获取数据ID
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用以太坊服务检查访问权限
	ethereum := services.NewEthereumService()
	hasAccess, err := ethereum.CheckAccess(dataID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查访问权限失败: " + err.Error()})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该数据"})
		return
	}

	// 获取数据详情
	data, err := ethereum.GetDataById(dataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败: " + err.Error()})
		return
	}

	// 从Fabric获取更详细的信息
	fabric := services.NewFabricService()
	detail, err := fabric.GetDataDetail(dataID)
	if err == nil && detail != nil {
		// 合并详细信息
		data.AccessCount = detail.AccessCount
		data.LastAccessTime = detail.LastAccessTime
		data.AccessHistory = detail.AccessHistory
	}

	// 记录访问日志
	fabric.LogAccess(dataID, userID.(string), "view")

	// 返回数据详情
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// AuthorizeData 授权其他用户访问数据
func AuthorizeData(c *gin.Context) {
	// 获取数据ID
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据ID不能为空"})
		return
	}

	// 解析请求体
	var req AuthorizeDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 设置默认的开始时间和结束时间
	startTime := req.StartTime
	if startTime == 0 {
		startTime = time.Now().Unix()
	}

	// 调用以太坊服务授权数据
	ethereum := services.NewEthereumService()
	err := ethereum.GrantAccess(dataID, userID.(string), req.AuthorizedUserID, startTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "授权失败: " + err.Error()})
		return
	}

	// 同步授权到Fabric
	fabric := services.NewFabricService()
	err = fabric.GrantAccess(dataID, userID.(string), req.AuthorizedUserID, startTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步授权到Fabric失败: " + err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "授权成功",
		"dataId": dataID,
		"authorizedUserId": req.AuthorizedUserID,
	})
}

// RevokeAuthorization 撤销用户的数据访问权限
func RevokeAuthorization(c *gin.Context) {
	// 获取数据ID和被授权用户ID
	dataID := c.Param("id")
	authorizedUserID := c.Param("userId")

	if dataID == "" || authorizedUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据ID和用户ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用以太坊服务撤销授权
	ethereum := services.NewEthereumService()
	err := ethereum.RevokeAccess(dataID, userID.(string), authorizedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销授权失败: " + err.Error()})
		return
	}

	// 同步撤销授权到Fabric
	fabric := services.NewFabricService()
	err = fabric.RevokeAccess(dataID, userID.(string), authorizedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步撤销授权到Fabric失败: " + err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "撤销授权成功",
		"dataId": dataID,
		"authorizedUserId": authorizedUserID,
	})
}

// ListAuthorizations 获取数据的所有授权信息
func ListAuthorizations(c *gin.Context) {
	// 获取数据ID
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用以太坊服务获取授权列表
	ethereum := services.NewEthereumService()
	authorizations, err := ethereum.GetDataAuthorizations(dataID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取授权列表失败: " + err.Error()})
		return
	}

	// 返回授权列表
	c.JSON(http.StatusOK, gin.H{
		"dataId": dataID,
		"authorizations": authorizations,
	})
}

// GetAccessLogs 获取数据的访问记录
func GetAccessLogs(c *gin.Context) {
	// 获取数据ID
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用Fabric服务获取访问日志
	fabric := services.NewFabricService()
	accessLogs, err := fabric.GetAccessLogs(dataID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取访问日志失败: " + err.Error()})
		return
	}

	// 返回访问日志
	c.JSON(http.StatusOK, gin.H{
		"dataId": dataID,
		"accessLogs": accessLogs,
	})
}