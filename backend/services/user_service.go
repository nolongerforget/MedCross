package services

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"medcross/models"
	"medcross/utils"
)

// UserService 用户服务
type UserService struct {
	// 在实际应用中，这里应该有数据库连接
	// 为了演示，我们使用内存存储
	users         map[string]*models.User
	usernameIndex map[string]string // username -> id 映射
}

// NewUserService 创建新的用户服务
func NewUserService() *UserService {
	service := &UserService{
		users:         make(map[string]*models.User),
		usernameIndex: make(map[string]string),
	}

	// 添加一个测试用户
	testUserID := uuid.New().String()
	hashedPassword, _ := utils.HashPassword("password123")
	service.users[testUserID] = &models.User{
		ID:         testUserID,
		Username:   "testuser",
		Password:   hashedPassword,
		Name:       "测试用户",
		Role:       "doctor",
		Hospital:   "协和医院",
		Department: "内科",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	service.usernameIndex["testuser"] = testUserID

	return service
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(userData models.UserRegister) (string, error) {
	// 检查用户名是否已存在
	if s.UsernameExists(userData.Username) {
		return "", errors.New("用户名已存在")
	}

	// 生成唯一ID
	userID := uuid.New().String()

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		log.Printf("密码哈希失败: %v", err)
		return "", err
	}

	// 创建用户
	user := &models.User{
		ID:         userID,
		Username:   userData.Username,
		Password:   hashedPassword,
		Name:       userData.Name,
		Role:       userData.Role,
		Hospital:   userData.Hospital,
		Department: userData.Department,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存用户
	s.users[userID] = user
	s.usernameIndex[userData.Username] = userID

	return userID, nil
}

// VerifyUser 验证用户凭据
func (s *UserService) VerifyUser(username, password string) (*models.User, error) {
	// 查找用户
	userID, exists := s.usernameIndex[username]
	if !exists {
		return nil, errors.New("用户不存在")
	}

	user := s.users[userID]

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("用户不存在")
	}

	return user, nil
}

// UsernameExists 检查用户名是否已存在
func (s *UserService) UsernameExists(username string) bool {
	_, exists := s.usernameIndex[username]
	return exists
}
