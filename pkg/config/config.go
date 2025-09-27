package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port string
	Env  string

	DBDriver string
	DBDSN    string

	JWTAccessSecret   string
	JWTRefreshSecret  string
	JWTAccessTTLMin   int
	JWTRefreshTTLDays int
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func atoi(s string, def int) int {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return def
	}
	return v
}

func Load() Config {
	// optional: load from .env if you use a loader like godotenv
	cfg := Config{
		Port: getenv("PORT", "8080"),
		Env:  getenv("ENV", "dev"),

		DBDriver: getenv("DB_DRIVER", "sqlite3"),
		DBDSN:    getenv("DB_DSN", "./chat.db"),

		JWTAccessSecret:   getenv("JWT_ACCESS_SECRET", "dev-access"),
		JWTRefreshSecret:  getenv("JWT_REFRESH_SECRET", "dev-refresh"),
		JWTAccessTTLMin:   atoi(getenv("JWT_ACCESS_TTL_MIN", "15"), 15),
		JWTRefreshTTLDays: atoi(getenv("JWT_REFRESH_TTL_DAYS", "7"), 7),
	}

	log.Printf("loaded config: env=%s driver=%s", cfg.Env, cfg.DBDriver)
	return cfg
}
