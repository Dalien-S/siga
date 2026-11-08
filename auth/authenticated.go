package main

import (
	"log"
	"net/http"
	"uuid"

	"github.com/gboncoffee/siga/auth/repository"
)

func authenticated(w http.ResponseWriter, r *http.Request) {
	token, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	if err := repository.SessionExists(token); err != nil {
		log.Println(err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
