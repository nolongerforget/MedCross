package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"medcross/controllers"
	"medcross/middleware"
	"medcross/services"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env文件未找到，将使用系统环境变量")
	}

	// 设置运行模式
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	configureCors(r)

	// 初始化服务
	userService := services.NewUserService()
	dataService := services.NewDataService()
	gatewayService := services.NewGatewayService()

	// 初始化控制器
	authController := controllers.NewAuthController(userService)
	dataController := controllers.NewDataController(dataService, gatewayService)

	// 注册路由
	setupRoutes(r, authController, dataController)

	// 获取端口
	port := getEnv("PORT", "8000")

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}

// 配置CORS中间件
func configureCors(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{getEnv("CORS_ALLOW_ORIGINS", "*")}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
}

// 设置路由
func setupRoutes(r *gin.Engine, authController *controllers.AuthController, dataController *controllers.DataController) {
	// API版本组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "服务正常运行",
			})
		})

		// 注册认证路由
		setupAuthRoutes(api, authController)

		// 注册数据路由
		setupDataRoutes(api, dataController)
	}
}

// 设置认证相关路由
func setupAuthRoutes(rg *gin.RouterGroup, authController *controllers.AuthController) {
	auth := rg.Group("/")
	{
		// 登录
		auth.POST("/login", authController.Login)

		// 注册
		auth.POST("/register", authController.Register)

		// 获取用户信息（需要认证）
		auth.GET("/user", middleware.AuthMiddleware(), authController.GetCurrentUser)
	}
}

// 设置数据相关路由
func setupDataRoutes(rg *gin.RouterGroup, dataController *controllers.DataController) {
	data := rg.Group("/")
	{
		// 数据查询
		data.GET("/query", dataController.QueryData)

		// 数据上传（需要认证）
		data.POST("/upload", middleware.AuthMiddleware(), dataController.UploadData)

		// 获取数据类型列表
		data.GET("/data-types", dataController.GetDataTypes)

		// 获取数据详情
		data.GET("/data/:id", dataController.GetDataDetail)

		// 获取统计数据
		data.GET("/statistics", dataController.GetStatistics)
	}
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
