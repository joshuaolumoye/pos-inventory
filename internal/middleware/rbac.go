package middleware

import (
	"net/http"
)

func RBACMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Check user role from context
			next.ServeHTTP(w, r)
		})
	}
}
