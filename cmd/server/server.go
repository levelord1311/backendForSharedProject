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
		fmt.Println("error reading main config:", err)
		os.Exit(1)
	}
	fmt.Println("mainConfig", mainConfig)

	//загрузка конфига с параметрами БД
	psqlConfig, err := c.LoadPSQLConfig(mainConfig.Path)
	if err != nil {
		fmt.Println("error reading psql config:", err)
		os.Exit(2)
	}

	//форматирование параметров для подключения к БД
	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		psqlConfig.Host, psqlConfig.Port, psqlConfig.User, psqlConfig.Password, psqlConfig.DBName)

	//открытие БД
	db, err := sql.Open("postgres", psqlConnectionString)
	if err != nil {
		fmt.Println("error opening database:", err)
		os.Exit(3)
	}
	//??? нужно ли обрабатывать ошибки при закрытии элементов при помощи defer в main() программе?
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	//подключение к БД
	err = db.Ping()
	if err != nil {
		fmt.Println("error connecting to database:", err)
		os.Exit(4)
	}
	fmt.Println("DB successfully connected.")

	http.HandleFunc("/home", j.VerifyJWT(h.HandlePage))
	http.HandleFunc("/", h.DefaultHandler)
	http.HandleFunc("/auth", h.AuthPage)

	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Error listening to port %s : %s", port, err)
	}

}
