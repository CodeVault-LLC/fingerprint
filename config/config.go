package config

import (
	"os"

	"github.com/codevault-llc/fingerprint/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type InternalConfig struct {
	ServicePort string
	ServiceHost string

	// Database
	DatabaseAddr string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

var Config *InternalConfig

func NewInternalConfig() (*InternalConfig, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Error loading .env file", zap.Error(err))
		return nil, err
	}

	config := &InternalConfig{
		ServicePort: getEnv("SERVICE_PORT", "50051"),
		ServiceHost: getEnv("SERVICE_HOST", "localhost"),

		// Database
		DatabaseAddr: getEnv("DATABASE_ADDR", "localhost:9000"),
		DatabaseUser: getEnv("DATABASE_USER", "default"),
		DatabasePass: getEnv("DATABASE_PASS", ""),
		DatabaseName: getEnv("DATABASE_NAME", "fingerprint"),
	}

	Config = config
	return config, nil
}
