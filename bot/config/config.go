package config

import (
	"log/slog"
	"os"
	"strings"
)

const (
	defaultDatabase = "data.db"
	defaultToken    = ""
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
	Token    string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func New(version string) *Config {
	prefix := strings.ToUpper(serviceName) + "_"

	return &Config{
		LogLevel: slog.LevelDebug,
		Database: getEnv(prefix+"DATABASE", defaultDatabase),
		Token:    getEnv(prefix+"TOKEN", defaultToken),
		Service: Service{
			Name:    serviceName,
			Version: version,
		},
	}
}
