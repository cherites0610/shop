package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL 驅動
)

// 取得所有商品
func GetAllCommodity() ([]CommodityListResponse, error) {
	var commodities []Commodity
	// 資料庫取得
	if err := db.Preload("CommoditySpecifications").Find(&commodities).Error; err != nil {
		return nil, err
	}

	var response []CommodityListResponse
	for _, commodity := range commodities {
		var minPirce, maxPrice float64
		var totalStock uint
		var pictureURL string

		// 遍歷該商品所有規格，找到最大最小價格及total庫存
		if len(commodity.CommoditySpecifications) > 0 {
			minPirce = commodity.CommoditySpecifications[0].Price
			maxPrice = minPirce
			pictureURL = "" // 預設為空字串

			for _, spec := range commodity.CommoditySpecifications {
				if spec.Price < minPirce {
					minPirce = spec.Price
				}
				if spec.Price > maxPrice {
					maxPrice = spec.Price
				}
				totalStock += spec.Stock

				if pictureURL == "" && spec.PictureUrl != nil {
					pictureURL = *spec.PictureUrl
				}
			}
		}

		// 轉換成API輸出格式
		response = append(response, CommodityListResponse{
			CommodityID:   commodity.CommodityID,
			CommodityName: commodity.CommodityName,
			PriceRange: struct {
				Min float64 "json:\"min\""
				Max float64 "json:\"max\""
			}{Min: minPirce, Max: maxPrice},
			TotalStock:          totalStock,
			SpecificationsCount: len(commodity.CommoditySpecifications),
			PictureURL:          pictureURL, // 加入圖片URL
		})
	}

	return response, nil
}

// 取得單個商品
func GetCommodityDetail(id uint) (CommodityDetailResponse, error) {
	var commodity Commodity
	if err := db.Preload("SpecificationTypes.SpecificationValues").
		Preload("CommoditySpecifications.SpecValue1").
		Preload("CommoditySpecifications.SpecValue2").
		First(&commodity, id).Error; err != nil {
		return CommodityDetailResponse{}, err
	}

	specTypes := make([]SpecTypeResponse, len(commodity.SpecificationTypes))
	for i, st := range commodity.SpecificationTypes {
		values := make([]SpecValueResponse, len(st.SpecificationValues))
		for j, sv := range st.SpecificationValues {
			values[j] = SpecValueResponse{SpecValueID: sv.SpecValueId, SpecValue: sv.SpecValue}
		}
		specTypes[i] = SpecTypeResponse{
			SpecTypeID:          st.SpecTypeId,
			SpecTypeName:        st.SpecTypeName,
			SpecificationValues: values,
		}
	}

	commoditySpecs := make([]CommoditySpecResponse, len(commodity.CommoditySpecifications))
	for i, cs := range commodity.CommoditySpecifications {
		var specValue2 *string
		if cs.SpecValue2ID != nil {
			sv2 := cs.SpecValue2.SpecValue
			specValue2 = &sv2
		}

		pictureUrl := ""
		if cs.PictureUrl != nil {
			pictureUrl = *cs.PictureUrl
		}

		commoditySpecs[i] = CommoditySpecResponse{
			CommoditySpecID: cs.CommoditySpecificationsID,
			SpecValue1ID:    cs.SpecValue1ID,
			SpecValue1:      cs.SpecValue1.SpecValue,
			SpecValue2ID:    cs.SpecValue2ID,
			SpecValue2:      specValue2,
			Stock:           cs.Stock,
			Price:           cs.Price,
			PictureURL:      pictureUrl,
		}
	}

	response := CommodityDetailResponse{
		CommodityID:        commodity.CommodityID,
		CommodityName:      commodity.CommodityName,
		SpecificationTypes: specTypes,
		CommoditySpecs:     commoditySpecs,
	}

	return response, nil
}

// 保存商品
func SaveCommodity(CommodityName string, CommodityID *uint) (Commodity, error) {
	// 創建商品
	commodity := Commodity{
		CommodityName: CommodityName,
	}
	if CommodityID != nil {
		commodity.CommodityID = *CommodityID
	}

	if err := db.Save(&commodity).Error; err != nil {
		return Commodity{}, fmt.Errorf("failed to create commodity: %v", err)
	}

	return commodity, nil
}

// 刪除商品
func DeleteCommodity(id uint) error {
	if err := db.Delete(&Commodity{}, id).Error; err != nil {
		return err
	}
	return nil
}
