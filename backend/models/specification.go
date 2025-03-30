package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 保存規格類型
func CreateSpecificationType(CommodityID uint, SpecTypeID *uint, SpecName string, SpecValue []string) error {
	// 創建規格類型
	specType := SpecificationType{
		CommodityID:  CommodityID,
		SpecTypeName: SpecName,
	}

	if SpecTypeID != nil {
		specType.SpecTypeId = *SpecTypeID
	}

	if err := db.Save(&specType).Error; err != nil {
		return err
	}

	// 創建規格值
	var specValues []SpecificationValues
	for _, value := range SpecValue {
		specValues = append(specValues, SpecificationValues{
			SpecTypeId: specType.SpecTypeId,
			SpecValue:  value,
		})
	}
	if len(specValues) > 0 {
		if err := db.Save(&specValues).Error; err != nil {
			return err
		}
	}

	return nil
}

func SaveSpecificationType(CommodityID uint, SpecTypeID uint, SpecName string) error {
	specType := SpecificationType{
		CommodityID:  CommodityID,
		SpecTypeId:   SpecTypeID,
		SpecTypeName: SpecName,
	}

	if err := db.Save(&specType).Error; err != nil {
		return err
	}

	return nil
}

// 刪除規格類型
func DeleteSpecificationType(id uint) error {
	// 查找規格類型
	var specType SpecificationType
	if err := db.First(&specType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("specification type not found")
		}
		return err
	}

	// 查找所有依賴該規格類型的商品規格組合
	var commoditySpecs []CommoditySpecifications
	if err := db.Where("spec_value_1_id IN (?) OR spec_value_2_id IN (?)",
		db.Model(&SpecificationValues{}).Where("spec_type_id = ?", id).Select("spec_value_id"),
		db.Model(&SpecificationValues{}).Where("spec_type_id = ?", id).Select("spec_value_id")).
		Find(&commoditySpecs).Error; err != nil {
		return fmt.Errorf("failed to query dependent commodity specifications: %v", err)
	}

	// 如果有依賴，逐一刪除
	if len(commoditySpecs) > 0 {
		for _, cs := range commoditySpecs {
			if _, err := DeleteSKU(cs.CommodityID, cs.CommoditySpecificationsID); err != nil {
				return fmt.Errorf("failed to delete commodity specification %d: %v", cs.CommoditySpecificationsID, err)
			}
		}
	}

	// 刪除規格值
	if err := db.Where("spec_type_id = ?", id).Delete(&SpecificationValues{}).Error; err != nil {
		return fmt.Errorf("failed to delete specification values: %v", err)
	}

	// 刪除規格類型
	if err := db.Delete(&specType).Error; err != nil {
		return fmt.Errorf("failed to delete specification type: %v", err)
	}

	return nil
}

// 取得規格
func GetSpecificationTypeBySpecTypeID(SpecificationID uint) (SpecificationType, error) {
	var specificationType SpecificationType
	if err := db.Find(&SpecificationType{}, SpecificationID).Error; err != nil {
		return SpecificationType{}, err
	}
	return specificationType, nil
}

// 更新規格Value值
func UpdateSpecValue(SpecValueID, SpecTypeID uint, SpecValueName string) error {
	if err := db.Save(&SpecificationValues{SpecValueId: SpecValueID, SpecTypeId: SpecTypeID, SpecValue: SpecValueName}).Error; err != nil {
		return err
	}
	return nil
}

// 刪除規格Value值
func DeleteSpecValue(SpecValueID uint) error {
	// 查找所有依賴該規格類型的商品規格組合
	var commoditySpecs []CommoditySpecifications
	if err := db.Where("spec_value_1_id = (?) OR spec_value_2_id = (?)", SpecValueID, SpecValueID).
		Find(&commoditySpecs).Error; err != nil {
		return fmt.Errorf("failed to query dependent commodity specifications: %v", err)
	}

	// 如果有依賴，逐一刪除
	if len(commoditySpecs) > 0 {
		for _, cs := range commoditySpecs {
			if _, err := DeleteSKU(cs.CommodityID, cs.CommoditySpecificationsID); err != nil {
				return fmt.Errorf("failed to delete commodity specification %d: %v", cs.CommoditySpecificationsID, err)
			}
		}
	}

	if err := db.Delete(&SpecificationValues{SpecValueId: SpecValueID}).Error; err != nil {
		return err
	}
	return nil
}

// 取得SKU
func GetCommoditySpec(commoditySpecID uint) CommoditySpecifications {
	var commoditySpec CommoditySpecifications
	if err := db.Preload("SpecValue1").Preload("SpecValue2").Preload("Commodity").First(&commoditySpec, commoditySpecID).Error; err != nil {
		fmt.Println(err)
	}

	return commoditySpec
}

// 創建SKU
func CreateSKU(SKU *CommoditySpecifications) error {
	if err := db.Create(SKU).Error; err != nil {
		return err
	}
	return nil
}

// 創建SKU
func SaveSKU(SKU *CommoditySpecifications) error {
	if err := db.Save(SKU).Error; err != nil {
		return err
	}
	return nil
}

// 更新SKU
func UpdateCommoditySpecification(commodityID, commoditySpecID uint, req UpdateCommoditySpecRequest) (CommoditySpecResponse, error) {
	// 檢查規格組合是否存在並屬於該商品
	var commoditySpec CommoditySpecifications
	if err := db.Where("commodity_id = ? AND commodity_spec_id = ?", commodityID, commoditySpecID).
		Preload("SpecValue1").Preload("SpecValue2").First(&commoditySpec).Error; err != nil {
		return CommoditySpecResponse{}, fmt.Errorf("commodity specification not found: %v", err)
	}

	// 更新字段（僅更新有提供的字段）
	if req.Stock != nil {
		commoditySpec.Stock = *req.Stock
	}
	if req.Price != nil {
		commoditySpec.Price = *req.Price
	}
	if req.PictureURL != nil {
		commoditySpec.PictureUrl = req.PictureURL
	}
	if err := db.Save(&commoditySpec).Error; err != nil {
		return CommoditySpecResponse{}, fmt.Errorf("failed to update commodity specification: %v", err)
	}

	// 準備回傳資料
	var specValue2 *string
	if commoditySpec.SpecValue2ID != nil {
		specValue2 = &commoditySpec.SpecValue2.SpecValue
	}
	pictureURL := ""
	if commoditySpec.PictureUrl != nil {
		pictureURL = *commoditySpec.PictureUrl
	}

	return CommoditySpecResponse{
		CommoditySpecID: commoditySpec.CommoditySpecificationsID,
		SpecValue1ID:    commoditySpec.SpecValue1ID,
		SpecValue1:      commoditySpec.SpecValue1.SpecValue,
		SpecValue2ID:    commoditySpec.SpecValue2ID,
		SpecValue2:      specValue2,
		Stock:           commoditySpec.Stock,
		Price:           commoditySpec.Price,
		PictureURL:      pictureURL,
	}, nil
}

// 刪除SKU
func DeleteSKU(commodityID, SKUID uint) (map[string]string, error) {
	// 檢查規格組合是否存在並屬於該商品（好像有點多餘?）
	var commoditySpec CommoditySpecifications
	if err := db.Where("commodity_id = ? AND commodity_spec_id = ?", commodityID, SKUID).
		First(&commoditySpec).Error; err != nil {
		return nil, fmt.Errorf("commodity specification not found: %v", err)
	}

	// 刪除規格組合
	if err := db.Delete(&commoditySpec).Error; err != nil {
		return nil, fmt.Errorf("failed to delete commodity specification: %v", err)
	}

	return map[string]string{"message": "Commodity specification deleted successfully"}, nil
}

func GetCommoditySpecByCommodityID(commodityID uint) ([]CommoditySpecifications, error) {
	var commoditySpecifications []CommoditySpecifications
	if err := db.Find(&commoditySpecifications).Error; err != nil {
		return []CommoditySpecifications{}, err
	}
	return commoditySpecifications, nil
}
