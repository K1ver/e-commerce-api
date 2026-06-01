package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func Load() (*Config, error) {
	accessTTL, err := parseDurationEnv("JWT_ACCESS_EXPIRE", 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshTTL, err := parseDurationEnv("JWT_REFRESH_EXPIRE", 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := parseDurationEnv("POSTGRES_CONN_MAX_LIFETIME", time.Hour)
	if err != nil {
		return nil, err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtSecret == "" || jwtRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET and JWT_REFRESH_SECRET are required")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	maxIdleConns, _ := strconv.Atoi(getEnv("POSTGRES_MAX_IDLE_CONNS", "5"))
	maxOpenConns, _ := strconv.Atoi(getEnv("POSTGRES_MAX_OPEN_CONNS", "25"))
	redisPoolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))

	return &Config{
		Server: ServerConfig{
			InternalPort: getEnv("SERVER_INTERNAL_PORT", "8080"),
			ExternalPort: getEnv("SERVER_EXTERNAL_PORT", "8080"),
			RunMode:      getEnv("SERVER_RUN_MODE", "debug"),
			Domain:       getEnv("SERVER_DOMAIN", "localhost"),
		},
		Postgres: PostgresConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnv("POSTGRES_PORT", "5432"),
			User:            getEnv("POSTGRES_USER", "postgres"),
			Password:        os.Getenv("POSTGRES_PASSWORD"),
			DbName:          getEnv("POSTGRES_DB", "ecommerce"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxIdleConns:    maxIdleConns,
			MaxOpenConns:    maxOpenConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			Db:           redisDB,
			DialTimeout:  mustDuration(getEnv("REDIS_DIAL_TIMEOUT", "5s")),
			ReadTimeout:  mustDuration(getEnv("REDIS_READ_TIMEOUT", "3s")),
			WriteTimeout: mustDuration(getEnv("REDIS_WRITE_TIMEOUT", "3s")),
			PoolSize:     redisPoolSize,
			PoolTimeout:  mustDuration(getEnv("REDIS_POOL_TIMEOUT", "4s")),
		},
		Cors: CorsConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		},
		JWT: JWTConfig{
			AccessTokenExpireDuration:  accessTTL,
			RefreshTokenExpireDuration: refreshTTL,
			Secret:                     jwtSecret,
			RefreshSecret:              jwtRefreshSecret,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return time.ParseDuration(v)
}

func mustDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}
