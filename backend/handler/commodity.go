package handler

import (
	"bots/shop/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCommoditiesHandler(c *gin.Context) {
	var commodities []models.CommodityListResponse
	commodities, _ = models.GetAllCommodity()
	c.JSON(http.StatusOK, commodities)
}

func GetCommoditieyByIDHandler(c *gin.Context) {
	id := c.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	commodity, err := models.GetCommodityDetail(uint(uintID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commodity"})
		return
	}
	c.JSON(http.StatusOK, commodity)
}
