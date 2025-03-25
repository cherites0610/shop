package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // MySQL 驅動
)

// 初始化數據庫連接
func SetupDatabase() (*sql.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bots?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 自定義 JSON 結構
type SpecificationResponse struct {
	CommoditySpecID int      `json:"commodity_spec_id"`
	SpecValue       []string `json:"spec_value"`
	Stock           int      `json:"stock"`
	Price           float64  `json:"price"`
}

type CommodityResponse struct {
	ID             int                     `json:"id"`
	Name           string                  `json:"name"`
	Spec           map[string][]string     `json:"spec"`
	Specifications []SpecificationResponse `json:"specifications"`
}

// 查詢所有商品及其規格
func GetAllCommodities(db *sql.DB) ([]CommodityResponse, error) {
	// SQL 查詢
	query := `
		SELECT 
			c.commodity_id,
			c.commodity_name,
			cs.commodity_spec_id,
			cs.stock,
			cs.price,
			sv1.spec_value AS spec_value_1,
			sv2.spec_value AS spec_value_2,
			st1.spec_type_name AS spec_type_name_1,
			st2.spec_type_name AS spec_type_name_2
		FROM commodities c
		LEFT JOIN commodity_specifications cs ON c.commodity_id = cs.commodity_id
		LEFT JOIN specification_values sv1 ON cs.spec_value_1_id = sv1.spec_value_id
		LEFT JOIN specification_values sv2 ON cs.spec_value_2_id = sv2.spec_value_id
		LEFT JOIN specification_types st1 ON sv1.spec_type_id = st1.spec_type_id
		LEFT JOIN specification_types st2 ON sv2.spec_type_id = st2.spec_type_id
		ORDER BY c.commodity_id, cs.commodity_spec_id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 臨時結構儲存查詢結果
	type tempResult struct {
		CommodityID     int
		CommodityName   string
		CommoditySpecID int
		Stock           int
		Price           float64
		SpecValue1      sql.NullString
		SpecValue2      sql.NullString
		SpecTypeName1   sql.NullString
		SpecTypeName2   sql.NullString
	}

	var results []tempResult
	for rows.Next() {
		var r tempResult
		err := rows.Scan(
			&r.CommodityID,
			&r.CommodityName,
			&r.CommoditySpecID,
			&r.Stock,
			&r.Price,
			&r.SpecValue1,
			&r.SpecValue2,
			&r.SpecTypeName1,
			&r.SpecTypeName2,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 封裝為 CommodityResponse
	var response []CommodityResponse
	currentCommodity := CommodityResponse{}
	specs := []SpecificationResponse{}
	specMap := make(map[string]map[string]bool) // 用於收集唯一的 spec 值

	for i, r := range results {
		// 如果是新的 commodity_id，開始新的 CommodityResponse
		if i == 0 || r.CommodityID != currentCommodity.ID {
			if i > 0 { // 添加前一個 commodity 到 response
				currentCommodity.Spec = make(map[string][]string)
				for key, values := range specMap {
					uniqueValues := []string{}
					for val := range values {
						uniqueValues = append(uniqueValues, val)
					}
					currentCommodity.Spec[key] = uniqueValues
				}
				currentCommodity.Specifications = specs
				response = append(response, currentCommodity)
				specs = []SpecificationResponse{}
				specMap = make(map[string]map[string]bool)
			}
			currentCommodity = CommodityResponse{
				ID:   r.CommodityID,
				Name: r.CommodityName,
			}
		}

		// 處理規格值
		var specValues []string
		if r.SpecValue1.Valid {
			specValues = append(specValues, r.SpecValue1.String)
			if r.SpecTypeName1.Valid {
				if specMap[r.SpecTypeName1.String] == nil {
					specMap[r.SpecTypeName1.String] = make(map[string]bool)
				}
				specMap[r.SpecTypeName1.String][r.SpecValue1.String] = true
			}
		}
		if r.SpecValue2.Valid {
			specValues = append(specValues, r.SpecValue2.String)
			if r.SpecTypeName2.Valid {
				if specMap[r.SpecTypeName2.String] == nil {
					specMap[r.SpecTypeName2.String] = make(map[string]bool)
				}
				specMap[r.SpecTypeName2.String][r.SpecValue2.String] = true
			}
		}

		// 添加規格到當前 commodity
		specs = append(specs, SpecificationResponse{
			CommoditySpecID: r.CommoditySpecID,
			SpecValue:       specValues,
			Stock:           r.Stock,
			Price:           r.Price,
		})
	}

	// 添加最後一個 commodity
	if len(specs) > 0 {
		currentCommodity.Spec = make(map[string][]string)
		for key, values := range specMap {
			uniqueValues := []string{}
			for val := range values {
				uniqueValues = append(uniqueValues, val)
			}
			currentCommodity.Spec[key] = uniqueValues
		}
		currentCommodity.Specifications = specs
		response = append(response, currentCommodity)
	}

	return response, nil
}
