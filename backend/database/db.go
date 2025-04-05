package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/medcross/backend/config"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	log.Println("数据库连接成功")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("获取数据库连接失败: %v\n", err)
			return
		}
		sqlDB.Close()
		log.Println("数据库连接已关闭")
	}
}

// AutoMigrateModels 自动迁移数据库模型
func AutoMigrateModels(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}