package database

import (
	"log"

	"github.com/medcross/backend/models"
)

// MigrateDB 执行数据库迁移，创建表结构
func MigrateDB() {
	log.Println("开始数据库迁移...")

	// 自动迁移数据库模型
	models := []interface{}{
		&models.User{},
		&models.MedicalData{},
		&models.Tag{},
		&models.DataTag{},
		&models.Authorization{},
		&models.AccessRecord{},
		&models.CrossChainReference{},
	}

	if err := DB.AutoMigrate(models...); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")

	// 创建初始数据
	createInitialData()
}

// createInitialData 创建初始数据
func createInitialData() {
	// 检查是否已存在管理员用户
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		// 创建默认管理员用户
		admin := models.User{
			Username: "admin",
			Password: "$2a$10$EqKCjGBKQDGMVLzRr1D.6.9UOLWRZGUYfPZB.wqyeAIKQKlc5I5Vy", // 默认密码: admin123
			Email:    "admin@medcross.org",
			FullName: "System Administrator",
			OrgName:  "MedCross",
			Role:     "admin",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("创建管理员用户失败: %v\n", err)
		} else {
			log.Println("已创建默认管理员用户")
		}
	}

	// 创建默认数据标签
	defaultTags := []string{"电子病历", "影像数据", "检验报告", "处方信息", "基因组数据"}
	for _, tagName := range defaultTags {
		var tagCount int64
		DB.Model(&models.Tag{}).Where("name = ?", tagName).Count(&tagCount)

		if tagCount == 0 {
			tag := models.Tag{Name: tagName}
			if err := DB.Create(&tag).Error; err != nil {
				log.Printf("创建标签 '%s' 失败: %v\n", tagName, err)
			}
		}
	}

	log.Println("初始数据创建完成")
}