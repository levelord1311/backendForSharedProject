package main

import (
	"fmt"
	"net/http"
)

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request. Note that "127.0.0.1:443" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, "https://127.0.0.1:443"+r.RequestURI, http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
	fmt.Println(r.RequestURI)
}

func main() {
	http.HandleFunc("/", handler)
	// Start the HTTPS server in a goroutine
	go http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
