package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

type contextKey string

const (
	businessIDKey contextKey = "business_id"
	userIDKey     contextKey = "user_id"
	roleKey       contextKey = "role"
)

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
			return
		}

		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			utils.Logger.Warn("Invalid Authorization format - rejecting request")
			writeAuthError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[7:]
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.Logger.Warn("Token parsing failed", zap.Error(err))
			writeAuthError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		utils.Logger.Info("Token validated",
			zap.String("userID", claims.UserID),
			zap.String("role", claims.Role))

		// Set business_id, user_id, and role for all authenticated users
		ctx := r.Context()
		ctx = contextWithBusinessID(ctx, claims.UserID)
		ctx = contextWithUserID(ctx, claims.UserID)
		ctx = contextWithRole(ctx, claims.Role)

		utils.Logger.Info("Context values set",
			zap.String("businessID", claims.UserID),
			zap.String("userID", claims.UserID),
			zap.String("role", claims.Role))

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

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func contextWithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

// GetBusinessIDFromContext retrieves business_id from context
func GetBusinessIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(businessIDKey)
	if id, ok := val.(string); ok {
		return id, true
	}
	return "", false
}

// GetUserIDFromContext retrieves user_id from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(userIDKey)
	if id, ok := val.(string); ok {
		return id, true
	}
	return "", false
}

// GetRoleFromContext retrieves role from context
func GetRoleFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(roleKey)
	if role, ok := val.(string); ok {
		return role, true
	}
	return "", false
}

// NoQueryParamsMiddleware returns 404 if any query parameters are present
// Use this for POST/PUT endpoints that should only accept data in the request body
func NoQueryParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Info("NoQueryParamsMiddleware called",
			zap.String("path", r.URL.Path),
			zap.Int("query_params_count", len(r.URL.Query())),
			zap.Any("query_params", r.URL.Query()))

		if len(r.URL.Query()) > 0 {
			utils.Logger.Warn("Request has query parameters - rejecting with 404")
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		utils.Logger.Info("NoQueryParamsMiddleware passed - proceeding to handler")
		next.ServeHTTP(w, r)
	})
}
