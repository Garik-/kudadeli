package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	defaultDatabase       = "data.db"
	defaultToken          = ""
	defaultAllowUsers     = ""
	serviceName           = "kudadeli"
	defaultHTTPAddr       = ":8080"
	defaultEnableBot      = true
	defaultAllowedOrigins = "http://localhost:3000,http://localhost:5173"
)

type Service struct {
	Name    string
	Version string
}

type Config struct {
	Addr           string
	Database       string
	Service        Service
	LogLevel       slog.Level
	Token          string
	AllowedUsers   []int64
	EnableBot      bool
	AllowedOrigins []string
}

func envString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func envBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}

	return defaultValue
}

func envStringSlice(key string, defaultValue []string, sep string) []string {
	if value, exists := os.LookupEnv(key); exists {
		parts := strings.Split(value, sep)

		result := make([]string, 0, len(parts))

		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}

		if len(result) == 0 {
			return defaultValue
		}

		return result
	}

	return defaultValue
}

func parseAllowedUsers(input string) []int64 {
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
		Addr:     envString(prefix+"ADDR", defaultHTTPAddr),
		Database: envString(prefix+"DATABASE", defaultDatabase),
		Token:    envString(prefix+"TOKEN", defaultToken),
		Service: Service{
			Name:    serviceName,
			Version: version,
		},
		AllowedUsers: parseAllowedUsers(envString(prefix+"USERS", defaultAllowUsers)),
		EnableBot:    envBool(prefix+"ENABLE_BOT", defaultEnableBot),
		AllowedOrigins: envStringSlice(prefix+"ALLOWED_ORIGINS", []string{
			"http://localhost:3000",
			"http://localhost:5173",
		}, ","),
	}
}
