package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查接口
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "MedCross API",
		"version":   "1.0.0",
	})
}