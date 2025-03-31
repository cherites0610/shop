package service

import (
	"bots/shop/models"
	"errors"
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

func CreateCommodityService(commodity *models.Commodity, skus []models.CommoditySpecifications) error {
	if err := SaveCommoditySerive(commodity); err != nil {
		return err
	}

	combo, err := GetSpecValueIDCombinations(commodity.CommodityID)
	if err != nil {
		return err
	}

	if len(combo) != len(skus) {
		// 手動回滾
		DeleteCommodityService(commodity.CommodityID)
		return errors.New("sku數量不對")
	}

	index := 0
	for _, array := range combo {
		skus[index].CommodityID = commodity.CommodityID
		skus[index].SpecValue1ID = uint(array[0]) // 第一個值給 SpecValue1ID
		if len(array) > 1 && array[1] != 0 {      // 檢查是否有第二個值且不為 0
			skus[index].SpecValue2ID = new(uint) // 初始化指針
			*skus[index].SpecValue2ID = uint(array[1])
		}
		index++
	}

	commodity.CommoditySpecifications = skus

	if err := SaveCommoditySerive(commodity); err != nil {
		return err
	}

	return nil
}

func PutCommodityService(commodity *models.Commodity) error {
	oldCommodity, err := models.GetCommodityDetail(commodity.CommodityID)
	if err != nil {
		return err
	}

	oldCommodity_t := models.Commodity{CommodityID: commodity.CommodityID}
	err = models.GetCommodity(&oldCommodity_t)
	if err != nil {
		return err
	}
	// 先處理名字

	// 會插入原本
	if err := UpdateCommodityNameService(commodity.CommodityID, commodity.CommodityName); err != nil {
		return err
	}

	// 在處理規格

	// 在新增新的type

	// 先刪除已經不在的typeID
	index := 0
	for _, old_spec_type := range oldCommodity.SpecificationTypes {
		isLive := false
		for _, new_spec_type := range commodity.SpecificationTypes {
			if old_spec_type.SpecTypeID == new_spec_type.SpecTypeId {
				isLive = true
				break
			}
		}

		if !isLive {
			// 刪除已經不在的type
			fmt.Println("delete")
			fmt.Println(old_spec_type.SpecTypeID)
			models.DeleteSpecificationType(old_spec_type.SpecTypeID)
		}

		if isLive {
			// 修改原本就在的typeID
			fmt.Println("Updete")
			fmt.Println(commodity.SpecificationTypes[index])

			// 處理value
			insider_index := 0
			for _, old_spec_value := range old_spec_type.SpecificationValues {
				isLive := false
				for _, new_spec_value := range commodity.SpecificationTypes[index].SpecificationValues {
					if old_spec_value.SpecValueID == new_spec_value.SpecValueId {
						isLive = true
						break
					}
				}

				if !isLive {
					fmt.Println("Delete.Value")
					fmt.Println(old_spec_value.SpecValueID)
					models.DeleteSpecValue(old_spec_value.SpecValueID)
				} else {
					fmt.Println("Update.Value")
					fmt.Println(commodity.SpecificationTypes[index].SpecificationValues[insider_index].SpecValueId)
					models.SaveSpecValue(models.SpecificationValues{SpecValueId: commodity.SpecificationTypes[index].SpecificationValues[insider_index].SpecValueId, SpecTypeId: commodity.SpecificationTypes[index].SpecTypeId, SpecValue: commodity.SpecificationTypes[index].SpecificationValues[insider_index].SpecValue})
					insider_index++
				}
			}

			models.SaveSpecType(&commodity.SpecificationTypes[index])
			index++
		}
	}

	// 新增新的typeID
	for _, new_spec_type := range commodity.SpecificationTypes {
		if new_spec_type.SpecTypeId == 0 {
			fmt.Println("new")
			fmt.Println(new_spec_type)
			models.SaveSpecType(&new_spec_type)
		}
	}

	// 在處理sku
	combo, _ := GetSpecValueIDCombinations(commodity.CommodityID)
	fmt.Println("combo")
	fmt.Println(combo)

	new_skus := make(map[uint]models.CommoditySpecifications)
	for _, new_sku := range commodity.CommoditySpecifications {
		new_skus[new_sku.CommoditySpecificationsID] = new_sku
	}

	fmt.Println("new_skus")
	fmt.Println(new_skus)

	index = 0
	for _, old_sku := range oldCommodity_t.CommoditySpecifications {
		isLive := false
		for _, new_sku := range commodity.CommoditySpecifications {
			if old_sku.CommoditySpecificationsID == new_sku.CommoditySpecificationsID {
				isLive = true
				break
			}
		}

		if !isLive {
			// 刪除
			fmt.Println("DELETE")
			fmt.Println(old_sku.CommoditySpecificationsID)
			models.DeleteSKU(commodity.CommodityID, old_sku.CommoditySpecificationsID)
		} else {
			// 更改
			fmt.Println("UPDATE")
			fmt.Println(new_skus[old_sku.CommoditySpecificationsID])
			sku := new_skus[old_sku.CommoditySpecificationsID]
			models.SaveSKU(&sku)
			index++
		}
	}

	// 新增沒有id的sku
	for _, new_sku := range commodity.CommoditySpecifications {
		if new_sku.CommoditySpecificationsID == 0 {
			new_sku.SpecValue1ID = uint(combo[index][0])
			if len(combo[index]) > 1 {
				new_sku.SpecValue2ID = new(uint)
				*new_sku.SpecValue2ID = uint(combo[index][1])
			}
			fmt.Println("CREATE")
			fmt.Println(new_sku)
			models.SaveSKU(&new_sku)
			index++
		}
	}

	// 先刪除舊的已經不在新的sku

	return nil
}

// 保存商品
func SaveCommoditySerive(commodity *models.Commodity) error {
	if err := models.SaveCommodity(commodity); err != nil {
		return err
	}
	return nil
}

func UpdateCommodityNameService(commodityID uint, commodityName string) error {
	if err := models.SaveCommodity(&models.Commodity{CommodityID: commodityID, CommodityName: commodityName}); err != nil {
		return err
	}
	return nil
}
