package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Secret string `envconfig:"secret"`
}

// GetYahooAuthConfig returns yahoo japan auth configurations from environment variables.
func GetAppConfig() (AppConfig, error) {
	var c AppConfig
	err := envconfig.Process("atwell_app", &c)

	if c.Secret == "" {
		return AppConfig{}, fmt.Errorf("failed to get auth configurations. Please set them in environments variables. It is now %v", c)
	}

	return c, err
}
