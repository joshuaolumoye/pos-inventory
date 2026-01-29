package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/dto"
)

// ProductUC is injected by main.go
var ProductUC *usecase.ProductUsecase

// AddProductHandler handles product creation
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context (set by AuthMiddleware)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Extract staff_id from JWT if available (for staff users)
	// In this codebase, businessID in context is set to claims.UserID, which for staff is their staff ID, for owner is business ID
	// So we need to determine if this is a staff or owner
	// Let's try to fetch staff by this ID; if found, use their businessID, else treat as owner
	staff, _ := StaffUC.StaffRepo.GetStaffByID(businessID)
	var realBusinessID, createdBy string
	var updatedBy *string
	if staff != nil {
		// Staff user
		realBusinessID = staff.BusinessID
		createdBy = staff.ID
		updatedBy = &staff.ID
		// Only allow certain roles
		switch staff.Role {
		case domain.RoleOwner, domain.RoleManager, domain.RoleCashier, domain.RoleInventory:
			// allowed
		default:
			http.Error(w, "unauthorized staff role", http.StatusForbidden)
			return
		}
	} else {
		// Owner user
		realBusinessID = businessID
		createdBy = businessID
		updatedBy = &businessID
	}

	var barcodePtr *string
	if req.BarcodeValue != "" {
		barcodePtr = &req.BarcodeValue
	} else {
		barcodePtr = nil
	}

	product := &domain.Product{
		ProductName:       req.ProductName,
		ProductCategory:   req.ProductCategory,
		SellingPrice:      req.SellingPrice,
		CostPrice:         req.CostPrice,
		QuantityInStock:   req.QuantityInStock,
		LowStockThreshold: req.LowStockThreshold,
		BarcodeValue:      barcodePtr,
		NAFDACRegNumber:   func() *string { v := req.NAFDACRegNumber; return &v }(),
		ExpiryDate:        req.ExpiryDate,
		ProductImageURL:   req.ProductImageURL,
		BranchID:          req.BranchID,
		BusinessID:        realBusinessID, // Always from context or staff
		CreatedBy:         createdBy,      // staff ID or business ID
		UpdatedBy:         updatedBy,      // staff ID or business ID
	}

	err := ProductUC.AddProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		NAFDACRegNumber: func() string {
			if product.NAFDACRegNumber != nil {
				return *product.NAFDACRegNumber
			}
			return ""
		}(),
		SellingPrice:    product.SellingPrice,
		QuantityLeft:    product.QuantityInStock,
		ProductImageURL: product.ProductImageURL,
		Message:         "Product added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateProductHandler handles product update by product_id in URL
// Route: /api/product/{id}
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get product_id from chi URL param
	productID := chi.URLParam(r, "id")
	if productID == "" {
		http.Error(w, "product id is required in URL", http.StatusBadRequest)
		return
	}

	// Get business_id and user_id from JWT context (NOT from query params!)
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Updated_by should be the user_id from JWT
	updatedBy := &businessID // Using businessID as updater for now

	product := &domain.Product{
		ID:                productID,
		ProductName:       req.ProductName,
		ProductCategory:   req.ProductCategory,
		SellingPrice:      req.SellingPrice,
		CostPrice:         req.CostPrice,
		QuantityInStock:   req.QuantityInStock,
		LowStockThreshold: req.LowStockThreshold,
		BarcodeValue:      func() *string { v := req.BarcodeValue; return &v }(),
		NAFDACRegNumber:   func() *string { v := req.NAFDACRegNumber; return &v }(),
		ExpiryDate:        req.ExpiryDate,
		ProductImageURL:   req.ProductImageURL,
		BranchID:          req.BranchID,
		BusinessID:        businessID, // ✅ From JWT
		UpdatedBy:         updatedBy,  // ✅ From JWT
	}

	err := ProductUC.UpdateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		NAFDACRegNumber: func() string {
			if product.NAFDACRegNumber != nil {
				return *product.NAFDACRegNumber
			}
			return ""
		}(),
		SellingPrice:    product.SellingPrice,
		QuantityLeft:    product.QuantityInStock,
		ProductImageURL: product.ProductImageURL,
		Message:         "Product updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetProductHandler retrieves a single product by ID from chi URL param
// Route: /api/product/{id}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		http.Error(w, "product id is required in URL", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT to ensure user can only view their own business's products
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	product, err := ProductUC.GetProductByID(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify product belongs to the authenticated business
	if product.BusinessID != businessID {
		http.Error(w, "unauthorized access", http.StatusForbidden)
		return
	}

	resp := dto.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		NAFDACRegNumber: func() string {
			if product.NAFDACRegNumber != nil {
				return *product.NAFDACRegNumber
			}
			return ""
		}(),
		SellingPrice:    product.SellingPrice,
		QuantityLeft:    product.QuantityInStock,
		ProductImageURL: product.ProductImageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetProductsHandler retrieves all products for the authenticated business
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Get business_id from JWT context
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Optional: filter by branch_id
	branchID := r.URL.Query().Get("branch_id")

	var (
		products []*domain.Product
		err      error
	)

	if branchID != "" {
		products, err = ProductUC.GetProductsByBranchID(businessID, branchID)
	} else {
		products, err = ProductUC.GetProductsByBusinessID(businessID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var productResponses []dto.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, dto.ProductResponse{
			ID:          p.ID,
			ProductName: p.ProductName,
			NAFDACRegNumber: func() string {
				if p.NAFDACRegNumber != nil {
					return *p.NAFDACRegNumber
				}
				return ""
			}(),
			SellingPrice:    p.SellingPrice,
			QuantityLeft:    p.QuantityInStock,
			ProductImageURL: p.ProductImageURL,
		})
	}

	resp := dto.ProductListResponse{
		Products: productResponses,
		Count:    len(productResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteProductHandler handles product deletion (soft delete)
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	// Get product_id from query params
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "product_id is required", http.StatusBadRequest)
		return
	}

	// Get business_id from JWT context
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	err := ProductUC.DeleteProduct(productID, businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
