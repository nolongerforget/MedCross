package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// 用户结构
type User struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	PasswordHash   string    `json:"passwordHash,omitempty"`
	Email          string    `json:"email"`
	EthereumWallet string    `json:"ethereumWallet"` // 以太坊钱包地址
	FabricID       string    `json:"fabricId"`       // Fabric身份ID
	CreatedAt      time.Time `json:"createdAt"`
}

// JWT声明结构
type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// 内存中存储用户数据（实际应用中应使用数据库）
var users = make(map[string]User)

// JWT密钥（实际应用中应从环境变量或配置文件中读取）
var jwtKey = []byte("your-secret-key")

// 生成随机ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 生成以太坊钱包地址（模拟）
func generateEthereumWallet() string {
	b := make([]byte, 20)
	rand.Read(b)
	return "0x" + hex.EncodeToString(b)
}

// 生成Fabric ID（模拟）
func generateFabricID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "fabric-" + hex.EncodeToString(b)
}

// 用户注册
func register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查用户名是否已存在
	for _, u := range users {
		if u.Username == input.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	userID := generateID()
	newUser := User{
		ID:             userID,
		Username:       input.Username,
		PasswordHash:   string(hashedPassword),
		Email:          input.Email,
		EthereumWallet: generateEthereumWallet(),
		FabricID:       generateFabricID(),
		CreatedAt:      time.Now(),
	}

	// 存储用户
	users[userID] = newUser

	// 创建JWT令牌
	token, err := createToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 返回用户信息（不包含密码哈希）
	newUser.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{
		"user":  newUser,
		"token": token,
	})
}

// 用户登录
func login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 查找用户
	var foundUser User
	var found bool
	for _, u := range users {
		if u.Username == input.Username {
			foundUser = u
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 创建JWT令牌
	token, err := createToken(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 返回用户信息（不包含密码哈希）
	foundUser.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{
		"user":  foundUser,
		"token": token,
	})
}

// 创建JWT令牌
func createToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// 验证JWT令牌
func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// 认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}

		tokenString := authorization[7:] // 移除 "Bearer " 前缀
		claims, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			return
		}

		// 检查用户是否存在
		user, exists := users[claims.UserID]
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user", user)
		c.Next()
	}
}

// 获取当前用户信息
func getCurrentUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// 代理跨链网关的数据查询请求
func proxyQueryData(c *gin.Context) {
	// 从跨链网关获取数据
	// 在实际应用中，这里应该调用跨链网关的API
	// 为了演示，我们直接返回模拟数据

	// 获取查询参数
	keyword := c.Query("keyword")
	dataType := c.Query("dataType")
	chainSource := c.Query("chain")

	log.Printf("代理查询请求: 关键词=%s, 类型=%s, 链源=%s", keyword, dataType, chainSource)

	// 这里应该是调用跨链网关的代码
	// 为了演示，我们返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"totalCount": 4,
		"data": []map[string]interface{}{
			{
				"id":        "eth-001",
				"owner":     "0x1234567890abcdef",
				"dataHash":  "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o",
				"dataType":  "影像数据",
				"metadata":  map[string]string{"patientId": "P12345", "hospital": "协和医院"},
				"timestamp": time.Now().Add(-24 * time.Hour),
				"keywords":  "肺部,CT,影像",
				"chain":     "ethereum",
			},
			{
				"id":        "eth-002",
				"owner":     "0xabcdef1234567890",
				"dataType":  "电子病历",
				"dataHash":  "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn",
				"metadata":  map[string]string{"patientId": "P54321", "hospital": "人民医院"},
				"timestamp": time.Now().Add(-48 * time.Hour),
				"keywords":  "糖尿病,慢性病,病历",
				"chain":     "ethereum",
			},
			{
				"id":        "fab-001",
				"owner":     "user1",
				"dataHash":  "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ",
				"dataType":  "基因组数据",
				"metadata":  map[string]string{"patientId": "P98765", "hospital": "医学研究中心"},
				"timestamp": time.Now().Add(-12 * time.Hour),
				"keywords":  "基因,癌症,研究",
				"chain":     "fabric",
			},
			{
				"id":        "fab-002",
				"owner":     "user2",
				"dataHash":  "QmT8CUmNPMYGe8P9G2XKZHUuWaq9ZqCTGGYVqx57FuLSdT",
				"dataType":  "影像数据",
				"metadata":  map[string]string{"patientId": "P24680", "hospital": "第三医院"},
				"timestamp": time.Now().Add(-36 * time.Hour),
				"keywords":  "脑部,MRI,影像",
				"chain":     "fabric",
			},
		},
	})
}

// 代理跨链网关的数据上传请求
func proxyUploadData(c *gin.Context) {
	// 获取当前用户
	user, _ := c.Get("user")
	currentUser, ok := user.(User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息获取失败"})
		return
	}

	// 解析请求体
	var uploadData struct {
		DataHash string `json:"dataHash" binding:"required"`
		DataType string `json:"dataType" binding:"required"`
		Metadata string `json:"metadata"`
		Keywords string `json:"keywords"`
		Chain    string `json:"chain" binding:"required"` // 目标区块链
	}

	if err := c.ShouldBindJSON(&uploadData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证目标链
	if uploadData.Chain != "ethereum" && uploadData.Chain != "fabric" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标区块链，必须是 'ethereum' 或 'fabric'"})
		return
	}

	// 根据目标链选择用户的区块链身份
	var owner string
	if uploadData.Chain == "ethereum" {
		owner = currentUser.EthereumWallet
	} else {
		owner = currentUser.FabricID
	}

	// 准备发送到跨链网关的数据
	gatewayData := map[string]interface{}{
		"dataHash": uploadData.DataHash,
		"dataType": uploadData.DataType,
		"metadata": uploadData.Metadata,
		"keywords": uploadData.Keywords,
		"chain":    uploadData.Chain,
		"owner":    owner,
	}

	// 这里应该是调用跨链网关的代码
	// 为了演示，我们只记录日志并返回成功
	log.Printf("代理上传请求: %+v", gatewayData)

	// 生成模拟ID
	var id string
	if uploadData.Chain == "ethereum" {
		id = fmt.Sprintf("eth-%d", time.Now().Unix())
	} else {
		id = fmt.Sprintf("fab-%d", time.Now().Unix())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
		"message": fmt.Sprintf("数据已成功上传到%s链", uploadData.Chain),
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

	// 这里应该是调用跨链网关的代码
	// 为了演示，我们返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"owner":     "0x1234567890abcdef",
		"dataHash":  "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o",
		"dataType":  "影像数据",
		"metadata":  map[string]string{"patientId": "P12345", "hospital": "协和医院", "description": "详细描述..."},
		"timestamp": time.Now().Add(-24 * time.Hour),
		"keywords":  "肺部,CT,影像",
		"chain":     "ethereum",
	})
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

	// 公开路由
	r.POST("/api/register", register)           // 用户注册
	r.POST("/api/login", login)                // 用户登录
	r.GET("/api/data-types", getDataTypes)     // 获取数据类型列表

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(authMiddleware())
	{
		auth.GET("/user", getCurrentUser)         // 获取当前用户信息
		auth.GET("/query", proxyQueryData)        // 查询数据
		auth.POST("/upload", proxyUploadData)      // 上传数据
		auth.GET("/data/:id", getDataDetail)      // 获取数据详情
	}

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // 默认端口
	}

	// 启动服务器
	log.Printf("后端API服务启动在 :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}