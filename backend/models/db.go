package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var CmsDb *gorm.DB

func InitDB(dataSourceName string) error {
	var err error
	CmsDb, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	return err
}
