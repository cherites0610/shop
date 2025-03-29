package models

import (
	_ "github.com/go-sql-driver/mysql" // MySQL 驅動
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db, _ = SetupDatabase()
}

func SetupDatabase() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bots?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()
	return db, err
}

type Commodity struct {
	CommodityID             uint                      `gorm:"column:commodity_id;primaryKey;autoIncrement"`
	CommodityName           string                    `gorm:"column:commodity_name"`
	SpecificationTypes      []SpecificationType       `gorm:"foreignKey:CommodityID;references:CommodityID"`
	CommoditySpecifications []CommoditySpecifications `gorm:"foreignKey:CommodityID;references:CommodityID"` // 添加關聯
}

type SpecificationType struct {
	SpecTypeId          uint                  `gorm:"column:spec_type_id;primaryKey;autoIncrement"`
	CommodityID         uint                  `gorm:"column:commodity_id"`
	SpecTypeName        string                `gorm:"column:spec_type_name;size:50"`
	Commodity           Commodity             `gorm:"foreignKey:CommodityID;references:CommodityID"`
	SpecificationValues []SpecificationValues `gorm:"foreignKey:SpecTypeId;references:SpecTypeId"`
}

type SpecificationValues struct {
	SpecValueId       uint              `gorm:"column:spec_value_id;primaryKey;autoIncrement"`
	SpecTypeId        uint              `gorm:"column:spec_type_id"`
	SpecificationType SpecificationType `gorm:"foreignKey:SpecTypeId;references:SpecTypeId"`
	SpecValue         string            `gorm:"column:spec_value"`
}

type CommoditySpecifications struct {
	CommoditySpecificationsID uint                `gorm:"column:commodity_spec_id;primaryKey;autoIncrement"`
	CommodityID               uint                `gorm:"column:commodity_id"`
	Commodity                 Commodity           `gorm:"foreignKey:CommodityID;references:CommodityID"`
	SpecValue1ID              uint                `gorm:"column:spec_value_1_id;"`
	SpecValue1                SpecificationValues `gorm:"foreignKey:SpecValue1ID;references:SpecValueId"`
	SpecValue2ID              *uint               `gorm:"column:spec_value_2_id;"`
	SpecValue2                SpecificationValues `gorm:"foreignKey:SpecValue2ID;references:SpecValueId"`
	Stock                     uint
	Price                     float64
	PictureUrl                *string
}

// 以上皆為實體類定義

// 定義getAllCommodityAPI輸出格式
type CommodityListResponse struct {
	CommodityID   uint   `json:"commodity_id"`
	CommodityName string `json:"commodity_name"`
	PriceRange    struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"price_range"`
	TotalStock          uint   `json:"total_stock"`
	SpecificationsCount int    `json:"specifications_count"`
	PictureURL          string `json:"picture_url"`
}

type CommodityDetailResponse struct {
	CommodityID        uint                    `json:"commodity_id"`
	CommodityName      string                  `json:"commodity_name"`
	SpecificationTypes []SpecTypeResponse      `json:"specification_types"`
	CommoditySpecs     []CommoditySpecResponse `json:"commodity_specifications"`
}

type SpecTypeResponse struct {
	SpecTypeID          uint                `json:"spec_type_id"`
	SpecTypeName        string              `json:"spec_type_name"`
	SpecificationValues []SpecValueResponse `json:"specification_values"`
}

type SpecValueResponse struct {
	SpecValueID uint   `json:"spec_value_id"`
	SpecValue   string `json:"spec_value"`
}

type CommoditySpecResponse struct {
	CommoditySpecID uint    `json:"commodity_spec_id"`
	SpecValue1ID    uint    `json:"spec_value_1_id"`
	SpecValue1      string  `json:"spec_value_1"`
	SpecValue2ID    *uint   `json:"spec_value_2_id,omitempty"`
	SpecValue2      *string `json:"spec_value_2,omitempty"`
	Stock           uint    `json:"stock"`
	Price           float64 `json:"price"`
	PictureURL      string  `json:"picture_url"`
}

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
	// fmt.Println(commodity.SpecificationTypes)
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
