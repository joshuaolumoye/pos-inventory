package repository

import (
	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/infrastructure"
	"gorm.io/gorm"
)

type BranchRepo struct {
	DB *gorm.DB
}

func (r *BranchRepo) CreateBranch(b *domain.Branch) error {
	infra := infrastructure.Branch{
		ID:            b.ID,
		BusinessID:    b.BusinessID,
		BranchName:    b.BranchName,
		BranchAddress: b.BranchAddress,
		IsMainBranch:  b.IsMainBranch,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
	err := r.DB.Create(&infra).Error
	if err == nil {
		b.ID = infra.ID
	}
	return err
}

func (r *BranchRepo) GetBranchByID(id string) (*domain.Branch, error) {
	var infra infrastructure.Branch
	err := r.DB.First(&infra, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	b := &domain.Branch{
		ID:            infra.ID,
		BusinessID:    infra.BusinessID,
		BranchName:    infra.BranchName,
		BranchAddress: infra.BranchAddress,
		IsMainBranch:  infra.IsMainBranch,
		CreatedAt:     infra.CreatedAt,
		UpdatedAt:     infra.UpdatedAt,
	}
	return b, nil
}

func (r *BranchRepo) GetBranchesByBusinessID(businessID string) ([]*domain.Branch, error) {
	var infras []*infrastructure.Branch
	err := r.DB.Where("business_id = ?", businessID).Find(&infras).Error
	if err != nil {
		return nil, err
	}
	var branches []*domain.Branch
	for _, infra := range infras {
		branches = append(branches, &domain.Branch{
			ID:            infra.ID,
			BusinessID:    infra.BusinessID,
			BranchName:    infra.BranchName,
			BranchAddress: infra.BranchAddress,
			IsMainBranch:  infra.IsMainBranch,
			CreatedAt:     infra.CreatedAt,
			UpdatedAt:     infra.UpdatedAt,
		})
	}
	return branches, nil
}

func (r *BranchRepo) DeleteBranch(id string) error {
	return r.DB.Delete(&infrastructure.Branch{}, "id = ?", id).Error
}

func (r *BranchRepo) UpdateBranch(b *domain.Branch) error {
	infra := infrastructure.Branch{
		ID:            b.ID,
		BusinessID:    b.BusinessID,
		BranchName:    b.BranchName,
		BranchAddress: b.BranchAddress,
		IsMainBranch:  b.IsMainBranch,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
	return r.DB.Save(&infra).Error
}
