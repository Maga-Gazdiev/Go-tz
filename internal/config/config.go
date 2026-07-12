package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type LogConfig struct {
	Level string
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func Load() (*Config, error) {
	loadDotEnv()

	return &Config{
		Server: ServerConfig{
			Host: envOr("SERVER_HOST", "0.0.0.0"),
			Port: envIntOr("SERVER_PORT", 8099),
		},
		Database: DatabaseConfig{
			Host:     envOr("DB_HOST", "localhost"),
			Port:     envIntOr("DB_PORT", 5432),
			User:     envOr("DB_USER", "postgres"),
			Password: envOr("DB_PASSWORD", "postgres"),
			Name:     envOr("DB_NAME", "subscriptions"),
			SSLMode:  envOr("DB_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level: envOr("LOG_LEVEL", "info"),
		},
	}, nil
}

func loadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOr(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			return p
		}
	}
	return def
}
