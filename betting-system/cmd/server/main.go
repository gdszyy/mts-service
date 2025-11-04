package main

import (
	"log"
	"os"

	"github.com/gdsZyy/betting-system/internal/database"
	"github.com/gdsZyy/betting-system/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Betting System...")

	// 加载数据库配置
	dbConfig := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "betting_system"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// 连接数据库
	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// 自动迁移
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated")

	// 创建索引
	if err := database.CreateIndexes(); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	log.Println("Indexes created")

	// 创建 Gin 路由
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

