package repository

import (
	"time"

	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type AuthRepo struct{}

func (a *AuthRepo) CreateToken(userID string, role string) (string, string, error) {
	// Create access token (valid for 24 hours)
	accessToken, err := utils.GenerateJWT(userID, role, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	// Create refresh token (valid for 7 days)
	refreshToken, err := utils.GenerateJWT(userID, role, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *AuthRepo) ValidateToken(token string) (string, string, error) {
	claims, err := utils.ParseJWT(token)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}
