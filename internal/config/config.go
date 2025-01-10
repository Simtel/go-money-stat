package config

import (
	"os"
)

type ZenMoneyConfig struct {
	Token string
}

type Config struct {
	ZenMoney ZenMoneyConfig
}

func New() *Config {
	return &Config{
		ZenMoney: ZenMoneyConfig{
			Token: getEnv("ZENMONEY_TOKEN", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
