package main

import (
	"crypto/sha512"
	"encoding/json"
	"net/http"

	"github.com/gboncoffee/siga/auth/repository"
)

type CredentialsRequest struct {
	Username string
	Password string
}

func credentials(w http.ResponseWriter, r *http.Request) {
	var req CredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	c := repository.Credential{
		Username: req.Username,
		Password: sha512.Sum512([]byte(req.Password)),
	}
	if err := repository.NewCredential(c); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
