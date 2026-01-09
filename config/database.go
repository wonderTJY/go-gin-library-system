package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"trae-go/models"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Book{}, &models.Student{}, &models.Book_Student{}); err != nil {
		return nil, err
	}
	return db, nil
}
