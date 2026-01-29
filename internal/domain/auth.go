package domain

type AuthRepository interface {
	CreateToken(userID string, role string) (accessToken string, refreshToken string, err error)
	ValidateToken(token string) (userID string, role string, err error)
}
