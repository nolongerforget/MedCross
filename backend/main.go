package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"d.com/MedCross/backend/routes"
	"d.com/MedCross/backend/config"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	log.Println("MedCross后端服务启动在 :8080 端口")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}