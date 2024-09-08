package middlewares

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type contextKey string

var TokenKey contextKey = "accessToken"

func AddTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if access token exists.
		token, exists := r.Context().Value("access_token").(string)
		if !exists {
			logrus.Fatal("Access token was not found: ")
			return
		}

		//New context with access token.
		ctx := context.WithValue(r.Context(), TokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
