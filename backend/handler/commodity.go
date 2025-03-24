package handler

import (
	"bots/shop/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	// 初始化数据库
	db, err = models.SetupDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
}

func GetCommoditiesHandler(c *gin.Context) {
	var commodities []models.Commodity
	db.Find(&commodities)
	c.JSON(http.StatusOK, commodities)
}

func GetCommoditieyByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var commodity models.Commodity
	if err := db.Where("id = ?", id).First(&commodity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Commodity not found"})
		return
	}
	c.JSON(http.StatusOK, commodity)
}
