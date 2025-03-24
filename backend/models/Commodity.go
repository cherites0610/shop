package models

import (
	"encoding/json"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Commodity 商品模型
type Commodity struct {
	ID            uint            `json:"id" gorm:"primary_key"`
	Name          string          `json:"name"`
	Price         float64         `json:"price"`
	Specification json.RawMessage `json:"specification"`
	Stock         int             `json:"stock"`
}

// 初始化数据库连接
func SetupDatabase() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bots?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Commodity{})
	return db, nil
}
