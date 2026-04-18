package api

import (
	"context"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("reelix_auth_token")

		if err != nil {
			http.Error(w, "Unauthorized 1", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		log.Printf("auth token %v", tokenStr)

		claims, err := validateJWT(tokenStr)
		if err != nil {
			log.Printf("err %v", err.Error())
			http.Error(w, "Unauthorized 2", http.StatusUnauthorized)
			return
		}

		// store user info in context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
