package models

import (
	_ "github.com/go-sql-driver/mysql" // MySQL 驅動
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = SetupDatabase()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
}

func SetupDatabase() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bots?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db = db.Debug()
	return db, err
}
