package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		GatewayAddr            string
		GatewayWriteTimeout    int
		GatewayReadTimeout     int
		GatewaySubAddr         string
		GatewayAuthAddr        string
		GatewayAnalyticAddr    string
		GatewayAnalyticOutAddr string
	}
)

func ParseEnv() (*Config, error) {
	writetimeout, err := strconv.Atoi(get("gateway_write_timeout_second", "15"))
	readtimeout, err := strconv.Atoi(get("gateway_read_timeout_second", "15"))
	if err != nil {
		return nil, err
	}

	return &Config{
		GatewayAddr:            get("gateway_addr", "0.0.0.0:80"),
		GatewayWriteTimeout:    writetimeout,
		GatewayReadTimeout:     readtimeout,
		GatewaySubAddr:         get("gateway_subscribes_addr", ""),
		GatewayAuthAddr:        get("gateway_authorization_addr", ""),
		GatewayAnalyticAddr:    get("gateway_analytics_addr", ""),
		GatewayAnalyticOutAddr: get("gateway_analytics_out_addr", "0.0.0.0:90"),
	}, nil
}

func get(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
