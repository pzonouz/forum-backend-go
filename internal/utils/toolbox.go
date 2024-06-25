package utils

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func ReadJSON[T any](w http.ResponseWriter, r *http.Request) T {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	decoder := json.NewDecoder(r.Body)
	data := new(T)
	err := decoder.Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return *data
}

func WriteJSON(w http.ResponseWriter, data any) {
	JSON := json.NewEncoder(w)
	err := JSON.Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type MyClaims struct {
	ID          int64
	Email       string
	Address     string
	PhoneNumber string
	Role        string
	Expired     int64
	jwt.RegisteredClaims
}
