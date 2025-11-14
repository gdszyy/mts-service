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
	VirtualHost  string
	AccessToken  string // Optional: UOF Access Token for whoami.xml
	Production   bool

	// OAuth
		AuthURL string
		UOFAPIBaseURL string // UOF API base URL for whoami.xml
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
			ClientID:     getEnv("MTS_CLIENT_ID", "rOCi1mVU6UhH7I2lu3t1ME4dyfkiPFVk"),
			ClientSecret: getEnv("MTS_CLIENT_SECRET", "IWGVLxyH9XcoJpZNEDNjCSFSvDAp49c_7kOi_iCuxQzitOsMfY8X4HMmw3Dcydcr"),
			BookmakerID:  getEnv("MTS_BOOKMAKER_ID", "45426"),
			VirtualHost:  getEnv("MTS_VIRTUAL_HOST", "mbs-dp-non-prod-wss"),
			AccessToken:  getEnv("UOF_ACCESS_TOKEN", ""),
		Production:   getEnvBool("MTS_PRODUCTION", false),
				AuthURL:      "", // Will be set after Production is determined
			UOFAPIBaseURL: getEnv("UOF_API_BASE_URL", "https://wss.dataplane-nonprod.sportradar.dev"),
		}

		// Set AuthURL after Production is determined
		cfg.AuthURL = getAuthURL(cfg.Production)

		// If MTS_CLIENT_ID is not set, we use the hardcoded values for testing.
		// We skip the required checks here since we are using hardcoded values.

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

func getAuthURL(production bool) string {
	if production {
		return "https://mts-api.betradar.com/api/v1/oauth/token"
	}
	return "https://mts-api-ci.betradar.com/api/v1/oauth/token"
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

