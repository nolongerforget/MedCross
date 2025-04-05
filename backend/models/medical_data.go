package models

import (
	"time"

	"gorm.io/gorm"
)

// MedicalData 医疗数据模型
// 用于在以太坊和Fabric区块链之间共享的医疗数据结构
type MedicalData struct {
	DataID         string         `gorm:"column:data_id;primaryKey" json:"dataId"`         // 数据唯一标识符
	DataHash       string         `gorm:"column:data_hash;type:varchar(255);not null" json:"dataHash"`       // 数据哈希值（IPFS或其他存储系统的引用）
	DataType       string         `gorm:"column:data_type;type:varchar(50);not null" json:"dataType"`       // 数据类型（如：电子病历、影像数据、基因组数据等）
	Description    string         `gorm:"column:description;type:text" json:"description"`    // 数据描述
	Tags           []string       `gorm:"-" json:"tags"`           // 数据标签（通过关联表存储）
	Owner          string         `gorm:"column:owner_id;not null" json:"owner"`          // 数据所有者（用户ID）
	Timestamp      int64          `gorm:"column:timestamp;not null" json:"timestamp"`      // 上传时间戳
	IsConfidential bool           `gorm:"column:is_confidential;default:false" json:"isConfidential"` // 是否为机密数据
	EthereumTxID   string         `gorm:"column:ethereum_tx_id;type:varchar(100)" json:"ethereumTxId"`   // 以太坊链上对应的交易ID
	FabricTxID     string         `gorm:"column:fabric_tx_id;type:varchar(100)" json:"fabricTxId"`     // Fabric链上对应的交易ID
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Authorization 授权信息模型
// 用于管理医疗数据的访问权限
type Authorization struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	DataID         string         `gorm:"column:data_id;not null" json:"dataId"`         // 被授权的数据ID
	OwnerID        string         `gorm:"column:owner_id;not null" json:"ownerId"`        // 数据所有者ID
	AuthorizedUser string         `gorm:"column:authorized_user_id;not null" json:"authorizedUser"` // 被授权的用户ID
	StartTime      int64          `gorm:"column:start_time;not null" json:"startTime"`      // 授权开始时间
	EndTime        int64          `gorm:"column:end_time;not null" json:"endTime"`        // 授权结束时间（0表示永久授权）
	IsActive       bool           `gorm:"column:is_active;default:true" json:"isActive"`       // 授权是否有效
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// AccessRecord 访问记录模型
// 用于记录医疗数据的访问历史
type AccessRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DataID     string    `gorm:"column:data_id;not null" json:"dataId"`    // 访问的数据ID
	Accessor   string    `gorm:"column:accessor_id;not null" json:"accessor"`  // 访问者ID
	Timestamp  int64     `gorm:"column:timestamp;not null" json:"timestamp"` // 访问时间戳
	Operation  string    `gorm:"column:operation;type:varchar(20);not null" json:"operation"` // 操作类型（查看、下载等）
	ChainType  string    `gorm:"column:chain_type;type:varchar(10);not null" json:"chainType"` // 区块链类型（Ethereum或Fabric）
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// CrossChainReference 跨链引用模型
// 用于维护不同区块链上同一数据的引用关系
type CrossChainReference struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DataID       string         `gorm:"column:data_id;not null" json:"dataId"`       // 数据ID
	EthereumTxID string         `gorm:"column:ethereum_tx_id;type:varchar(100)" json:"ethereumTxId"` // 以太坊交易ID
	FabricTxID   string         `gorm:"column:fabric_tx_id;type:varchar(100)" json:"fabricTxId"`   // Fabric交易ID
	SyncStatus   string         `gorm:"column:sync_status;type:varchar(20);not null" json:"syncStatus"`   // 同步状态（已同步、同步中、同步失败）
	LastSyncTime int64          `gorm:"column:last_sync_time;not null" json:"lastSyncTime"` // 最后同步时间
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Tag 标签模型
type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
}

// DataTag 数据标签关联模型
type DataTag struct {
	DataID string `gorm:"primaryKey;column:data_id" json:"dataId"`
	TagID  uint   `gorm:"primaryKey;column:tag_id" json:"tagId"`
}