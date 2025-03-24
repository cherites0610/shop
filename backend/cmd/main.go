package main

import (
	"bots/shop/handler"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 從 .env 讀取變數
	token := os.Getenv("LINE_ACCESS_TOKEN")
	clientID := os.Getenv("LINE_client_id")
	clientSecret := os.Getenv("LINE_client_secret")

	// 確保變數不為空
	if token == "" || clientID == "" || clientSecret == "" {
		log.Fatal("❌ Missing LINE_client_secret or LINE_client_id or LINE_ACCESS_TOKEN in .env")
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                // 允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},     // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // 允许的头
		AllowCredentials: true,
	}))
	router.GET("/commodities", handler.GetCommoditiesHandler)
	router.GET("/commodities/:id", handler.GetCommoditieyByIDHandler)
	router.GET("/lineLogin", handler.LineLoginURLHandler)
	router.POST("/LineAcess", handler.LineAuthHandler)
	router.POST("/sendMessage", handler.SendMessage)

	router.Run(":8081")

}
