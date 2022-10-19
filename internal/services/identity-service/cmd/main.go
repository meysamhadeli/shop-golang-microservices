package main

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/server"
	log "github.com/sirupsen/logrus"
	"os"
)

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

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
