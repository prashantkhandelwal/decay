package config

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type Server struct {
	PORT           string   `yaml:"port"`
	EnableLogs     bool     `yaml:"enableLogs"`
	TrustedProxies []string `yaml:"trustedProxies"`
}

type Config struct {
	Environment string  `yaml:"environment"`
	Server      *Server `yaml:"server"`
}

func InitConfig() (*Config, error) {
	const configPath = "config.yml"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := &Config{
			Environment: "Debug",
			Server: &Server{
				PORT:           "8989",
				EnableLogs:     false,
				TrustedProxies: []string{"127.0.0.1", "::1"},
			},
		}

		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("error marshaling default config: %w", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("error writing config file: %w", err)
		}

		log.Println("Default config file created.")
	}

	// Load config from file
	fileData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
