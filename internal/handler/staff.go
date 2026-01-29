package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/dto"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

var StaffUC *usecase.StaffUsecase

func CreateStaffHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.StaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context (NOT from request!)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Check that branch exists and belongs to this business
	branch, err := BranchUC.BranchRepo.GetBranchByID(req.BranchID)
	if err != nil || branch == nil || branch.BusinessID != businessID {
		http.Error(w, "invalid branch_id for this business", http.StatusBadRequest)
		return
	}

	staff := &domain.Staff{
		StaffID:     req.StaffID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Role:        domain.StaffRole(req.Role),
		BranchID:    req.BranchID,
		BusinessID:  businessID, // ✅ From JWT, not request
		Status:      req.Status,
		PhotoURL:    req.PhotoURL,
	}

	err = StaffUC.CreateStaff(staff, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.StaffResponse{
		ID:          staff.ID,
		StaffID:     staff.StaffID,
		FullName:    staff.FullName,
		PhoneNumber: staff.PhoneNumber,
		Role:        string(staff.Role),
		BranchID:    staff.BranchID,
		BusinessID:  staff.BusinessID,
		Status:      staff.Status,
		PhotoURL:    staff.PhotoURL,
		CreatedAt:   staff.CreatedAt,
		UpdatedAt:   staff.UpdatedAt,
		Message:     "Staff created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func StaffLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.StaffLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	businessIdentifyer := strings.TrimSpace(req.BusinessIdentifyer)
	staffID := strings.ToLower(strings.TrimSpace(req.StaffID))
	password := req.Password

	if businessIdentifyer == "" || staffID == "" || password == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	// Fetch business by identifyer
	business, err := BusinessUC.BusinessRepo.GetBusinessByIdentifyer(businessIdentifyer)
	if err != nil || business == nil {
		http.Error(w, "invalid business identifier", http.StatusUnauthorized)
		return
	}

	// Fetch staff by staffID and businessID
	staff, err := StaffUC.StaffRepo.GetStaffByStaffID(staffID)
	if err != nil || staff == nil || staff.BusinessID != business.ID {
		http.Error(w, "invalid staff credentials", http.StatusUnauthorized)
		return
	}

	// Check password securely
	if !utils.CheckPasswordHash(password, staff.PasswordHash) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	token, refreshToken, err := AuthRepo.CreateToken(staff.ID, string(staff.Role))
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	resp := dto.StaffLoginResponse{
		StaffID:      staffID,
		StaffName:    staff.FullName,
		Role:         string(staff.Role),
		Token:        token,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetStaffByIDHandler(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staff_id")
	if staffID == "" {
		http.Error(w, "missing staff_id", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT to ensure user can only view their own business's staff
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	staff, err := StaffUC.GetStaffByID(staffID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify staff belongs to the authenticated business
	if staff.BusinessID != businessID {
		http.Error(w, "unauthorized access", http.StatusForbidden)
		return
	}

	resp := dto.StaffResponse{
		ID:          staff.ID,
		StaffID:     staff.StaffID,
		FullName:    staff.FullName,
		PhoneNumber: staff.PhoneNumber,
		Role:        string(staff.Role),
		BranchID:    staff.BranchID,
		BusinessID:  staff.BusinessID,
		Status:      staff.Status,
		PhotoURL:    staff.PhotoURL,
		CreatedAt:   staff.CreatedAt,
		UpdatedAt:   staff.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetStaffListHandler(w http.ResponseWriter, r *http.Request) {
	// Get business_id from JWT context
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Optional: filter by branch_id
	branchID := r.URL.Query().Get("branch_id")

	var (
		staff []*domain.Staff
		err   error
	)
	if branchID != "" {
		staff, err = StaffUC.GetStaffByBranchID(businessID, branchID)
	} else {
		staff, err = StaffUC.GetStaffByBusinessID(businessID)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var staffResponses []dto.StaffResponse
	for _, s := range staff {
		staffResponses = append(staffResponses, dto.StaffResponse{
			ID:          s.ID,
			StaffID:     s.StaffID,
			FullName:    s.FullName,
			PhoneNumber: s.PhoneNumber,
			Role:        string(s.Role),
			BranchID:    s.BranchID,
			BusinessID:  s.BusinessID,
			Status:      s.Status,
			PhotoURL:    s.PhotoURL,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	resp := dto.StaffListResponse{
		Staff: staffResponses,
		Count: len(staffResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func UpdateStaffHandler(w http.ResponseWriter, r *http.Request) {
	// Get staff ID from query params
	staffID := r.URL.Query().Get("id")
	if staffID == "" {
		http.Error(w, "missing staff id", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context (NOT from query params!)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	var req dto.StaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	staff := &domain.Staff{
		ID:          staffID,
		BusinessID:  businessID, // ✅ From JWT, not request
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Role:        domain.StaffRole(req.Role),
		BranchID:    req.BranchID,
		Status:      req.Status,
		PhotoURL:    req.PhotoURL,
	}

	err := StaffUC.UpdateStaff(staff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch updated staff to return complete data
	updatedStaff, err := StaffUC.StaffRepo.GetStaffByID(staffID)
	if err != nil {
		http.Error(w, "failed to fetch updated staff", http.StatusInternalServerError)
		return
	}

	resp := dto.StaffResponse{
		ID:          updatedStaff.ID,
		StaffID:     updatedStaff.StaffID,
		FullName:    updatedStaff.FullName,
		PhoneNumber: updatedStaff.PhoneNumber,
		Role:        string(updatedStaff.Role),
		BranchID:    updatedStaff.BranchID,
		BusinessID:  updatedStaff.BusinessID,
		Status:      updatedStaff.Status,
		PhotoURL:    updatedStaff.PhotoURL,
		CreatedAt:   updatedStaff.CreatedAt,
		UpdatedAt:   updatedStaff.UpdatedAt,
		Message:     "Staff updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
