package main

import (
	"flag"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/server"
	"log"
	"math/rand"
	"os"
)

const dev = "development"
const production = "production"

func main() {

	a := rand.Intn(10)
	for a < 100 {
		if a%5 == 0 {
			goto done
		}
		a = a*2 + 1
	}

	fmt.Println("do something when the loop completes normally")
done:
	fmt.Println("do complicated stuff no matter why we left the loop")
	fmt.Println(a)

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
	appLogger.Info(config.GetMicroserviceName(cfg))

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
