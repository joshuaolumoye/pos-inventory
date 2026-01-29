package usecase

import (
	"fmt"
	"time"

	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type AuthRepo struct{}

func (a *AuthRepo) CreateToken(userID uint, role string) (string, string, error) {
	userIDStr := fmt.Sprintf("%d", userID)

	accessToken, err := utils.GenerateJWT(userIDStr, role, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateJWT(userIDStr, role, 7*24*time.Hour)
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
