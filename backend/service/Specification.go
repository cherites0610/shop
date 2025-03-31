package service

import (
	"bots/shop/models"
	"errors"
	"fmt"
)

// 獲取所有 SpecValueID 組合
func GetSpecValueIDCombinations(commodityID uint) ([][]int, error) {
	commodity, err := models.GetCommodityDetail(commodityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get commodity detail: %v", err)
	}

	// 提取每個 SpecificationType 的 SpecValueID
	var specValueIDs [][]int
	for _, specType := range commodity.SpecificationTypes {
		var ids []int
		for _, specValue := range specType.SpecificationValues {
			ids = append(ids, int(specValue.SpecValueID)) // 直接使用 int，無需強制轉型
		}
		specValueIDs = append(specValueIDs, ids)
	}

	// 計算所有 SpecValueID 的組合
	combinations := cartesianProduct(specValueIDs)
	return combinations, nil
}

// 計算笛卡爾積的函數
func cartesianProduct(values [][]int) [][]int {
	if len(values) == 0 {
		return [][]int{}
	}

	result := [][]int{{}}
	for _, pool := range values {
		var temp [][]int
		for _, res := range result {
			for _, val := range pool {
				newCombo := append([]int{}, res...) // 複製當前組合
				newCombo = append(newCombo, val)    // 添加新值
				temp = append(temp, newCombo)
			}
		}
		result = temp
	}
	return result
}

// 新增規格
func SaveCommoditySpecTypeService(SpecType *models.SpecificationType) error {
	// 檢查商品是否存在
	commodity, err := models.GetCommodityDetail(SpecType.CommodityID)
	if err != nil {
		return err
	}

	// 如果已有兩個規格，不予新增
	if len(commodity.SpecificationTypes) == 2 {
		return errors.New("超過規格數量")
	}

	// 新增規格
	// 先保存SpecType資訊
	if err := models.SaveSpecType(SpecType); err != nil {
		return err
	}

	return nil

}

// 修改規格
func UpdateCommoditySpecTypeService(SpecType *models.SpecificationType) error {
	// 取得舊有的 SpecificationType 資料
	oldSpecType, err := models.GetSpecificationTypeBySpecTypeID(SpecType.SpecTypeId)
	if err != nil {
		return err // 如果找不到舊資料，返回錯誤
	}

	// 更新 SpecificationType 的基本資訊（例如 SpecTypeName）
	oldSpecType.SpecTypeName = SpecType.SpecTypeName
	if err := models.SaveSpecType(&oldSpecType); err != nil {
		return err
	}

	// 比較舊的 SpecificationValues 和新的 SpecificationValues
	oldValues := oldSpecType.SpecificationValues
	newValues := SpecType.SpecificationValues

	// 情況 1：舊的比新的多
	if len(oldValues) > len(newValues) {
		// 保留並更新與新值對應的部分
		for i := 0; i < len(newValues); i++ {
			oldValues[i].SpecValue = newValues[i].SpecValue
			if err := models.SaveSpecValue(oldValues[i]); err != nil {
				return err
			}
		}
		// 刪除多餘的舊值
		for i := len(newValues); i < len(oldValues); i++ {
			if err := models.DeleteSpecValue(oldValues[i].SpecValueId); err != nil {
				return err
			}
		}
	}

	// 情況 2：新的比舊的多或相等
	if len(oldValues) <= len(newValues) {
		// 更新舊有的值
		for i := 0; i < len(oldValues); i++ {
			oldValues[i].SpecValue = newValues[i].SpecValue
			if err := models.SaveSpecValue(oldValues[i]); err != nil {
				return err
			}
		}
		// 新增額外的值
		for i := len(oldValues); i < len(newValues); i++ {
			newValues[i].SpecTypeId = SpecType.SpecTypeId // 設置外鍵
			if err := models.SaveSpecValue(newValues[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// 刪除規格
func DeleteCommoditySpecTypeService(CommodityID uint, SpecTypeID uint) error {
	// 刪除規格
	if err := models.DeleteSpecificationType(SpecTypeID); err != nil {
		return err
	}

	return nil
}

// 新增SKU
func CreateSKUService(sku *models.CommoditySpecifications) error {
	if err := models.SaveSKU(sku); err != nil {
		return err
	}

	return nil
}

// 更改SKU
func UpdatewSKUService(sku *models.CommoditySpecifications) error {
	if err := models.SaveSKU(sku); err != nil {
		return err
	}

	return nil
}
