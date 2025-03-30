package handler

import (
	"bots/shop/models"
	"bots/shop/service"
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
	id := c.Param("commodity_id")
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

// 購買Handler
func BuyHandler(c *gin.Context) {
	var req models.BuyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.Buy(req.UserID, req.CommodityID, req.SpecTypeID, req.Num)
	c.JSON(http.StatusOK, req)
}

// 創建商品Handler
func CreateCommodityHandler(c *gin.Context) {
	var req models.CreateCommodityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := service.SaveCommoditySerive(req.CommodityName, nil, req.SpecificationTypes, req.SKU)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "")
}

// 更改商品Handler
func UpdateCommodityHandler(c *gin.Context) {
	type UpdateCommodityRequest struct {
		CommodityName string `json:"commodity_name"`
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	var req UpdateCommodityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	commodityIDPtr := uint(commodityID)
	err = service.SaveCommoditySerive(req.CommodityName, &commodityIDPtr, []models.CreateSpecTypeRequest{}, []models.CreateSpecTypeSKURequest{})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "")
}

// 刪除商品Handler
func DeleteCommodityHandler(c *gin.Context) {
	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	if err := service.DeleteCommodityService(uint(commodityID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, "")
}
