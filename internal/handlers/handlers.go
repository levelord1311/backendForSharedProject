package handlers

import (
	j "backendForSharedProject/internal/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
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

func RedirectToTls(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	u := r.URL
	u.Host = net.JoinHostPort(host, ":443")
	u.Scheme = "https"
	log.Println(u.String())
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

func AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	var a Authn
	err := DecodeJSONBody(w, r, &a)
	if err != nil {
		var mr *MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "cheking login: %+v\n", a.Login)

	// encrypt password
	encrPass, err := HashPassword(a.Password)
	if err != nil {
		err := "Password Encryption  failed"
		fmt.Println(err)
	}
	fmt.Fprintf(w, "hashed pw: %+v\n", encrPass)

}
