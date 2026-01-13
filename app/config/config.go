package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

const (
	envBaseURL   = "ENSYNC_BASE_URL"
	envDebug     = "ENSYNC_DEBUG"
	envConfigDir = "ENSYNC_CONFIG_DIR"

	defaultConfigDirName = ".ensync"
	configFileName       = "config"
	configFileType       = "yaml"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	Debug   bool   `mapstructure:"debug"`
}

func Load() (*Config, error) {
	if err := initViperPaths(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	applyEnvironmentOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return errors.New("base_url is required: set via config file or ENSYNC_BASE_URL environment variable")
	}
	return nil
}

func initViperPaths() error {
	viper.AddConfigPath(configDir())
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return fmt.Errorf("read config file: %w", err)
		}
	}

	return nil
}

func applyEnvironmentOverrides(cfg *Config) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv(envBaseURL)
	}

	if !cfg.Debug {
		if val := os.Getenv(envDebug); val != "" {
			if parsed, err := strconv.ParseBool(val); err == nil {
				cfg.Debug = parsed
			}
		}
	}
}

func configDir() string {
	if dir := os.Getenv(envConfigDir); dir != "" {
		return dir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(home, defaultConfigDirName)
}
