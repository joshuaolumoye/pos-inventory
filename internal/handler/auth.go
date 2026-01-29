package handler

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/joshuaolumoye/pos-backend/pkg/dto"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	// Validate the refresh token
	userID, role, err := AuthRepo.ValidateToken(req.RefreshToken)
	if err != nil {
		utils.Logger.Error("Invalid refresh token", zap.Error(err))
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	utils.Logger.Info("Refresh token validated", zap.String("userID", userID), zap.String("role", role))

	// Generate new tokens
	newAccessToken, newRefreshToken, err := AuthRepo.CreateToken(userID, role)
	if err != nil {
		utils.Logger.Error("Failed to generate new tokens", zap.Error(err))
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	// Convert userID from string to uint
	uid64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		utils.Logger.Error("Invalid userID format", zap.Error(err))
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	uid := uint(uid64)

	// Get business details
	business, err := BusinessUC.BusinessRepo.GetBusinessByID(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		utils.Logger.Error("Failed to fetch business", zap.Error(err))
		http.Error(w, "failed to fetch user details", http.StatusInternalServerError)
		return
	}

	businessIDUint := uint(0)
	if parsed, err := strconv.ParseUint(business.ID, 10, 64); err == nil {
		businessIDUint = uint(parsed)
	}
	resp := dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		BusinessID:   businessIDUint,
		Role:         role,
	}

	businessIDUint64, _ := strconv.ParseUint(business.ID, 10, 64)
	utils.Logger.Info("Tokens refreshed successfully", zap.String("businessID", strconv.FormatUint(businessIDUint64, 10)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
