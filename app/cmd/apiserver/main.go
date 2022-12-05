package main

import (
	"backendForSharedProject/internal/app/apiserver"
	"backendForSharedProject/internal/config"
	"backendForSharedProject/pkg/logging"
	"os"
)

func main() {
	// TODO "Пионерский код. Логер без закрытия в defer c потерей данных логирования" - проверить, исправить
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("initializing config...")
	cfg := config.GetConfig()

	if err := apiserver.StartMainHTTP(cfg); err != nil {
		logger.Errorln("error starting http server: ", err)
		os.Exit(2)
	}
}

// TODO Makefile
// TODO swagger
// TODO отрефакторить из монолита в микросервисы
