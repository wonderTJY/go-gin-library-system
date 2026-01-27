package config

import (
	"context"
	"time"

	"trae-go/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(AppConfig.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Book{}, &models.Student{}, &models.Book_Student{}, &models.User{}); err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
