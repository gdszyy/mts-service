package config

import (
	"fmt"
	"os"
	"strconv"
	"log"
	"github.com/gdsZyy/mts-service/internal/client"
)

type Config struct {
	// Server
	Port string

		// MTS
		ClientID     string
		ClientSecret string
			BookmakerID  string
			LimitID      string
			OperatorID   int64 // Changed to int64
			VirtualHost  string
			WSURL        string
			WSAudience   string
		AccessToken  string // Optional: UOF Access Token for whoami.xml
		Production   bool

	// OAuth
		AuthURL string
		UOFAPIBaseURL string // UOF API base URL for whoami.xml
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
			ClientID:     getEnv("MTS_CLIENT_ID", ""),
				ClientSecret: getEnv("MTS_CLIENT_SECRET", ""),
					BookmakerID:  getEnv("MTS_BOOKMAKER_ID", ""),
					LimitID:      getEnv("MTS_LIMIT_ID", ""),
					OperatorID:   getEnvInt64("MTS_OPERATOR_ID", 0), // Changed to getEnvInt64
					VirtualHost:  getEnv("MTS_VIRTUAL_HOST", ""),
				WSURL:        getEnv("MTS_WS_URL", "wss://wss.dataplane-nonprod.sportradar.dev"),
				WSAudience:   getEnv("MTS_WS_AUDIENCE", "mbs-dp-non-prod-wss"),
			AccessToken:  getEnv("UOF_ACCESS_TOKEN", ""),
		Production:   getEnvBool("MTS_PRODUCTION", false),
				AuthURL:      getEnv("MTS_AUTH_URL", "https://auth.sportradar.com/oauth/token"),
			UOFAPIBaseURL: getEnv("UOF_API_BASE_URL", "https://global.api.betradar.com"),
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

	// If AccessToken is provided, try to fetch Bookmaker ID and VirtualHost
	if cfg.AccessToken != "" && (cfg.BookmakerID == "" || cfg.VirtualHost == "") {
			log.Println("Bookmaker ID or VirtualHost not provided, fetching from whoami.xml...")
			bookmakerID, virtualHost, err := client.FetchBookmakerInfo(cfg.AccessToken, cfg.UOFAPIBaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Bookmaker Info: %w", err)
		}
		if cfg.BookmakerID == "" {
			cfg.BookmakerID = bookmakerID
		}
		if cfg.VirtualHost == "" {
			cfg.VirtualHost = virtualHost
		}
		log.Printf("Bookmaker Info fetched successfully: BookmakerID=%s, VirtualHost=%s", cfg.BookmakerID, cfg.VirtualHost)
	}

	// Final check for required fields
	if cfg.BookmakerID == "" {
		return nil, fmt.Errorf("MTS_BOOKMAKER_ID is required and could not be fetched")
	}
	if cfg.VirtualHost == "" {
		return nil, fmt.Errorf("MTS_VIRTUAL_HOST is required and could not be fetched")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}



	func getEnvInt64(key string, defaultValue int64) int64 {
		if value := os.Getenv(key); value != "" {
			if i, err := strconv.ParseInt(value, 10, 64); err == nil {
				return i
			}
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

