package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type StaffUsecase struct {
	StaffRepo domain.StaffRepository
}

func (u *StaffUsecase) CreateStaff(s *domain.Staff, password string) error {
	if s.FullName == "" || s.StaffID == "" || s.PhoneNumber == "" || password == "" || s.Role == "" || s.BranchID == "" || s.BusinessID == "" || s.Status == "" {
		return errors.New("missing required fields")
	}

	s.FullName = utils.Sanitize(s.FullName)
	s.StaffID = strings.ToLower(utils.Sanitize(s.StaffID)) // Assuming StaffID is now a string
	s.PhoneNumber = utils.Sanitize(s.PhoneNumber)
	s.Status = utils.Sanitize(s.Status)

	if s.PhotoURL != nil {
		url := utils.Sanitize(*s.PhotoURL)
		s.PhotoURL = &url
	}

	switch s.Role {
	case domain.RoleOwner, domain.RoleManager, domain.RoleCashier, domain.RoleInventory:
		// valid
	default:
		return errors.New("invalid staff role")
	}

	existing, _ := u.StaffRepo.GetStaffByStaffID(s.StaffID)
	if existing != nil {
		return errors.New("staff_id already exists")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	s.PasswordHash = hash

	now := time.Now().Unix()
	s.CreatedAt = now
	s.UpdatedAt = now

	return u.StaffRepo.CreateStaff(s) // Assuming CreateStaff method is updated to accept string IDs
}

func (u *StaffUsecase) Login(staffID, password string, authRepo domain.AuthRepository) (string, string, string, string, error) {
	if staffID == "" || password == "" {
		return "", "", "", "", errors.New("missing credentials")
	}

	staffID = strings.ToLower(strings.TrimSpace(utils.Sanitize(staffID)))

	staff, err := u.StaffRepo.GetStaffByStaffID(staffID)
	if err != nil || staff == nil {
		return "", "", "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, staff.PasswordHash) {
		return "", "", "", "", errors.New("invalid credentials")
	}

	access, refresh, err := authRepo.CreateToken(staff.ID, string(staff.Role))
	if err != nil {
		return "", "", "", "", errors.New("token generation failed")
	}

	return access, refresh, staff.ID, string(staff.Role), nil
}

func (u *StaffUsecase) GetStaffByID(staffID string) (*domain.Staff, error) {
	if staffID == "" {
		return nil, errors.New("missing staff_id")
	}
	return u.StaffRepo.GetStaffByStaffID(staffID) // Assuming GetStaffByStaffID method is updated to accept string IDs
}

func (u *StaffUsecase) GetStaffByBusinessID(businessID string) ([]*domain.Staff, error) {
	if businessID == "" {
		return nil, errors.New("missing business_id")
	}
	return u.StaffRepo.GetStaffByBusinessID(businessID)
}

func (u *StaffUsecase) GetStaffByBranchID(businessID, branchID string) ([]*domain.Staff, error) {
	if businessID == "" || branchID == "" {
		return nil, errors.New("missing business_id or branch_id")
	}
	return u.StaffRepo.GetStaffByBranchID(businessID, branchID)
}

func (u *StaffUsecase) UpdateStaff(s *domain.Staff) error {
	if s.ID == "" || s.BusinessID == "" {
		return errors.New("missing staff id or business_id")
	}

	existing, err := u.StaffRepo.GetStaffByID(s.ID)
	if err != nil {
		return errors.New("staff not found")
	}

	if existing.BusinessID != s.BusinessID {
		return errors.New("unauthorized")
	}

	s.FullName = utils.Sanitize(s.FullName)
	s.PhoneNumber = utils.Sanitize(s.PhoneNumber)
	s.Status = utils.Sanitize(s.Status)

	if s.PhotoURL != nil {
		url := utils.Sanitize(*s.PhotoURL)
		s.PhotoURL = &url
	}

	switch s.Role {
	case domain.RoleOwner, domain.RoleManager, domain.RoleCashier, domain.RoleInventory:
		// valid
	default:
		return fmt.Errorf("invalid staff role")
	}

	s.UpdatedAt = time.Now().Unix()
	s.CreatedAt = existing.CreatedAt
	s.StaffID = existing.StaffID
	s.PasswordHash = existing.PasswordHash

	return u.StaffRepo.UpdateStaff(s) // Assuming UpdateStaff method is updated to accept string IDs
}
