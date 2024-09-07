package middlewares

import (
	"context"
	"net/http"
)

type contextKey string

var tokenKey contextKey = "accessToken"

func AddTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if access token exists.
		token, exists := r.Context().Value("access_token").(string)
		if !exists {
			http.Error(w, "access token not found", http.StatusNotFound)
			return
		}

		//New context with access token.
		ctx := context.WithValue(r.Context(), tokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
