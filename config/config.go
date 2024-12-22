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
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePass     string
	DatabaseKeyspace string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

var Config *InternalConfig

func NewInternalConfig(path string) (*InternalConfig, error) {
	err := godotenv.Load(path)
	if err != nil {
		logger.Log.Error("Error loading .env file", zap.Error(err))
		return nil, err
	}

	config := &InternalConfig{
		ServicePort: getEnv("SERVICE_PORT", "50051"),
		ServiceHost: getEnv("SERVICE_HOST", "localhost"),

		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "9042"),
		DatabaseUser:     getEnv("DATABASE_USER", "cassandra"),
		DatabasePass:     getEnv("DATABASE_PASS", "cassandra"),
		DatabaseKeyspace: getEnv("DATABASE_KEYSPACE", "fingerprint"),
	}

	Config = config
	return config, nil
}
