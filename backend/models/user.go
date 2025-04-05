package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(100);not null" json:"password,omitempty"` // 密码不应在JSON响应中返回
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	FullName  string         `gorm:"type:varchar(100);not null" json:"fullName"`
	OrgName   string         `gorm:"type:varchar(100);not null" json:"orgName"` // 组织名称（如医院、研究机构等）
	Role      string         `gorm:"type:varchar(20);not null" json:"role"`    // 用户角色（如医生、研究员、管理员等）
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}