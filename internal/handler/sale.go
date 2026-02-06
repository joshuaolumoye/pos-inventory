package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

var SaleUC *usecase.SaleUsecase

func CreateSaleHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info("=== CreateSaleHandler ENTERED ===")

	var req usecase.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error("Failed to decode request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "invalid input",
		})
		return
	}

	utils.Logger.Info("Request decoded successfully", zap.Any("request", req))

	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	utils.Logger.Info("Getting businessID from context", zap.Bool("ok", ok), zap.String("businessID", businessID))
	if !ok || businessID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "unauthorized: business ID not found",
		})
		return
	}

	cashierID, ok := middleware.GetUserIDFromContext(r.Context())
	utils.Logger.Info("Getting userID from context", zap.Bool("ok", ok), zap.String("userID", cashierID))
	if !ok || cashierID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "unauthorized: user ID not found",
		})
		return
	}

	role, ok := middleware.GetRoleFromContext(r.Context())
	utils.Logger.Info("Getting role from context", zap.Bool("ok", ok), zap.String("role", role))
	if !ok || role == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "unauthorized: role not found",
		})
		return
	}

	// All authenticated users can create sales - no role restriction
	utils.Logger.Info("Creating sale",
		zap.String("businessID", businessID),
		zap.String("cashierID", cashierID),
		zap.String("role", role))

	resp, err := SaleUC.CreateSale(&req, businessID, cashierID)
	if err != nil {
		msg := err.Error()
		status := http.StatusInternalServerError
		if msg == "unauthorized" {
			status = http.StatusUnauthorized
		} else if msg == "invalid input" {
			status = http.StatusBadRequest
		} else if msg == "quantity must be greater than 0" || msg == "insufficient stock" || msg == "product does not belong to business" || msg == "stock would become negative" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": msg,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
