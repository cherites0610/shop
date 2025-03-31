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
	var req models.SpecificationTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CommodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	specValue := []models.SpecificationValues{}

	for _, value := range req.SpecTypeValue {
		temp := models.SpecificationValues{
			SpecValue: value,
		}

		specValue = append(specValue, temp)
	}

	temp := models.SpecificationType{
		CommodityID:         uint(CommodityID),
		SpecTypeName:        req.SpecTypeName,
		SpecificationValues: specValue,
	}

	if err := service.SaveCommoditySpecTypeService(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, temp)
}

// 修改商品規格Handler
func UpdateCommoditySpecTypeHandler(c *gin.Context) {
	var req models.SpecificationTypeRequest
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_type_id"})
		return
	}

	specValue := []models.SpecificationValues{}

	for _, value := range req.SpecTypeValue {
		temp := models.SpecificationValues{
			SpecTypeId: uint(SpecTypeId),
			SpecValue:  value,
		}

		specValue = append(specValue, temp)
	}

	temp := models.SpecificationType{
		SpecTypeId:          uint(SpecTypeId),
		CommodityID:         uint(CommodityID),
		SpecTypeName:        req.SpecTypeName,
		SpecificationValues: specValue,
	}

	err = service.UpdateCommoditySpecTypeService(&temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, temp)
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
	var req models.SKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	sku := models.CommoditySpecifications{
		CommodityID:  uint(commodityID),
		SpecValue1ID: req.SpecValue1ID,
		SpecValue2ID: req.SpecValue2ID,
		Stock:        req.Stock,
		Price:        req.Price,
		PictureUrl:   &req.PictureURL,
	}

	if err := models.SaveSKU(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, sku)
}

// 更新單個SKU Handler
func UpdateSKUHandler(c *gin.Context) {
	var req models.SKURequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sku_id, err := strconv.ParseUint(c.Param("sku_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_type_id"})
		return
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	sku := models.CommoditySpecifications{
		CommoditySpecificationsID: uint(sku_id),
		CommodityID:               uint(commodityID),
		SpecValue1ID:              req.SpecValue1ID,
		SpecValue2ID:              req.SpecValue2ID,
		Stock:                     req.Stock,
		Price:                     req.Price,
		PictureUrl:                &req.PictureURL,
	}

	if err := service.UpdatewSKUService(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, sku)
}
