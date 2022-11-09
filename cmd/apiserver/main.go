package main

import (
	"backendForSharedProject/internal/app/apiserver"
	"log"
	"os"
)

func main() {

	config, err := apiserver.NewConfig()
	if err != nil {
		log.Println("error creating config file:", err)
		os.Exit(1)
	}

	//запуск главного HTTP сервера
	if err := apiserver.StartMainHTTP(config); err != nil {
		log.Println("error starting http server: ", err)
		os.Exit(2)
	}
}
