package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"` // 密码不会在JSON中返回
	Name       string    `json:"name"`
	Role       string    `json:"role"` // 角色：doctor, researcher, admin等
	Hospital   string    `json:"hospital,omitempty"`
	Department string    `json:"department,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// UserLogin 用户登录请求
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegister 用户注册请求
type UserRegister struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Hospital   string `json:"hospital,omitempty"`
	Department string `json:"department,omitempty"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Hospital   string    `json:"hospital,omitempty"`
	Department string    `json:"department,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
