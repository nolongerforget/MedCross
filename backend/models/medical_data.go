package models

import (
	"time"
)

// MedicalData 医疗数据模型
type MedicalData struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`      // 数据所有者的用户ID
	DataHash  string    `json:"dataHash"`   // IPFS或其他存储系统的哈希值
	DataType  string    `json:"dataType"`   // 数据类型（如：影像数据、电子病历等）
	Metadata  string    `json:"metadata"`   // JSON格式的元数据
	Timestamp time.Time `json:"timestamp"`  // 上传时间戳
	Keywords  string    `json:"keywords"`   // 关键词，用于搜索，以逗号分隔
	Chain     string    `json:"chain"`      // 标识数据来源的区块链: "ethereum" 或 "fabric"
}

// MedicalDataUpload 医疗数据上传请求
type MedicalDataUpload struct {
	File        []byte `json:"file" binding:"required"`      // 文件数据（Base64编码）
	FileName    string `json:"fileName" binding:"required"`  // 文件名
	DataType    string `json:"dataType" binding:"required"`  // 数据类型
	Description string `json:"description" binding:"required"` // 数据描述
	Keywords    string `json:"keywords"`                     // 关键词，用逗号分隔
	TargetChain string `json:"targetChain" binding:"required"` // 目标区块链
}

// MedicalDataQuery 医疗数据查询请求
type MedicalDataQuery struct {
	Keyword    string `form:"keyword"`    // 搜索关键词
	DataType   string `form:"dataType"`   // 数据类型筛选
	Chain      string `form:"chain"`      // 区块链筛选
	StartDate  string `form:"startDate"`  // 开始日期
	EndDate    string `form:"endDate"`    // 结束日期
	SortBy     string `form:"sortBy"`     // 排序方式
	Page       int    `form:"page"`       // 页码
	PageSize   int    `form:"pageSize"`   // 每页大小
}

// QueryResult 查询结果
type QueryResult struct {
	TotalCount int           `json:"totalCount"`
	Data       []MedicalData `json:"data"`
	Errors     []string      `json:"errors,omitempty"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	DataHash string `json:"dataHash"`
	Chain    string `json:"chain"`
}

// Statistics 统计数据
type Statistics struct {
	TotalRecords         int            `json:"totalRecords"`
	EthereumRecords      int            `json:"ethereumRecords"`
	FabricRecords        int            `json:"fabricRecords"`
	DataTypeDistribution map[string]int `json:"dataTypeDistribution"`
	UploadTrend          []DailyUpload  `json:"uploadTrend"`
}

// DailyUpload 每日上传统计
type DailyUpload struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}