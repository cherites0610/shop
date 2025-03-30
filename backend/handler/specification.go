package handler

import (
	"bots/shop/models"
	"bots/shop/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 創建商品規格Handler
func CreateCommoditySpecTypeHandler(c *gin.Context) {
	var req models.CreateSpecTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	// 從1個規格變兩個，故sku必定是4個
	if len(req.SKU) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sku number"})
		return
	}

	// 處理新增規則
	err = service.SaveCommoditySpecTypeService(uint(commodityID), req.SpecTypeName, req.SpecValues)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	service.CreateSKUAUTOServie(uint(commodityID), req.SKU)

	c.JSON(http.StatusCreated, "")
}

// 修改多個商品規格Handler
func UpdateCommoditySpecTypesHandler(c *gin.Context) {
	var req []models.CreateSpecTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CommodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	err = service.UpdateCommoditySpecTypeService(uint(CommodityID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, "")
}

// 修改商品規格Handler
func UpdateCommoditySpecTypeHandler(c *gin.Context) {
	var req models.CreateSpecTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CommodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	SpecTypeId, err := strconv.ParseUint(c.Param("spec_type_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	specTypeIDValue := uint(SpecTypeId)
	req.SpecTypeID = &specTypeIDValue

	err = service.UpdateCommoditySpecTypeService(uint(CommodityID), []models.CreateSpecTypeRequest{req})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	c.JSON(http.StatusCreated, "")
}

// 刪除商品規格Handler
func DeleteCommoditySpecTypeHandler(c *gin.Context) {
	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	specTypeID, err := strconv.ParseUint(c.Param("spec_type_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_type_id"})
		return
	}

	if err := service.DeleteCommoditySpecTypeService(uint(commodityID), uint(specTypeID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "")
}

// 創建商品SKU Handler
func CreateSKUHandler(c *gin.Context) {
	var sku models.CreateSKURequest
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.CreateSKUService(sku.CommodityID, sku.SpecValue1ID, sku.SpecValue2ID, sku.Stock, sku.Price, sku.PictureURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, sku)
}

func UpdateSKUSHandler(c *gin.Context) {
	var skus []models.CreateSKURequest

	if err := c.ShouldBindJSON(&skus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	if err := service.UpdatewSKUSService(uint(commodityID), skus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, skus)
}

func UpdateSKUHandler(c *gin.Context) {
	var sku models.CreateSKURequest

	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sku_id, err := strconv.ParseUint(c.Param("sku_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_type_id"})
		return
	}

	if err := service.UpdatewSKUService(sku.CommodityID, sku.SpecValue1ID, uint(sku_id), sku.SpecValue2ID, sku.Stock, sku.Price, sku.PictureURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, sku)
}
