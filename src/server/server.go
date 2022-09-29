package main

import (
	c "backendForSharedProject/src/config"
	h "backendForSharedProject/src/handlers"
	j "backendForSharedProject/src/jwt"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mainConfig, err := c.LoadMainConfig("../config")
	if err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}
	fmt.Println("mainConfig", mainConfig)

	http.HandleFunc("/home", j.VerifyJWT(h.HandlePage))
	http.HandleFunc("/", h.DefaultHandler)
	http.HandleFunc("/auth", h.AuthPage)

	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Error listening to port %s : %s", port, err)
	}

}
