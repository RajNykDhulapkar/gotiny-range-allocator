package config

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// RangeConfig holds range-specific configuration
type RangeConfig struct {
	DefaultSize int64 `mapstructure:"default_size"`
	MinSize     int64 `mapstructure:"min_size"`
	MaxSize     int64 `mapstructure:"max_size"`
}

type Config struct {
	GRPCPort    string      `mapstructure:"grpc_port"`
	DatabaseURL string      `mapstructure:"database_url"`
	Range       RangeConfig `mapstructure:"range"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig loads configuration once and caches it
func LoadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		v := viper.New()

		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Info: .env file not found, proceeding with environment variables and defaults")
		} else {
			log.Println("Info: .env file loaded successfully")
		}

		// Read from environment variables
		v.SetEnvPrefix("RANGE_ALLOCATOR")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		// Set default values
		v.SetDefault("grpc_port", "50051")
		v.SetDefault("database_url", "")
		// range
		v.SetDefault("range.default_size", 1000)
		v.SetDefault("range.min_size", 100)
		v.SetDefault("range.max_size", 10000)

		config = &Config{}

		// Unmarshal the configuration into the Config struct
		if err = v.Unmarshal(config); err != nil {
			log.Printf("Error unmarshaling config: %v", err)
			config = nil
			return
		}
	})

	return config, err
}

func ValidateConfig(cfg *Config) error {
	if config == nil {
		return fmt.Errorf("Configuration is not properly initialized")
	}

	if cfg.GRPCPort == "" {
		return fmt.Errorf("RANGE_ALLOCATOR_GRPC_PORT is required")
	}
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("RANGE_ALLOCATOR_DATABASE_URL is required")
	}
	if cfg.Range.DefaultSize <= 0 {
		return fmt.Errorf("range_default_size must be positive")
	}
	if cfg.Range.MinSize <= 0 {
		return fmt.Errorf("range_min_size must be positive")
	}
	if cfg.Range.MaxSize < cfg.Range.MinSize {
		return fmt.Errorf("range_max_size must be greater than or equal to range_min_size")
	}
	return nil
}
