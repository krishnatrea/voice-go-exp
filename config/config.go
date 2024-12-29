package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg     *Config
	once    sync.Once
	cfgLock sync.RWMutex
)

type Config struct {
	SpeechSDK SpeechSDKConfig `mapstructure:"speech_sdk"`
}

// Microsoft Speech SDK Configuration
type SpeechSDKConfig struct {
	SubscriptionKey string `mapstructure:"subscription_key"`
	Region          string `mapstructure:"region"`
}

// Logging Configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// LoadConfig initializes the configuration singleton.
func LoadConfig(path string) (*Config, error) {
	var err error

	// Ensure configuration is loaded only once
	once.Do(func() {
		viper.SetConfigFile(path)

		if err = viper.ReadInConfig(); err != nil {
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}

		var tempConfig Config
		if err = viper.Unmarshal(&tempConfig); err != nil {
			err = fmt.Errorf("error unmarshalling config: %w", err)
			return
		}

		cfgLock.Lock()
		defer cfgLock.Unlock()
		cfg = &tempConfig
	})

	return cfg, err
}

// GetConfig retrieves the global configuration instance.
func GetConfig() *Config {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}
