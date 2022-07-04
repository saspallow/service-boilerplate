package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
)

func JwtPayLoad(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "SOME_FIELD", claims["SOME_FIELD"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
