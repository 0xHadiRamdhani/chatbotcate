package database

import (
	"fmt"
	"time"

	"whatsapp-bot/internal/config"
	"whatsapp-bot/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var (
	db          *gorm.DB
	redisClient *redis.Client
)

func Initialize(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	database, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	database.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	database.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	database.DB().SetConnMaxLifetime(5 * time.Minute)

	// Enable logging
	database.LogMode(true)

	// Auto migrate models
	if err := migrate(database); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	db = database
	return database, nil
}

func InitializeRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := client.Context()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	redisClient = client
	return client, nil
}

func migrate(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.UserPreferences{},
		&models.Contact{},
		&models.Message{},
		&models.AutoReply{},
		&models.Broadcast{},
		&models.BroadcastRecipient{},
		&models.Group{},
		&models.GroupMember{},
		&models.GameScore{},
		&models.Quiz{},
		&models.QuizSession{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Reminder{},
		&models.BlockedWord{},
		&models.SpamReport{},
		&models.WeatherData{},
		&models.CurrencyRate{},
		&models.Analytics{},
		&models.CustomCommand{},
		&models.Template{},
		&models.SystemLog{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}

	logrus.Info("Database migration completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func GetRedis() *redis.Client {
	return redisClient
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}