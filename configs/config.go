package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type AppConfig struct {
	Name string `envconfig:"APP_NAME" default:"vocabunny-core-api"`
	Env  string `envconfig:"APP_ENV" default:"local"`
}

type HTTPConfig struct {
	Port            string        `envconfig:"HTTP_PORT" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	WriteTimeout    time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"15s"`
	ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"10s"`
}

type DatabaseConfig struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            string        `envconfig:"DB_PORT" default:"5432"`
	User            string        `envconfig:"DB_USER" default:"postgres"`
	Password        string        `envconfig:"DB_PASSWORD" default:"postgres"`
	Name            string        `envconfig:"DB_NAME" default:"vocabunny"`
	SSLMode         string        `envconfig:"DB_SSLMODE" default:"disable"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"50"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"30m"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
	)
}

type RedisConfig struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type JWTConfig struct {
	Secret          string        `envconfig:"JWT_SECRET" default:"change-me"`
	Issuer          string        `envconfig:"JWT_ISSUER" default:"vocabunny-core-api"`
	Audience        string        `envconfig:"JWT_AUDIENCE" default:"vocabunny-clients"`
	AccessTokenTTL  time.Duration `envconfig:"JWT_ACCESS_TOKEN_TTL" default:"24h"`
	RefreshTokenTTL time.Duration `envconfig:"JWT_REFRESH_TOKEN_TTL" default:"720h"`
}

type StorageConfig struct {
	Provider string `envconfig:"STORAGE_PROVIDER" default:"local"`
	BasePath string `envconfig:"STORAGE_BASE_PATH" default:"./storage"`
	BaseURL  string `envconfig:"STORAGE_BASE_URL" default:"http://localhost:8080/storage"`
}

func Load() (*Config, error) {
	envFilePath := os.Getenv("ENV_FILE_PATH")
	if envFilePath == "" {
		envFilePath = ".env"
	}

	if err := godotenv.Load(envFilePath); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, fmt.Errorf("load env file: %w", err)
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
