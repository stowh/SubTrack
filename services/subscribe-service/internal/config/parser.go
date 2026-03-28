package config

import (
	"os"
)

type (
	Config struct {
		ServerAddr string

		PostgresAddr string
		PostgresUser string
		PostgresPass string
		PostgresName string
	}
)

func ParseEnv() (*Config, error) {
	return &Config{
		ServerAddr:   get("subscribes_addr", "0.0.0.0:80"),
		PostgresAddr: get("postgres_addr", ""),
		PostgresUser: get("postgres_user", ""),
		PostgresPass: get("postgres_pass", ""),
		PostgresName: get("postgres_name", ""),
	}, nil
}

func get(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
