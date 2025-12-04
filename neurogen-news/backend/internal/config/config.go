package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	ServeStatic bool   `mapstructure:"SERVE_STATIC"`
	BaseURL     string `mapstructure:"BASE_URL"`

	// Database
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	// Redis
	RedisURL string `mapstructure:"REDIS_URL"`

	// JWT
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTAccessExpiry  int    `mapstructure:"JWT_ACCESS_EXPIRY"`  // minutes
	JWTRefreshExpiry int    `mapstructure:"JWT_REFRESH_EXPIRY"` // days

	// CORS
	CORSOrigins string `mapstructure:"CORS_ORIGINS"`

	// Meilisearch
	MeilisearchURL string `mapstructure:"MEILISEARCH_URL"`
	MeilisearchKey string `mapstructure:"MEILISEARCH_KEY"`

	// S3 / Object Storage
	S3Enabled         bool   `mapstructure:"S3_ENABLED"`
	S3Region          string `mapstructure:"S3_REGION"`
	S3Bucket          string `mapstructure:"S3_BUCKET"`
	S3AccessKeyID     string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3Endpoint        string `mapstructure:"S3_ENDPOINT"` // For S3-compatible (MinIO, DO Spaces)
	S3CDNBaseURL      string `mapstructure:"S3_CDN_BASE_URL"`

	// Uploads (local fallback)
	UploadPath    string `mapstructure:"UPLOAD_PATH"`
	MaxUploadSize int64  `mapstructure:"MAX_UPLOAD_SIZE"` // bytes

	// OAuth - Google
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`

	// OAuth - VK
	VKClientID     string `mapstructure:"VK_CLIENT_ID"`
	VKClientSecret string `mapstructure:"VK_CLIENT_SECRET"`
	VKRedirectURL  string `mapstructure:"VK_REDIRECT_URL"`

	// OAuth - Telegram
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`

	// OAuth - GitHub
	GithubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURL  string `mapstructure:"GITHUB_REDIRECT_URL"`

	// Email (for notifications)
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	// Logging
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("SERVE_STATIC", true)
	viper.SetDefault("BASE_URL", "http://localhost:8080")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/neurogen?sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379/0")
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")
	viper.SetDefault("JWT_ACCESS_EXPIRY", 15)  // 15 minutes
	viper.SetDefault("JWT_REFRESH_EXPIRY", 7)  // 7 days
	viper.SetDefault("CORS_ORIGINS", "http://localhost:3000")
	viper.SetDefault("MEILISEARCH_URL", "http://localhost:7700")
	viper.SetDefault("MEILISEARCH_KEY", "")
	viper.SetDefault("S3_ENABLED", false)
	viper.SetDefault("S3_REGION", "us-east-1")
	viper.SetDefault("UPLOAD_PATH", "./uploads")
	viper.SetDefault("MAX_UPLOAD_SIZE", 10*1024*1024) // 10MB
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SMTP_PORT", 587)

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read from config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is not an error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

