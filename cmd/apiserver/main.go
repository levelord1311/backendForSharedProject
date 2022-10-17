package main

import (
	"backendForSharedProject/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")

}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	//создание новых сертификатов
	//certificates.CreateCertAndKey()

	//конкурентный запуск http сервера для редиректа на TLS соединение
	go func() {
		if err := apiserver.StartHTTP(config); err != nil {
			log.Println("error starting http server: ", err)
			os.Exit(1)
		}
	}()

	//запуск TLS сервера
	if err := apiserver.StartTLS(config); err != nil {
		log.Println("error starting TLS server: ", err)
		os.Exit(2)
	} else {
		log.Println("OK")
	}
}
