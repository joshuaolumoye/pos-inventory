package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/dto"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

var BranchUC *usecase.BranchUsecase

func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.BranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get business_id from context (set by AuthMiddleware)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	branch := &domain.Branch{
		ID:            utils.GenerateUUID(),
		BusinessID:    businessID,
		BranchName:    req.BranchName,
		BranchAddress: req.BranchAddress,
		IsMainBranch:  req.IsMainBranch,
	}
	err := BranchUC.CreateBranch(branch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := dto.BranchResponse{
		BranchID:      branch.ID,
		BusinessID:    branch.BusinessID,
		BranchName:    branch.BranchName,
		BranchAddress: branch.BranchAddress,
		IsMainBranch:  branch.IsMainBranch,
		CreatedAt:     branch.CreatedAt,
		UpdatedAt:     branch.UpdatedAt,
		Message:       "Branch created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetBranchesHandler(w http.ResponseWriter, r *http.Request) {
	// Get business_id from JWT context
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	search := r.URL.Query().Get("search")
	branches, err := BranchUC.GetBranchesByBusinessID(businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var branchResponses []dto.BranchResponse
	for _, branch := range branches {
		if search == "" || utils.FuzzyMatch(branch.BranchName, search) {
			branchResponses = append(branchResponses, dto.BranchResponse{
				BranchID:      branch.ID,
				BusinessID:    branch.BusinessID,
				BranchName:    branch.BranchName,
				BranchAddress: branch.BranchAddress,
				IsMainBranch:  branch.IsMainBranch,
				CreatedAt:     branch.CreatedAt,
				UpdatedAt:     branch.UpdatedAt,
			})
		}
	}

	resp := dto.BranchListResponse{
		Branches: branchResponses,
		Count:    len(branchResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func UpdateBranchHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.BranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get branch_id from query params
	branchID := r.URL.Query().Get("branch_id")
	if branchID == "" {
		http.Error(w, "branch_id is required", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context (NOT from query params!)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	branch := &domain.Branch{
		ID:            branchID,
		BusinessID:    businessID, // ✅ From JWT, not query params
		BranchName:    req.BranchName,
		BranchAddress: req.BranchAddress,
		IsMainBranch:  req.IsMainBranch,
	}

	err := BranchUC.UpdateBranch(branch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.BranchResponse{
		BranchID:      branch.ID,
		BusinessID:    branch.BusinessID,
		BranchName:    branch.BranchName,
		BranchAddress: branch.BranchAddress,
		IsMainBranch:  branch.IsMainBranch,
		CreatedAt:     branch.CreatedAt,
		UpdatedAt:     branch.UpdatedAt,
		Message:       "Branch updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteBranchHandler(w http.ResponseWriter, r *http.Request) {
	// Get branch_id from query params
	branchID := r.URL.Query().Get("branch_id")
	if branchID == "" {
		http.Error(w, "branch_id is required", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context (NOT from query params!)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	err := BranchUC.DeleteBranch(branchID, businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
