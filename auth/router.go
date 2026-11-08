package main

import "net/http"

func router() {
	http.HandleFunc("POST /login", login)
	http.HandleFunc("POST /logout/{uuid}", logout)
	http.HandleFunc("POST /credentials", credentials)
	http.HandleFunc("GET /{uuid}", authenticated)
}
