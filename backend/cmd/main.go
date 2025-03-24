package main

import (
	"bots/shop/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 初始化数据库
	db, err := models.SetupDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// 获取商品列表
	router.GET("/commodities", func(c *gin.Context) {
		var commodities []models.Commodity
		db.Find(&commodities)
		c.JSON(http.StatusOK, commodities)
	})

	// 获取单个商品
	router.GET("/commodities/:id", func(c *gin.Context) {
		id := c.Param("id")
		var commodity models.Commodity
		if err := db.Where("id = ?", id).First(&commodity).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Commodity not found"})
			return
		}
		c.JSON(http.StatusOK, commodity)
	})

	// 创建新商品
	router.POST("/commodities", func(c *gin.Context) {
		var commodity models.Commodity
		if err := c.ShouldBindJSON(&commodity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&commodity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, commodity)
	})

	// 更新商品
	router.PUT("/commodities/:id", func(c *gin.Context) {
		id := c.Param("id")
		var commodity models.Commodity
		if err := db.Where("id = ?", id).First(&commodity).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Commodity not found"})
			return
		}

		if err := c.ShouldBindJSON(&commodity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&commodity)
		c.JSON(http.StatusOK, commodity)
	})

	// 删除商品
	router.DELETE("/commodities/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Where("id = ?", id).Delete(&models.Commodity{}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Commodity not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Commodity deleted"})
	})

	// 启动服务器
	router.Run(":8081")
}
