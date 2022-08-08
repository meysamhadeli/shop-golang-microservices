package main

import (
	"flag"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/server"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/web"
	"log"
	"os"
)

const dev = "development"
const production = "production"

//https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Meysam Hadeli
// @contact.url https://github.com/meysamhadeli
// @title Products Service Api.
// @version 1.0
// @description Products Service Api.
func main() {
	flag.Parse()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = dev
	}

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(web.GetMicroserviceName(cfg))

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
