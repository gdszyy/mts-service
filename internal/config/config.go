package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port string

	// MTS
	ClientID     string
	ClientSecret string
	BookmakerID  string
	AccessToken  string // Optional: UOF Access Token for whoami.xml
	Production   bool

	// OAuth
	AuthURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		ClientID:     getEnv("MTS_CLIENT_ID", ""),
		ClientSecret: getEnv("MTS_CLIENT_SECRET", ""),
		BookmakerID:  getEnv("MTS_BOOKMAKER_ID", ""),
		AccessToken:  getEnv("UOF_ACCESS_TOKEN", ""),
		Production:   getEnvBool("MTS_PRODUCTION", false),
		AuthURL:      "https://auth.sportradar.com/oauth/token",
	}

	if cfg.ClientID == "" {
		return nil, fmt.Errorf("MTS_CLIENT_ID is required")
	}
	if cfg.ClientSecret == "" {
		return nil, fmt.Errorf("MTS_CLIENT_SECRET is required")
	}
	// BookmakerID is optional if AccessToken is provided (will be fetched from whoami.xml)
	if cfg.BookmakerID == "" && cfg.AccessToken == "" {
		return nil, fmt.Errorf("either MTS_BOOKMAKER_ID or UOF_ACCESS_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

