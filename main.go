package main

// @title           Trae Go API
// @version         1.0
// @description     这是一个 Go 项目的 API 文档。
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
