package routes

import (
	"bots/shop/handler"

	"github.com/gin-gonic/gin"
)

func CommodityRoutes(router *gin.Engine) {
	// 商品相關路由
	commodityRoutes := router.Group("/commodities")
	{
		commodityRoutes.GET("", handler.GetCommoditiesHandler)
		commodityRoutes.GET("/:commodity_id", handler.GetCommoditieyByIDHandler)
		commodityRoutes.POST("", handler.CreateCommodityHandler)
		commodityRoutes.PUT("", handler.PutCommodityHandler)
		commodityRoutes.PUT("/:commodity_id", handler.UpdateCommodityHandler)
		commodityRoutes.DELETE("/:commodity_id", handler.DeleteCommodityHandler)

		// 購買商品
		commodityRoutes.POST("/buy", handler.BuyHandler)

		// 規格相關路由（規格是商品的子路由）
		specificationRoutes := commodityRoutes.Group("/:commodity_id/specification-types")
		{
			specificationRoutes.POST("", handler.CreateCommoditySpecTypeHandler)                 // 新增規格
			specificationRoutes.PUT("/:spec_type_id", handler.UpdateCommoditySpecTypeHandler)    // 修改規格
			specificationRoutes.DELETE("/:spec_type_id", handler.DeleteCommoditySpecTypeHandler) // 刪除規格
		}

		skuRoutes := commodityRoutes.Group("/:commodity_id/sku")
		{
			skuRoutes.POST("", handler.CreateSKUHandler)
			skuRoutes.PUT("/:sku_id", handler.UpdateSKUHandler)
		}
	}
}
