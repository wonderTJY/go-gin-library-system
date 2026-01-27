package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Cors      CorsConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"ratelimit"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
	DSN      string `mapstructure:"dsn"` // 兼容 SQLite 或直接提供 DSN 的情况

	MaxIdleConns    int64  `mapstructure:"MaxIdleConns"`
	MaxOpenConns    int64  `mapstructure:"MaxOpenConns"`
	ConnMaxLifetime string `mapstructure:"ConnMaxLifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AuthConfig struct {
	TokenExpireHours string `mapstructure:"token_expire_hours"`
}

type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type RateLimitConfig struct {
	GlobalLimit int64 `mapstructure:"global_limit"`
	IPLimit     int64 `mapstructure:"ip_limit"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	log.Println("Configuration loaded successfully")
}
