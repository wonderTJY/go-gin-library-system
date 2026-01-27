package main

import (
	"log"

	"trae-go/config"
	"trae-go/router"
)

func main() {
	config.InitConfig()

	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	rdb, err := config.InitRedis()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	defer rdb.Close()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()

	r := router.SetupRouter(db, rdb)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
