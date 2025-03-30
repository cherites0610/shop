package models

type Commodity struct {
	CommodityID             uint                      `gorm:"column:commodity_id;primaryKey;autoIncrement"`
	CommodityName           string                    `gorm:"column:commodity_name"`
	SpecificationTypes      []SpecificationType       `gorm:"foreignKey:CommodityID;references:CommodityID"`
	CommoditySpecifications []CommoditySpecifications `gorm:"foreignKey:CommodityID;references:CommodityID"`
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

type UpdateCommoditySpecRequest struct {
	Stock      *uint    `json:"stock,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	PictureURL *string  `json:"picture_url,omitempty"`
}

type CreateSKURequest struct {
	CommodityID               uint    `json:"commodity_id" binding:"required"`
	CommoditySpecificationsID *uint   `json:"commodity_spec_id"`
	SpecValue1ID              uint    `json:"spec_value_1_id"`
	SpecValue2ID              *uint   `json:"spec_value_2_id,omitempty"`
	Stock                     uint    `json:"stock" binding:"required"`
	Price                     float64 `json:"price" binding:"required"`
	PictureURL                string  `json:"picture_url" binding:"required"`
}

type CreateCommodityRequest struct {
	CommodityName      string                     `json:"commodity_name"`
	SpecificationTypes []CreateSpecTypeRequest    `json:"specification_types,omitempty"`
	SKU                []CreateSpecTypeSKURequest `json:"sku"`
}

type CreateSpecTypeRequest struct {
	SpecTypeID   *uint                      `json:"spec_type_id"`
	SpecTypeName string                     `json:"spec_type_name" binding:"required"`
	SpecValues   []string                   `json:"spec_type_values" binding:"required"`
	SKU          []CreateSpecTypeSKURequest `json:"sku"`
}

type CreateSpecTypeSKURequest struct {
	CommodityID uint    `json:"commodity_id" binding:"required"`
	Stock       uint    `json:"stock" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	PictureURL  string  `json:"picture_url" binding:"required"`
}

type BuyRequest struct {
	CommodityID uint   `json:"commodity_id" binding:"required"`
	SpecTypeID  uint   `json:"spec_type_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Num         uint   `json:"num" binding:"required"`
}
