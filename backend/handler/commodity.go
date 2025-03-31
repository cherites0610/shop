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
	var req models.CommodtiyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	specTypes := make([]models.SpecificationType, 0, len(req.SpecificationTypeRequest))
	for _, item := range req.SpecificationTypeRequest {
		specValues := make([]models.SpecificationValues, 0, len(item.SpecTypeValue))
		for _, value := range item.SpecTypeValue {
			specValues = append(specValues, models.SpecificationValues{SpecValue: value})
		}

		specTypes = append(specTypes, models.SpecificationType{
			SpecTypeName:        item.SpecTypeName,
			SpecificationValues: specValues,
		})
	}

	skus := make([]models.CommoditySpecifications, 0, len(req.SKUTypeRequest))
	for _, item := range req.SKUTypeRequest {
		sku := models.CommoditySpecifications{
			Stock:      item.Stock,
			Price:      item.Price,
			PictureUrl: &item.PictureURL,
		}
		skus = append(skus, sku)
	}

	commodity := models.Commodity{
		CommodityName:      req.CommodityName,
		SpecificationTypes: specTypes,
	}

	if err := service.CreateCommodityService(&commodity, skus); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, commodity)
}

// 更改商品名稱andler
func UpdateCommodityHandler(c *gin.Context) {
	var req models.CommodtiyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	commodityID, err := strconv.ParseUint(c.Param("commodity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commodity_id"})
		return
	}

	temp := models.Commodity{
		CommodityID:   uint(commodityID),
		CommodityName: req.CommodityName,
	}

	if err := models.SaveCommodity(&temp); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, temp)
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

func PutCommodityHandler(c *gin.Context) {
	var req models.CommodityDetailResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SpecificationTypes := []models.SpecificationType{}
	for _, spec_type := range req.SpecificationTypes {
		spec_values := []models.SpecificationValues{}
		for _, spec_value := range spec_type.SpecificationValues {
			spec_values = append(spec_values, models.SpecificationValues{
				SpecValueId: spec_value.SpecValueID,
				SpecValue:   spec_value.SpecValue,
			})
		}

		SpecificationTypes = append(SpecificationTypes, models.SpecificationType{
			CommodityID:         req.CommodityID,
			SpecTypeId:          spec_type.SpecTypeID,
			SpecTypeName:        spec_type.SpecTypeName,
			SpecificationValues: spec_values,
		})
	}

	skus := []models.CommoditySpecifications{}
	for _, req_sku := range req.CommoditySpecs {
		skus = append(skus, models.CommoditySpecifications{
			CommoditySpecificationsID: req_sku.CommoditySpecID,
			CommodityID:               req.CommodityID,
			SpecValue1ID:              req_sku.SpecValue1ID,
			SpecValue2ID:              req_sku.SpecValue2ID,
			Stock:                     req_sku.Stock,
			Price:                     req_sku.Price,
			PictureUrl:                &req_sku.PictureURL,
		})
	}

	commodity := models.Commodity{
		CommodityID:             req.CommodityID,
		CommodityName:           req.CommodityName,
		SpecificationTypes:      SpecificationTypes,
		CommoditySpecifications: skus,
	}

	if err := service.PutCommodityService(&commodity); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, commodity)
}
