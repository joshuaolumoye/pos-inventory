package repository

import (
	"strings"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/infrastructure"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
	"gorm.io/gorm"
)

type StaffRepo struct {
	DB *gorm.DB
}

func (r *StaffRepo) CreateStaff(s *domain.Staff) error {
	infra := infrastructure.Staff{
		ID:           s.ID,
		StaffID:      s.StaffID,
		FullName:     s.FullName,
		PhoneNumber:  s.PhoneNumber,
		PasswordHash: s.PasswordHash,
		Role:         infrastructure.StaffRole(s.Role),
		Status:       s.Status,
		PhotoURL:     s.PhotoURL,
		BranchID:     s.BranchID,
		BusinessID:   s.BusinessID,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
	if s.ID == "" {
		newID := utils.GenerateUUID()
		s.ID = newID
		infra.ID = newID
	}
	err := r.DB.Create(&infra).Error
	if err == nil {
		s.ID = infra.ID
	}
	return err
}

func (r *StaffRepo) GetStaffByStaffID(staffID string) (*domain.Staff, error) {
	var infra infrastructure.Staff
	err := r.DB.Where("LOWER(staff_id) = ?", strings.ToLower(staffID)).First(&infra).Error
	if err != nil {
		return nil, err
	}
	s := &domain.Staff{
		ID:           infra.ID,
		StaffID:      infra.StaffID,
		FullName:     infra.FullName,
		PhoneNumber:  infra.PhoneNumber,
		PasswordHash: infra.PasswordHash,
		Role:         domain.StaffRole(infra.Role),
		Status:       infra.Status,
		PhotoURL:     infra.PhotoURL,
		CreatedAt:    infra.CreatedAt,
		UpdatedAt:    infra.UpdatedAt,
		BranchID:     infra.BranchID,
		BusinessID:   infra.BusinessID,
	}
	return s, nil
}

func (r *StaffRepo) GetStaffByID(id string) (*domain.Staff, error) {
	var infra infrastructure.Staff
	err := r.DB.First(&infra, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	s := &domain.Staff{
		ID:           infra.ID,
		StaffID:      infra.StaffID,
		FullName:     infra.FullName,
		PhoneNumber:  infra.PhoneNumber,
		PasswordHash: infra.PasswordHash,
		Role:         domain.StaffRole(infra.Role),
		Status:       infra.Status,
		PhotoURL:     infra.PhotoURL,
		CreatedAt:    infra.CreatedAt,
		UpdatedAt:    infra.UpdatedAt,
		BranchID:     infra.BranchID,
		BusinessID:   infra.BusinessID,
	}
	return s, nil
}

func (r *StaffRepo) GetStaffByBusinessID(businessID string) ([]*domain.Staff, error) {
	var infras []*infrastructure.Staff
	// Only fetch staff that are not soft-deleted (DeletedAt is null or 0)
	err := r.DB.Where("business_id = ?", businessID).Find(&infras).Error
	if err != nil {
		return nil, err
	}
	var staffList []*domain.Staff
	for _, infra := range infras {
		staffList = append(staffList, &domain.Staff{
			ID:           infra.ID,
			StaffID:      infra.StaffID,
			FullName:     infra.FullName,
			PhoneNumber:  infra.PhoneNumber,
			PasswordHash: infra.PasswordHash,
			Role:         domain.StaffRole(infra.Role),
			Status:       infra.Status,
			PhotoURL:     infra.PhotoURL,
			BranchID:     infra.BranchID,
			BusinessID:   infra.BusinessID,
		})
	}
	return staffList, nil
}

func (r *StaffRepo) GetStaffByBranchID(businessID, branchID string) ([]*domain.Staff, error) {
	var infras []*infrastructure.Staff
	// Only fetch staff that are not soft-deleted (DeletedAt is null or 0)
	err := r.DB.Where("business_id = ? AND branch_id = ?", businessID, branchID).Find(&infras).Error
	if err != nil {
		return nil, err
	}
	var staffList []*domain.Staff
	for _, infra := range infras {
		staffList = append(staffList, &domain.Staff{
			ID:           infra.ID,
			StaffID:      infra.StaffID,
			FullName:     infra.FullName,
			PhoneNumber:  infra.PhoneNumber,
			PasswordHash: infra.PasswordHash,
			Role:         domain.StaffRole(infra.Role),
			Status:       infra.Status,
			PhotoURL:     infra.PhotoURL,
			BranchID:     infra.BranchID,
			BusinessID:   infra.BusinessID,
		})
	}
	return staffList, nil
}

func (r *StaffRepo) UpdateStaff(s *domain.Staff) error {
	infra := infrastructure.Staff{
		ID:           s.ID,
		StaffID:      s.StaffID,
		FullName:     s.FullName,
		PhoneNumber:  s.PhoneNumber,
		PasswordHash: s.PasswordHash,
		Role:         infrastructure.StaffRole(s.Role),
		Status:       s.Status,
		PhotoURL:     s.PhotoURL,
		BranchID:     s.BranchID,
		BusinessID:   s.BusinessID,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
	return r.DB.Save(&infra).Error
}
