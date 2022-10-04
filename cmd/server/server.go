package main

import (
	c "backendForSharedProject/config"
	h "backendForSharedProject/internal/handlers"
	j "backendForSharedProject/internal/jwt"
	"backendForSharedProject/internal/store"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	//загрузка главного конфига
	mainConfig, err := c.LoadMainConfig("../../config")
	if err != nil {
		log.Println("error reading main config:", err)
		os.Exit(1)
	}
	log.Println("mainConfig", mainConfig)

	//генерация нового сертификата и ключа
	//https.CreateCertAndKey()

	//подключение к БД
	db := store.ConnectToDB(mainConfig.Path)
	//??? нужно ли обрабатывать ошибки при закрытии элементов при помощи defer в main() программе?
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	//конкурентный запуск http сервера для редиректа на https соединение
	go func() {
		if err := http.ListenAndServe(mainConfig.HttpPort, http.HandlerFunc(h.RedirectToTls)); err != nil {
			log.Printf("Error listening to http port %s	: %s", mainConfig.HttpsPort, err)
			os.Exit(5)
		}
	}()
	datab := store.Database{db}

	http.HandleFunc("/home", j.VerifyJWT(h.HandlePage))
	http.HandleFunc("/", h.DefaultHandler)
	http.HandleFunc("/auth", datab.AuthorizationHandler)

	if err := http.ListenAndServeTLS(mainConfig.HttpsPort, "./cert.pem", "./key.pem", nil); err != nil {
		log.Printf("Error listening to https port %s	: %s", mainConfig.HttpsPort, err)
		os.Exit(6)
	}
}
