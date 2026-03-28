package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		ServerAddr         string
		JwtSecret          string
		JwtExpiredMin      int
		RefreshExpiredDays int

		PostgresAddr string
		PostgresUser string
		PostgresPass string
		PostgresName string
	}
)

func ParseEnv() (*Config, error) {
	expired, err := strconv.Atoi(get("authorization_jwt_expired_min", ""))
	if err != nil {
		return nil, err
	}

	expiredRefresh, err := strconv.Atoi(get("authorization_refresh_expired_days", ""))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr:         get("authorization_addr", "0.0.0.0:80"),
		JwtSecret:          get("authorization_jwt_secret", ""),
		JwtExpiredMin:      expired,
		RefreshExpiredDays: expiredRefresh,

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
