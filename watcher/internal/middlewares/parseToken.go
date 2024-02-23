package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SigningKey = []byte(os.Getenv("KEY"))
)

// ParseToken parses a token
func ParseToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			http.Error(w, "missing Authorization header", http.StatusBadRequest)
			return
		}

		split := strings.Split(bearerToken, "Bearer ")
		tokenString := split[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return SigningKey, nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["username"] == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Internal server error: invalid claims", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
