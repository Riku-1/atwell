package db

import (
	"atwell/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CreateGormDB creates gorm.DB struct from configurations.
func CreateGormDB(dc *config.DatabaseConfigurations) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
