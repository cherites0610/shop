package service

import (
	"bots/shop/models"
	"fmt"
)

// 刪除商品
func DeleteCommodityService(commodity_id uint) error {
	// 先查詢該商品所有訊息
	commodity, err := models.GetCommodityDetail(commodity_id)
	if err != nil {
		return err
	}

	// 若有的話先刪除所有依賴
	for _, spec_type := range commodity.SpecificationTypes {
		if err := models.DeleteSpecificationType(spec_type.SpecTypeID); err != nil {
			return err
		}
	}

	// 刪除其本體
	if err := models.DeleteCommodity(commodity_id); err != nil {
		return err
	}

	return nil
}

// 保存商品
func SaveCommoditySerive(CommodityName string, CommodityID *uint, SpecificationTypes []models.CreateSpecTypeRequest, SKUS []models.CreateSpecTypeSKURequest) error {
	// 先保存商品
	commodity, err := models.SaveCommodity(CommodityName, CommodityID)
	if err != nil {
		return fmt.Errorf("failed to create commodity: %v", err)
	}

	// 若為修改則不保存規格
	if CommodityID != nil {
		return nil
	}

	// 再保存規格
	for _, specTypeReq := range SpecificationTypes {
		err := SaveCommoditySpecTypeService(commodity.CommodityID, specTypeReq.SpecTypeName, specTypeReq.SpecValues)
		if err != nil {
			return fmt.Errorf("failed to create specification type: %v", err)
		}
	}

	// 在保存sku
	CreateSKUAUTOServie(commodity.CommodityID, SKUS)

	return nil
}
