package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// BearerToken represents a valid bearer token
type BearerToken struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// LoginHandler represents a dummy login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BearerToken{
		Token: tokenString,
		Type:  "Bearer",
	})
}
