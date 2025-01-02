package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Debug   bool   `mapstructure:"debug"`
}

func Load() (*Config, error) {
	config := &Config{}

	if err := loadConfigFile(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyEnvVariables(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigFile() error {
	configDir := getConfigDir()
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}
	return nil
}

func applyEnvVariables(config *Config) {
	if config.APIKey == "" {
		config.APIKey = os.Getenv("ENSYNC_API_KEY")
	}
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("ENSYNC_BASE_URL")
	}
	if !config.Debug {
		if debugEnv := os.Getenv("ENSYNC_DEBUG"); debugEnv != "" {
			if debugValue, err := strconv.ParseBool(debugEnv); err == nil {
				config.Debug = debugValue
			}
		}
	}
}

func validateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if config.BaseURL == "" {
		return fmt.Errorf("Base URL is required")
	}
	return nil
}

func getConfigDir() string {
	if configDir := os.Getenv("ENSYNC_CONFIG_DIR"); configDir != "" {
		return configDir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(home, ".ensync")
}
