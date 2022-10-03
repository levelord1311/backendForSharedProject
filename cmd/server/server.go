package main

import (
	c "backendForSharedProject/config"
	h "backendForSharedProject/internal/handlers"
	j "backendForSharedProject/internal/jwt"
	"database/sql"
	"fmt"
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
	db := connectToDB(mainConfig.Path)
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

	http.HandleFunc("/home", j.VerifyJWT(h.HandlePage))
	http.HandleFunc("/", h.DefaultHandler)
	http.HandleFunc("/auth", h.AuthorizationHandler)

	if err := http.ListenAndServeTLS(mainConfig.HttpsPort, "./cert.pem", "./key.pem", nil); err != nil {
		log.Printf("Error listening to https port %s	: %s", mainConfig.HttpsPort, err)
		os.Exit(6)
	}
}

func connectToDB(path string) *sql.DB {
	//загрузка конфига с параметрами БД
	psqlConfig, err := c.LoadPSQLConfig(path)
	if err != nil {
		log.Println("error reading psql config:", err)
		os.Exit(2)
	}

	//форматирование параметров для подключения к БД
	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		psqlConfig.Host, psqlConfig.Port, psqlConfig.User, psqlConfig.Password, psqlConfig.DBName)

	//открытие БД
	db, err := sql.Open("postgres", psqlConnectionString)
	if err != nil {
		log.Println("error opening database:", err)
		os.Exit(3)
	}

	//подключение к БД
	err = db.Ping()
	if err != nil {
		log.Println("error connecting to database:", err)
		os.Exit(4)
	}
	log.Println("DB successfully connected.")
	return db
}
