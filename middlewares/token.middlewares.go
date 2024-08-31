package middlewares

import (
	"context"
	"net/http"
)

type token struct {
	Token string
}

func AddTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if access token exists.
		token, exists := r.Context().Value("access_token").(*token)
		if !exists {
			http.Error(w, "access token not found", http.StatusNotFound)
			return
		}

		//New context with access token.
		ctx := context.WithValue(r.Context(), "accessToken", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
