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
func SaveCommoditySpecTypeService(CommodityID uint, SpecName string, SpecValue []string) error {
	// 檢查商品是否存在
	commodity, err := models.GetCommodityDetail(CommodityID)
	if err != nil {
		return err
	}

	// 如果已有兩個規格，不予新增
	if len(commodity.SpecificationTypes) == 2 {
		return errors.New("超過規格數量")
	}

	// 新增規格
	if err := models.CreateSpecificationType(CommodityID, nil, SpecName, SpecValue); err != nil {
		return err
	}

	return nil

}

// 修改規格
func UpdateCommoditySpecTypeService(CommodityID uint, req []models.CreateSpecTypeRequest) error {
	// 取得舊有的商品資料
	commodity, err := models.GetCommodityDetail(CommodityID)
	if err != nil {
		return err
	}

	// 取得舊有每個的ValueID
	for _, reqSpecType := range req {
		if reqSpecType.SpecTypeID != nil {
			// 修改

			// 尋找舊有valueID
			oldSpecValueIDs := []uint{}
			for _, commoditySpecType := range commodity.SpecificationTypes {
				if *reqSpecType.SpecTypeID == commoditySpecType.SpecTypeID {
					for _, item := range commoditySpecType.SpecificationValues {
						oldSpecValueIDs = append(oldSpecValueIDs, item.SpecValueID)
					}
				}
			}
			fmt.Println("oldSpecValueIDs", oldSpecValueIDs)
			// 更改每個SpecValue
			index := 0
			for _, specValue := range reqSpecType.SpecValues {
				fmt.Println(oldSpecValueIDs[index])
				models.UpdateSpecValue(oldSpecValueIDs[index], *reqSpecType.SpecTypeID, specValue)
				index++
			}

			if len(oldSpecValueIDs) > len(reqSpecType.SpecValues) {
				models.DeleteSpecValue(oldSpecValueIDs[1])
			}

			// 更改SpecTypeName
			models.SaveSpecificationType(commodity.CommodityID, *reqSpecType.SpecTypeID, reqSpecType.SpecTypeName)
		} else {
			// 新增

			// 逐步替換舊規格
			currentSpecs := len(commodity.SpecificationTypes)
			if currentSpecs == 2 {
				// 先刪除一個舊規格
				for _, commoditySpecType := range commodity.SpecificationTypes {
					isLive := false
					for _, reqSpecTypeInner := range req {
						if reqSpecTypeInner.SpecTypeID != nil && *reqSpecTypeInner.SpecTypeID == commoditySpecType.SpecTypeID {
							isLive = true
							break
						}
					}
					if !isLive {
						fmt.Println("Deleting SpecTypeID (step 1):", commoditySpecType.SpecTypeID)
						if err := DeleteCommoditySpecTypeService(CommodityID, commoditySpecType.SpecTypeID); err != nil {
							return err
						}
						break // 只刪除一個
					}
				}
			}

			// 新增新規格
			if err := SaveCommoditySpecTypeService(CommodityID, reqSpecType.SpecTypeName, reqSpecType.SpecValues); err != nil {
				return err
			}

			// 刪除剩餘的舊規格（如果有）
			for _, commoditySpecType := range commodity.SpecificationTypes {
				isLive := false
				for _, reqSpecTypeInner := range req {
					if reqSpecTypeInner.SpecTypeID != nil && *reqSpecTypeInner.SpecTypeID == commoditySpecType.SpecTypeID {
						isLive = true
						break
					}
				}
				if !isLive {
					fmt.Println("Deleting SpecTypeID (step 2):", commoditySpecType.SpecTypeID)
					if err := DeleteCommoditySpecTypeService(CommodityID, commoditySpecType.SpecTypeID); err != nil {
						fmt.Println("Delete failed (ignored):", err) // 可選擇忽略，因為新增已成功
					}
				}
			}
		}
	}
	return nil
}

// 刪除規格（一定要保留一個）
func DeleteCommoditySpecTypeService(CommodityID uint, SpecTypeID uint) error {
	// 檢查商品是否存在
	commodity, err := models.GetCommodityDetail(CommodityID)
	if err != nil {
		return err
	}

	// 如果已有兩個規格，不予新增
	if len(commodity.SpecificationTypes) == 1 {
		return errors.New("規格數量不能為0")
	}

	// 刪除規格
	if err := models.DeleteSpecificationType(SpecTypeID); err != nil {
		return err
	}

	return nil
}

func CreateSKUAUTOServie(CommodityID uint, SKUS []models.CreateSpecTypeSKURequest) error {
	// 取得所有組合確保sku相同數量
	combinations, err := GetSpecValueIDCombinations(CommodityID)
	if err != nil {
		return err
	}

	// 新增sku
	for i, combo := range combinations {
		specValue1ID := combo[0]
		var specValue2ID *uint
		if len(combo) > 1 {
			value := uint(combo[1])
			specValue2ID = &value
		}

		CreateSKUService(CommodityID, uint(specValue1ID), specValue2ID, SKUS[i].Stock, SKUS[i].Price, SKUS[i].PictureURL)
	}
	return nil
}

// 新增SKU
func CreateSKUService(CommodityID, SpecValue1ID uint, SpecValue2ID *uint, Stock uint, Price float64, PictureURL string) error {
	SKU := models.CommoditySpecifications{
		CommodityID:  CommodityID,
		SpecValue1ID: SpecValue1ID,
		SpecValue2ID: SpecValue2ID,
		Stock:        Stock,
		Price:        Price,
		PictureUrl:   &PictureURL,
	}

	if err := models.CreateSKU(&SKU); err != nil {
		return err
	}

	return nil
}

// 更改SKU
func UpdatewSKUSService(commodityID uint, skus []models.CreateSKURequest) error {
	// 先找資料庫中舊的
	oldSKUs, _ := models.GetCommoditySpecByCommodityID(commodityID)

	// 進行對比，不存在於新的的全部刪除
	for _, oldSKU := range oldSKUs {
		isLive := false
		for _, newSKU := range skus {
			if oldSKU.CommoditySpecificationsID == *newSKU.CommoditySpecificationsID {
				isLive = true
				break
			}
		}

		if !isLive {
			fmt.Printf("已經刪除sku ID:%d", oldSKU.CommoditySpecificationsID)
			models.DeleteSKU(commodityID, oldSKU.CommoditySpecificationsID)
		}

	}
	// 遍歷skus
	for _, sku := range skus {
		if sku.CommoditySpecificationsID != nil {
			// 有id的話是更改
			fmt.Print("UpdateSKU:")
			fmt.Println(sku)
		} else {
			// 沒有id的話是新建
			fmt.Print("NewSKU:")
			fmt.Println(sku)
		}
	}

	//

	return nil
}

// 更改SKU
func UpdatewSKUService(CommodityID, SpecValue1ID, CommoditySpecID uint, SpecValue2ID *uint, Stock uint, Price float64, PictureURL string) error {
	SKU := models.CommoditySpecifications{
		CommoditySpecificationsID: CommoditySpecID,
		CommodityID:               CommodityID,
		SpecValue1ID:              SpecValue1ID,
		SpecValue2ID:              SpecValue2ID,
		Stock:                     Stock,
		Price:                     Price,
		PictureUrl:                &PictureURL,
	}

	// 如果兩個value其中一個 其中依舊存在 先刪除

	if err := models.SaveSKU(&SKU); err != nil {
		return err
	}

	return nil
}
