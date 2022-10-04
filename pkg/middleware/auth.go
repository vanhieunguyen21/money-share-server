package middleware

import (
	"fmt"
	"money_share/pkg/auth"
	"net/http"
)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot validate token: %s", err), http.StatusUnauthorized)
			return
		}
		r.Header.Set("username", claims.Username)
		h.ServeHTTP(w, r)
	})
}
