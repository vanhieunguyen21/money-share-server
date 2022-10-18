package middleware

import (
	"fmt"
	"money_share/pkg/auth"
	"money_share/pkg/controller"
	"net/http"
	"strconv"
)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			controller.ResponseError(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateAccessToken(tokenStr)
		if err != nil {
			controller.ResponseError(w, fmt.Sprintf("Cannot validate token: %s", err), http.StatusUnauthorized)
			return
		}
		r.Header.Set("userID", strconv.Itoa(int(claims.UserID)))
		r.Header.Set("username", claims.Username)
		h.ServeHTTP(w, r)
	})
}
