package config

import (
	"log/slog"
	"os"
	"strings"
)

const (
	defaultDatabase = "data.db"
	serviceName     = "kudadeli"
)

type Service struct {
	Name    string
	Version string
}

type Config struct {
	Database string
	Service  Service
	LogLevel slog.Level
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func New(version string) *Config {
	prefix := strings.ToUpper(serviceName) + "_"

	return &Config{
		LogLevel: slog.LevelDebug,
		Database: getEnv(prefix+"DATABASE", defaultDatabase),
		Service: Service{
			Name:    serviceName,
			Version: version,
		},
	}
}
