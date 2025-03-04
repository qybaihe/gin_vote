package database

import (
	"log"
	"vote-demo/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite3", "vote.db")
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 启用日志
	DB.LogMode(true)

	// 自动迁移数据库结构
	autoMigrate()
}

// 自动迁移数据库结构
func autoMigrate() {
	DB.AutoMigrate(&models.Poll{}, &models.Option{}, &models.Vote{}, &models.User{}, &models.Comment{})
	log.Println("数据库迁移完成")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
} 