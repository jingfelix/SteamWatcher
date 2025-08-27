package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config defines the application configuration
type Config struct {
	SteamAPIKey  string
	SteamUserIDs []string

	// Bark Configs
	DeviceKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file (non-fatal if fails)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Unable to load .env file: %v", err)
	}

	config := &Config{}

	// Get Steam API Key from environment variable
	config.SteamAPIKey = os.Getenv("STEAM_API_KEY")

	// Get Bark Device Key from environment variable
	config.DeviceKey = os.Getenv("BARK_DEVICE_KEY")

	// Get Steam User IDs from environment variable
	userIDsStr := os.Getenv("STEAM_USER_IDS")
	if userIDsStr != "" {
		config.SteamUserIDs = strings.Split(userIDsStr, ",")
		// Clean up whitespace
		for i, id := range config.SteamUserIDs {
			config.SteamUserIDs[i] = strings.TrimSpace(id)
		}
	}

	// Validate required configuration
	if config.SteamAPIKey == "" {
		return nil, fmt.Errorf("STEAM_API_KEY is a required configuration")
	}

	if len(config.SteamUserIDs) == 0 {
		return nil, fmt.Errorf("STEAM_USER_IDS is a required configuration")
	}

	if config.DeviceKey == "" {
		return nil, fmt.Errorf("BARK_DEVICE_KEY is a required configuration")
	}

	return config, nil
}
