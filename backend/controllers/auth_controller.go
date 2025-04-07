package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"medcross/models"
	"medcross/services"
	"medcross/utils"
)

// AuthController 处理认证相关请求
type AuthController struct {
	userService *services.UserService
}

// NewAuthController 创建新的认证控制器
func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// Register 处理用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var registerData models.UserRegister

	// 绑定请求数据
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查用户名是否已存在
	exists := ac.userService.UsernameExists(registerData.Username)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建用户
	userID, err := ac.userService.CreateUser(registerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"userId":  userID,
	})
}

// Login 处理用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var loginData models.UserLogin

	// 绑定请求数据
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证用户凭据
	user, err := ac.userService.VerifyUser(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 构建响应
	userResponse := models.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Name:       user.Name,
		Role:       user.Role,
		Hospital:   user.Hospital,
		Department: user.Department,
		CreatedAt:  user.CreatedAt,
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}

// GetCurrentUser 获取当前用户信息
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取用户信息
	user, err := ac.userService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 构建响应
	userResponse := models.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Name:       user.Name,
		Role:       user.Role,
		Hospital:   user.Hospital,
		Department: user.Department,
		CreatedAt:  user.CreatedAt,
	}

	c.JSON(http.StatusOK, userResponse)
}

// 使用utils包中的GenerateJWT函数

// 使用utils包中的密码哈希函数
