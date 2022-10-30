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
	//deprecated: HTTPS не поддерживается на бесплатном heroku
	//config := apiserver.NewConfig()
	//_, err := toml.DecodeFile(configPath, config)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//создание новых сертификатов
	//certificates.CreateCertAndKey()
	//
	//конкурентный запуск http сервера для редиректа на TLS соединение
	//go func() {
	//	if err := apiserver.StartHTTP(config); err != nil {
	//		log.Println("error starting http server: ", err)
	//		os.Exit(1)
	//	}
	//}()
	//
	//запуск TLS сервера
	//if err := apiserver.StartTLS(config); err != nil {
	//	log.Println("error starting TLS server: ", err)
	//	os.Exit(2)
	//} else {
	//	log.Println("OK")
	//}

	//запуск главного HTTP сервера
	if err := apiserver.StartMainHTTP(config); err != nil {
		log.Println("error starting http server: ", err)
		os.Exit(2)
	}
}
