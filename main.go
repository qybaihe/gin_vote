package main

import (
	"log"
	"vote-demo/database"
	"vote-demo/routes"
)

func main() {
	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	r.Run(":8080")
} 