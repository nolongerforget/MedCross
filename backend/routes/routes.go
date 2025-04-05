package routes

import (
	"github.com/gin-gonic/gin"
	
	"d.com/MedCross/backend/controllers"
	"d.com/MedCross/backend/middleware"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine) {
	// API版本前缀
	api := r.Group("/api/v1")

	// 公共路由
	public := api.Group("/")
	{
		// 健康检查
		public.GET("/health", controllers.HealthCheck)
		
		// 用户认证
		auth := public.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 医疗数据管理
		data := protected.Group("/data")
		{
			// 数据上传
			data.POST("/upload", controllers.UploadData)
			
			// 数据查询
			data.GET("/list", controllers.ListData)
			data.GET("/search", controllers.SearchData)
			data.GET("/:id", controllers.GetDataById)
			
			// 数据授权管理
			data.POST("/:id/authorize", controllers.AuthorizeData)
			data.DELETE("/:id/authorize/:userId", controllers.RevokeAuthorization)
			data.GET("/:id/authorizations", controllers.ListAuthorizations)
			
			// 数据访问记录
			data.GET("/:id/access-logs", controllers.GetAccessLogs)
		}
		
		// 区块链记录
		blockchain := protected.Group("/blockchain")
		{
			// 获取区块链交易记录
			blockchain.GET("/transactions", controllers.GetTransactions)
			blockchain.GET("/transactions/:txHash", controllers.GetTransactionByHash)
		}
		
		// 用户管理
		user := protected.Group("/user")
		{
			// 获取用户信息
			user.GET("/profile", controllers.GetUserProfile)
			// 更新用户信息
			user.PUT("/profile", controllers.UpdateUserProfile)
		}
		
		// 跨链操作
		crosschain := protected.Group("/crosschain")
		{
			// 跨链数据同步
			crosschain.POST("/sync/:id", controllers.SyncDataAcrossChains)
			// 获取跨链状态
			crosschain.GET("/status/:id", controllers.GetCrosschainStatus)
		}
	}
}