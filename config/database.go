package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"trae-go/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	//适配
	var dialector gorm.Dialector
	switch AppConfig.Database.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			AppConfig.Database.Host,
			AppConfig.Database.User,
			AppConfig.Database.Password,
			AppConfig.Database.DBName,
			AppConfig.Database.Port,
			AppConfig.Database.SSLMode,
			AppConfig.Database.TimeZone,
		)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(AppConfig.Database.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", AppConfig.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	d, err := time.ParseDuration(AppConfig.Database.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("MaxLifetime parse failed,use default setting")
		d = time.Hour
	}
	sqlDB.SetMaxIdleConns(int(AppConfig.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(AppConfig.Database.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(d)
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
