package main

import (
	h "backendForSharedProject/src/handlers"
	j "backendForSharedProject/src/jwt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/home", j.VerifyJWT(h.HandlePage))
	http.HandleFunc("/", h.DefaultHandler)
	http.HandleFunc("/auth", h.AuthPage)

	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Error listening to port %s : %s", port, err)
	}

}
