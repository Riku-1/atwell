package config

import (
	"github.com/kelseyhightower/envconfig"
)

// DatabaseConfigurations is configurations about db.
type DatabaseConfigurations struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// GetPrdDBConfig returns database configurations.
func GetPrdDBConfig() (DatabaseConfigurations, error) {
	return getDBConfig("atwell_db")
}

// GetTestDBConfig returns database configurations for tests.
func GetTestDBConfig() (DatabaseConfigurations, error) {
	return getDBConfig("atwell_test_db")
}

// getDBConfig return s database configurations from environment variables.
func getDBConfig(prefix string) (DatabaseConfigurations, error) {
	var dc DatabaseConfigurations
	err := envconfig.Process(prefix, &dc)

	return dc, err
}
