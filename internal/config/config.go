package config

import "os"

type Config struct {
	HTTPPort    string
	RedisAddr   string
	PostgresDSN string
}

func Load() Config {
	return Config{
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
