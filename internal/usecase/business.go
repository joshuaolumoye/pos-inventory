package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"go.uber.org/zap"
)

type BusinessUsecase struct {
	BusinessRepo domain.BusinessRepository
	BranchRepo   domain.BranchRepository
}

func (u *BusinessUsecase) RegisterBusiness(b *domain.Business, plainPassword string) error {
	if b.Name == "" || b.OwnerFullName == "" || b.Email == "" || b.PhoneNumber == "" || plainPassword == "" || b.StoreAddress == "" || b.BusinessCategory == "" || b.Currency == "" {
		return errors.New("missing required fields")
	}

	b.Name = utils.Sanitize(b.Name)
	b.OwnerFullName = utils.Sanitize(b.OwnerFullName)
	b.Email = strings.ToLower(utils.Sanitize(b.Email))
	b.PhoneNumber = utils.Sanitize(b.PhoneNumber)
	b.StoreAddress = utils.Sanitize(b.StoreAddress)
	b.BusinessCategory = utils.Sanitize(b.BusinessCategory)
	b.Currency = utils.Sanitize(b.Currency)
	if b.StoreIcon != nil {
		icon := utils.Sanitize(*b.StoreIcon)
		b.StoreIcon = &icon
	}

	// Generate UUID for business ID
	b.ID = utils.GenerateUUID()

	// Generate unique identifyer if not set
	if b.Identifyer == "" {
		prefix := ""
		if len(b.Name) >= 5 {
			prefix = b.Name[:5]
		} else {
			prefix = b.Name
		}
		prefix = strings.ToLower(prefix)
		for {
			candidate := prefix + utils.RandomDigits(5)
			existing, _ := u.BusinessRepo.GetBusinessByIdentifyer(candidate)
			if existing == nil {
				b.Identifyer = candidate
				break
			}
		}
	}

	existing, _ := u.BusinessRepo.GetBusinessByEmail(b.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	hash, err := utils.HashPassword(plainPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	b.PasswordHash = hash

	now := time.Now().Unix()
	b.CreatedAt = now
	b.UpdatedAt = now

	if err := u.BusinessRepo.CreateBusiness(b); err != nil {
		return err
	}

	mainBranch := &domain.Branch{
		ID:            utils.GenerateUUID(),
		BusinessID:    b.ID,
		BranchName:    "Main Branch",
		BranchAddress: b.StoreAddress,
		IsMainBranch:  true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := u.BranchRepo.CreateBranch(mainBranch); err != nil {
		return err
	}

	return nil
}

func (u *BusinessUsecase) Login(email, password string, authRepo domain.AuthRepository) (string, string, string, string, error) {
	if email == "" || password == "" {
		utils.Logger.Error("Login failed: missing credentials")
		return "", "", "", "", errors.New("missing credentials")
	}

	email = strings.ToLower(strings.TrimSpace(utils.Sanitize(email)))
	utils.Logger.Info("Login attempt with sanitized email", zap.String("email", email))

	business, err := u.BusinessRepo.GetBusinessByEmail(email)
	if err != nil {
		utils.Logger.Error("Login failed: database error", zap.String("email", email), zap.Error(err))
		return "", "", "", "", errors.New("invalid credentials")
	}
	if business == nil {
		utils.Logger.Error("Login failed: business is nil", zap.String("email", email))
		return "", "", "", "", errors.New("invalid credentials")
	}

	utils.Logger.Info("Business found, checking password",
		zap.String("email", email),
		zap.String("businessID", business.ID))

	if !utils.CheckPasswordHash(password, business.PasswordHash) {
		utils.Logger.Error("Login failed: password mismatch", zap.String("email", email))
		return "", "", "", "", errors.New("invalid credentials")
	}

	utils.Logger.Info("Password verified, generating tokens", zap.String("email", email))

	role := "owner"
	access, refresh, err := authRepo.CreateToken(
		business.ID,
		role,
	)
	if err != nil {
		utils.Logger.Error("Token generation failed", zap.Error(err))
		return "", "", "", "", errors.New("token generation failed")
	}

	utils.Logger.Info("Login successful",
		zap.String("email", email),
		zap.String("businessID", business.ID))

	return access, refresh, business.ID, role, nil
}
