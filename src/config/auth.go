package config

import (
	"os"
	"strings"
)

// AuthConfig handles authentication configuration
type AuthConfig struct {
	AppID   string
	Secret  string
	Public  string
	Origins []string
}

func createAuthConfig() *AuthConfig {
	config := &AuthConfig{}

	if value, ok := os.LookupEnv("AUTH_APP_ID"); ok {
		config.AppID = value
	}

	if value, ok := os.LookupEnv("AUTH_SECRET"); ok {
		config.Secret = value
	}

	if value, ok := os.LookupEnv("AUTH_PUBLIC"); ok {
		config.Public = value
	}

	if value, ok := os.LookupEnv("AUTH_ORIGINS"); ok {
		config.Origins = strings.Split(value, ",")
	}

	return config
}
