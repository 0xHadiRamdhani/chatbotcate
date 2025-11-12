package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	WhatsApp WhatsAppConfig
	JWT      JWTConfig
	Security SecurityConfig
	Features FeaturesConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Environment  string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type WhatsAppConfig struct {
	PhoneNumberID string
	AccessToken   string
	APIVersion    string
	BaseURL       string
	WebhookSecret string
}

type JWTConfig struct {
	Secret        string
	ExpireHours   int
	RefreshExpire int
}

type SecurityConfig struct {
	RateLimitPerMinute int
	BcryptCost         int
	EnableCORS         bool
}

type FeaturesConfig struct {
	EnableGames      bool
	EnableBusiness   bool
	EnableUtils      bool
	EnableModeration bool
	MaxBroadcastSize int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			DBName:       getEnv("DB_NAME", "whatsapp_bot"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getInt("DB_MAX_IDLE_CONNS", 5),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		WhatsApp: WhatsAppConfig{
			PhoneNumberID: getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
			AccessToken:   getEnv("WHATSAPP_ACCESS_TOKEN", ""),
			APIVersion:    getEnv("WHATSAPP_API_VERSION", "v18.0"),
			BaseURL:       getEnv("WHATSAPP_BASE_URL", "https://graph.facebook.com"),
			WebhookSecret: getEnv("WHATSAPP_WEBHOOK_SECRET", ""),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours:   getInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpire: getInt("JWT_REFRESH_EXPIRE", 168),
		},
		Security: SecurityConfig{
			RateLimitPerMinute: getInt("RATE_LIMIT_PER_MINUTE", 60),
			BcryptCost:         getInt("BCRYPT_COST", 10),
			EnableCORS:         getBool("ENABLE_CORS", true),
		},
		Features: FeaturesConfig{
			EnableGames:      getBool("ENABLE_GAMES", true),
			EnableBusiness:   getBool("ENABLE_BUSINESS", true),
			EnableUtils:      getBool("ENABLE_UTILS", true),
			EnableModeration: getBool("ENABLE_MODERATION", true),
			MaxBroadcastSize: getInt("MAX_BROADCAST_SIZE", 1000),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}