package models

import (
	"time"
)

// TransferRecord 跨链转移记录
type TransferRecord struct {
	ID              string    `json:"id"`
	SourceID        string    `json:"sourceId"`        // 源数据ID
	TargetID        string    `json:"targetId"`        // 目标数据ID
	SourceChain     string    `json:"sourceChain"`     // 源区块链
	TargetChain     string    `json:"targetChain"`     // 目标区块链
	Timestamp       time.Time `json:"timestamp"`       // 转移时间
	Status          string    `json:"status"`          // 转移状态: "pending", "completed", "failed"
	ErrorMessage    string    `json:"errorMessage"`    // 错误信息（如果有）
	TransactionHash string    `json:"transactionHash"` // 区块链交易哈希
	Verified        bool      `json:"verified"`        // 数据完整性是否已验证
}

// TransferRequest 跨链转移请求
type TransferRequest struct {
	DataID      string `json:"dataId" binding:"required"`      // 要转移的数据ID
	SourceChain string `json:"sourceChain" binding:"required"` // 源区块链
	TargetChain string `json:"targetChain" binding:"required"` // 目标区块链
}

// TransferResponse 跨链转移响应
type TransferResponse struct {
	TransferID      string `json:"transferId"`      // 转移记录ID
	SourceID        string `json:"sourceId"`        // 源数据ID
	TargetID        string `json:"targetId"`        // 目标数据ID（如果已完成）
	Status          string `json:"status"`          // 转移状态
	TransactionHash string `json:"transactionHash"` // 区块链交易哈希（如果已完成）
}

// TransferHistoryResponse 转移历史响应
type TransferHistoryResponse struct {
	Records []TransferRecord `json:"records"`
	Total   int              `json:"total"`
}
