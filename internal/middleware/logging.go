package middleware

import (
	"net/http"

	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
