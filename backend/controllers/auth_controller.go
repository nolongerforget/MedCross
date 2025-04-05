package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"d.com/MedCross/backend/config"
	"d.com/MedCross/backend/models"
)

// 用户登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户注册请求结构
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FullName  string `json:"fullName" binding:"required"`
	OrgName   string `json:"orgName" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

// 登录响应结构
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
	User      models.User  `json:"user"`
}

// Login 处理用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 在实际应用中，这里应该从数据库查询用户
	// 这里简化处理，使用硬编码的用户信息进行演示
	user := models.User{
		ID:       "1",
		Username: "admin",
		Password: "$2a$10$XOPbrlUPQdwdJUpSrIF6X.LG1dXgFOoCWRqCX48hKaYGtEZyGW.3O", // 加密的 "password"
		Email:    "admin@example.com",
		FullName: "管理员",
		OrgName:  "医院",
		Role:     "admin",
	}

	// 检查用户名
	if req.Username != user.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, expiresAt, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回登录成功响应
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// Register 处理用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 在实际应用中，这里应该检查用户名是否已存在
	// 并将用户信息保存到数据库
	// 这里简化处理，直接返回成功

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	newUser := models.User{
		ID:       "2", // 在实际应用中，这应该是自动生成的
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		FullName: req.FullName,
		OrgName:  req.OrgName,
		Role:     req.Role,
	}

	// 生成JWT令牌
	token, expiresAt, err := generateToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回注册成功响应
	c.JSON(http.StatusCreated, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      newUser,
	})
}

// 生成JWT令牌
func generateToken(user models.User) (string, int64, error) {
	cfg := config.GetConfig()

	// 设置过期时间
	expiresAt := time.Now().Add(time.Duration(cfg.JWT.ExpiresIn) * time.Second).Unix()

	// 创建JWT声明
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiresAt,
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}