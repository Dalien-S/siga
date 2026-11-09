package main

import (
	"net/http"
	"uuid"

	"github.com/gboncoffee/siga/auth/repository"
)

func logout(w http.ResponseWriter, r *http.Request) {
	token, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	repository.DeleteSession(token)
	w.WriteHeader(http.StatusOK)
}
