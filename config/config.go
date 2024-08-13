package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server struct {
		Port string
	}
	DB struct {
		Username string
		Password string
		Host     string
		PortDb   string
		DBName   string
		SSLMode  string
	}
}

func LoadConfig() *Config {
	cfg := &Config{
		Server: struct{ Port string }{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: struct {
			Username string
			Password string
			Host     string
			PortDb   string
			DBName   string
			SSLMode  string
		}{
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", ""),
			PortDb:   getEnv("DB_PORTDB", "5432"),
			DBName:   getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	return cfg
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB.Username,
		c.DB.Password,
		c.DB.Host,
		c.DB.PortDb,
		c.DB.DBName,
		c.DB.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
