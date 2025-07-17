package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	defaultDatabase   = "data.db"
	defaultToken      = ""
	defaultAllowUsers = ""
	serviceName       = "kudadeli"
	defaultHTTPAddr   = ":8080"
)

type Service struct {
	Name    string
	Version string
}

type Config struct {
	Addr       string
	Database   string
	Service    Service
	LogLevel   slog.Level
	Token      string
	AllowUsers []int64
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func parseAllowUsers(input string) []int64 {
	if input == "" {
		return nil
	}

	parts := strings.Split(input, ",")
	if len(parts) == 0 {
		return nil
	}

	// Заранее выделяем память под слайс
	userIDs := make([]int64, 0, len(parts)) // len=0, cap=len(parts)

	for _, part := range parts {
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			continue
		}

		userIDs = append(userIDs, id) // Добавляем в pre-allocated слайс
	}

	if len(userIDs) == 0 {
		return nil
	}

	return userIDs
}

func New(version string) *Config {
	prefix := strings.ToUpper(serviceName) + "_"

	return &Config{
		LogLevel: slog.LevelDebug,
		Addr:     getEnv(prefix+"ADDR", defaultHTTPAddr),
		Database: getEnv(prefix+"DATABASE", defaultDatabase),
		Token:    getEnv(prefix+"TOKEN", defaultToken),
		Service: Service{
			Name:    serviceName,
			Version: version,
		},
		AllowUsers: parseAllowUsers(getEnv(prefix+"USERS", defaultAllowUsers)),
	}
}
