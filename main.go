package main

import (
	"log"

	"trae-go/config"
	"trae-go/router"
)

func main() {
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()
	r := router.SetupRouter(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

