package main

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/server"
	"log"
	"os"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	err = server.NewServer(appLogger, cfg).Run()
	if err != nil {
		appLogger.Fatal(err)
	}
}
