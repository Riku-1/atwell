package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// YahooAuthConfigurations is configuration about yahoo japan auth information.
type YahooAuthConfigurations struct {
	ClientID    string `envconfig:"client_id"`
	Secret      string `envconfig:"secret"`
	RedirectURL string `envconfig:"redirect_url"`
}

// GetYahooAuthConfig returns yahoo japan auth configurations from environment variables.
func GetYahooAuthConfig() (YahooAuthConfigurations, error) {
	var c YahooAuthConfigurations
	err := envconfig.Process("atwell_yahoo_japan_auth", &c)

	if c.ClientID == "" || c.Secret == "" || c.RedirectURL == "" {
		return YahooAuthConfigurations{}, fmt.Errorf("failed to get database configurations. Please set them in environments variables. It is now %v", c)
	}

	return c, err
}
