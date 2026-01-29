package usecase

import (
	"errors"
	"time"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type BranchUsecase struct {
	BranchRepo domain.BranchRepository
}

func (u *BranchUsecase) CreateBranch(b *domain.Branch) error {
	if b.BranchName == "" || b.BranchAddress == "" || b.BusinessID == "" {
		return errors.New("missing required fields")
	}
	b.BranchName = utils.Sanitize(b.BranchName)
	b.BranchAddress = utils.Sanitize(b.BranchAddress)
	now := time.Now().Unix()
	b.CreatedAt = now
	b.UpdatedAt = now
	return u.BranchRepo.CreateBranch(b)
}

func (u *BranchUsecase) UpdateBranch(b *domain.Branch) error {
	if b.ID == "" || b.BusinessID == "" {
		return errors.New("missing branch_id or business_id")
	}
	existing, err := u.BranchRepo.GetBranchByID(b.ID)
	if err != nil {
		return errors.New("branch not found")
	}
	if existing.BusinessID != b.BusinessID {
		return errors.New("unauthorized")
	}
	b.BranchName = utils.Sanitize(b.BranchName)
	b.BranchAddress = utils.Sanitize(b.BranchAddress)
	b.UpdatedAt = time.Now().Unix()
	b.CreatedAt = existing.CreatedAt
	return u.BranchRepo.UpdateBranch(b)
}

func (u *BranchUsecase) DeleteBranch(branchID, businessID string) error {
	if branchID == "" || businessID == "" {
		return errors.New("missing branch_id or business_id")
	}
	branch, err := u.BranchRepo.GetBranchByID(branchID)
	if err != nil {
		return errors.New("branch not found")
	}
	if branch.BusinessID != businessID {
		return errors.New("unauthorized")
	}
	if branch.IsMainBranch {
		return errors.New("cannot delete main branch")
	}
	return u.BranchRepo.DeleteBranch(branchID)
}

func (u *BranchUsecase) GetBranchesByBusinessID(businessID string) ([]*domain.Branch, error) {
	if businessID == "" {
		return nil, errors.New("missing business_id")
	}
	return u.BranchRepo.GetBranchesByBusinessID(businessID)
}
