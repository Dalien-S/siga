package main

import (
	"crypto/sha512"
	"encoding/json"
	"log"
	"net/http"
	"uuid"

	"github.com/gboncoffee/siga/auth/repository"
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token uuid.UUID
}

func login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	c, err := repository.GetByUsername(req.Username)
	if err != nil {
		log.Println("here")
		log.Println(err)
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if c.Password != sha512.Sum512([]byte(req.Password)) {
		log.Println("there")
		log.Println(err)
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token := uuid.New()
	if err := repository.NewSession(token); err != nil {
		http.Error(w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
