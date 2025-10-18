package config

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type Server struct {
	URL            string   `yaml:"url"`
	PORT           uint16   `yaml:"port"`
	EnableLogs     bool     `yaml:"enableLogs"`
	TrustedProxies []string `yaml:"trustedProxies"`
	Monitoring     bool     `yaml:"monitoring"`
}

type FileSettings struct {
	MimeTypes []string `yaml:"mimeTypes"`
	MaxSize   int64    `yaml:"maxSize"`
	UploadDir string   `yaml:"uploadDir"`
}

type JobScheduler struct {
	JobName     string `yaml:"jobName"`
	JobInterval string `yaml:"jobInterval"`
}

type Config struct {
	Environment string        `yaml:"environment"`
	Server      *Server       `yaml:"server"`
	File        *FileSettings `yaml:"file"`
	Job         *JobScheduler `yaml:"job"`
	Debugging   *Debugging    `yaml:"debugging"`
}

type Debugging struct {
	EnablePprof bool   `yaml:"enablePprof"`
	Port        uint16 `yaml:"port"`
}

func InitConfig() (*Config, error) {
	const configPath = "config.yml"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := &Config{
			Environment: "Debug",
			Server: &Server{
				URL:            "localhost",
				PORT:           8989,
				EnableLogs:     false,
				TrustedProxies: []string{"127.0.0.1", "::1"},
				Monitoring:     true,
			},
			File: &FileSettings{
				MimeTypes: []string{"image/jpeg", "image/png", "application/pdf"},
				MaxSize:   10 * 1024 * 1024, // 10 MB
				UploadDir: "./uploads",
			},
			Job: &JobScheduler{
				JobName:     "FileCleanup",
				JobInterval: "15m",
			},
			Debugging: &Debugging{
				EnablePprof: true,
				Port:        6060,
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

func GetConfig() *Config {
	config, err := InitConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
	}
	return config
}
