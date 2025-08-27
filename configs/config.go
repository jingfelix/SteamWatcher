package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config defines the application configuration
type Config struct {
	SteamAPIKey  string   `mapstructure:"STEAM_API_KEY"`
	SteamUserIDs []string `mapstructure:"STEAM_USER_IDS"`

	// Bark Configs
	DeviceKey string `mapstructure:"BARK_DEVICE_KEY"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Initialize Viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Read config (optional)
	_ = viper.ReadInConfig()

	// Unmarshall configuration
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	// Validate required configuration
	if config.SteamAPIKey == "" {
		return nil, fmt.Errorf("STEAM_API_KEY is a required configuration")
	}

	// Process Steam User IDs (for comma-separated string support)
	if len(config.SteamUserIDs) == 1 && strings.Contains(config.SteamUserIDs[0], ",") {
		config.SteamUserIDs = strings.Split(config.SteamUserIDs[0], ",")
		// Clean up whitespace
		for i, id := range config.SteamUserIDs {
			config.SteamUserIDs[i] = strings.TrimSpace(id)
		}
	}

	if len(config.SteamUserIDs) == 0 {
		return nil, fmt.Errorf("STEAM_USER_IDS is a required configuration")
	}

	if config.DeviceKey == "" {
		return nil, fmt.Errorf("BARK_DEVICE_KEY is a required configuration")
	}

	return config, nil
}
