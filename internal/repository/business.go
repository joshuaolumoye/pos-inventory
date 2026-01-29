package repository

import (
	"strings"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/infrastructure"
	"gorm.io/gorm"
)

type BusinessRepo struct {
	DB *gorm.DB
}

func (r *BusinessRepo) CreateBusiness(b *domain.Business) error {
	infra := infrastructure.Business{
		ID:               b.ID,
		Name:             b.Name,
		OwnerFullName:    b.OwnerFullName,
		Email:            b.Email,
		PhoneNumber:      b.PhoneNumber,
		PasswordHash:     b.PasswordHash,
		StoreAddress:     b.StoreAddress,
		BusinessCategory: b.BusinessCategory,
		Currency:         b.Currency,
		StoreIcon:        b.StoreIcon,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
		Identifyer:       b.Identifyer,
	}
	err := r.DB.Create(&infra).Error
	if err == nil {
		b.ID = infra.ID
	}
	return err
}

func (r *BusinessRepo) GetBusinessByEmail(email string) (*domain.Business, error) {
	var infra infrastructure.Business
	err := r.DB.Where("LOWER(email) = ?", strings.ToLower(email)).First(&infra).Error
	if err != nil {
		return nil, err
	}
	b := &domain.Business{
		ID:               infra.ID,
		Name:             infra.Name,
		OwnerFullName:    infra.OwnerFullName,
		Email:            infra.Email,
		PhoneNumber:      infra.PhoneNumber,
		PasswordHash:     infra.PasswordHash,
		StoreAddress:     infra.StoreAddress,
		BusinessCategory: infra.BusinessCategory,
		Currency:         infra.Currency,
		StoreIcon:        infra.StoreIcon,
		CreatedAt:        infra.CreatedAt,
		UpdatedAt:        infra.UpdatedAt,
	}
	return b, nil
}

func (r *BusinessRepo) GetBusinessByID(id string) (*domain.Business, error) {
	var infra infrastructure.Business
	err := r.DB.First(&infra, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	b := &domain.Business{
		ID:               infra.ID,
		Name:             infra.Name,
		OwnerFullName:    infra.OwnerFullName,
		Email:            infra.Email,
		PhoneNumber:      infra.PhoneNumber,
		PasswordHash:     infra.PasswordHash,
		StoreAddress:     infra.StoreAddress,
		BusinessCategory: infra.BusinessCategory,
		Currency:         infra.Currency,
		StoreIcon:        infra.StoreIcon,
		CreatedAt:        infra.CreatedAt,
		UpdatedAt:        infra.UpdatedAt,
	}
	return b, nil
}

func (r *BusinessRepo) GetBusinessByIdentifyer(identifyer string) (*domain.Business, error) {
	var infra infrastructure.Business
	err := r.DB.Where("identifyer = ?", identifyer).First(&infra).Error
	if err != nil {
		return nil, err
	}
	b := &domain.Business{
		ID:               infra.ID,
		Name:             infra.Name,
		OwnerFullName:    infra.OwnerFullName,
		Email:            infra.Email,
		PhoneNumber:      infra.PhoneNumber,
		PasswordHash:     infra.PasswordHash,
		StoreAddress:     infra.StoreAddress,
		BusinessCategory: infra.BusinessCategory,
		Currency:         infra.Currency,
		StoreIcon:        infra.StoreIcon,
		CreatedAt:        infra.CreatedAt,
		UpdatedAt:        infra.UpdatedAt,
	}
	return b, nil
}
