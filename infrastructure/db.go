package infrastructure

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseConfigurations is configurations about db.
type DatabaseConfigurations struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// GetPrdGormDB returns gorm.DB struct for production.
func GetPrdGormDB() (*gorm.DB, error) {
	c := getDBConfig("atwell_db")
	db, err := createGormDB(&c)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetDevGormDB returns gorm.DB struct for production.
func GetDevGormDB() (*gorm.DB, error) {
	c := getDBConfig("atwell_test_db")
	db, err := createGormDB(&c)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getDBConfig return s database configurations from environment variables.
func getDBConfig(prefix string) DatabaseConfigurations {
	var c DatabaseConfigurations
	err := envconfig.Process(prefix, &c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Host == "" || c.Port == 0 || c.User == "" || c.Password == "" || c.DBName == "" {
		log.Fatalf("failed to get database configurations. Please set them in environments variables. It is now %v", c)
	}

	return c
}

// createGormDB creates gorm.DB struct from configurations.
func createGormDB(dc *DatabaseConfigurations) (*gorm.DB, error) {
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
