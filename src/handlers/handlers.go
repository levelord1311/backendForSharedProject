package handlers

import (
	j "backendForSharedProject/src/jwt"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(message)
	if err != nil {
		return
	}
}

func AuthPage(w http.ResponseWriter, r *http.Request) {
	token, err := j.GenerateJWT()
	if err != nil {
		return
	}
	client := &http.Client{}
	r, _ = http.NewRequest("POST", "<http://localhost:8080/>", nil)
	r.Header.Set("Token", token)
	if _, err := client.Do(r); err != nil {
		return
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "default handler, serving %s\n", r.Host)
}
