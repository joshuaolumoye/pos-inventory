package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

type contextKey string

const businessIDKey contextKey = "business_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Info("AuthMiddleware called",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))

		authHeader := r.Header.Get("Authorization")
		utils.Logger.Info("Authorization header", zap.String("value", authHeader))

		if authHeader == "" {
			utils.Logger.Warn("Missing Authorization header - rejecting request")
			writeAuthError(w, "Missing Authorization header", http.StatusUnauthorized)
			return // ← Make sure this return is here!
		}

		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			utils.Logger.Warn("Invalid Authorization format - rejecting request")
			writeAuthError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return // ← Make sure this return is here!
		}

		tokenStr := authHeader[7:]
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.Logger.Warn("Token parsing failed", zap.Error(err))
			writeAuthError(w, "Invalid or expired token", http.StatusUnauthorized)
			return // ← Make sure this return is here!
		}

		utils.Logger.Info("Token validated",
			zap.String("userID", claims.UserID),
			zap.String("role", claims.Role))

		// Set business_id for owner, manager, cashier, inventory_staff
		var ctx context.Context
		switch claims.Role {
		case "owner", "manager", "cashier", "inventory_staff":
			ctx = contextWithBusinessID(r.Context(), claims.UserID)
			utils.Logger.Info("BusinessID set in context", zap.String("businessID", claims.UserID))
		default:
			utils.Logger.Warn("Unauthorized role", zap.String("role", claims.Role))
			writeAuthError(w, "Unauthorized role", http.StatusForbidden)
			return // ← Make sure this return is here!
		}

		r = r.WithContext(ctx)
		utils.Logger.Info("Proceeding to handler")
		next.ServeHTTP(w, r)
	})
}

// writeAuthError writes a consistent JSON error response for auth failures
func writeAuthError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := map[string]interface{}{
		"error":   true,
		"message": message,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func contextWithBusinessID(ctx context.Context, businessID string) context.Context {
	return context.WithValue(ctx, businessIDKey, businessID)
}

// GetBusinessIDFromContext retrieves business_id from context
func GetBusinessIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(businessIDKey)
	if id, ok := val.(string); ok {
		return id, true
	}
	return "", false
}

// NoQueryParamsMiddleware returns 404 if any query parameters are present
// Use this for POST/PUT endpoints that should only accept data in the request body
func NoQueryParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()) > 0 {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
