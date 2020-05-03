package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var err error

func InitDB(connection string, enableLogging bool) {

	DB, err = gorm.Open("mysql", connection)
	if err != nil {
		panic("failed to connect database")
	}

	DB.LogMode(enableLogging)

	DB.AutoMigrate(&Wallet{})

}
