package config

import "time"

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Cors     CorsConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	InternalPort string
	ExternalPort string
	RunMode      string
	Domain       string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}

type CorsConfig struct {
	AllowOrigins string
}

type JWTConfig struct {
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
	Secret                     string
	RefreshSecret              string
}
