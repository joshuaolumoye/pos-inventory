package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/dto"
)

// ProductUC is injected by main.go
var ProductUC *usecase.ProductUsecase
var NotificationUC *usecase.NotificationUsecase

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

	// Check for unique product name in branch
	existingProducts, err := ProductUC.GetProductsByBranchID(realBusinessID, req.BranchID)
	if err == nil {
		for _, p := range existingProducts {
			if p.ProductName == req.ProductName {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{"error": "Product name already exists in this branch"})
				return
			}
		}
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

	err = ProductUC.AddProduct(product)
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

// NotificationProductResponse represents the notification product fields
type NotificationProductResponse struct {
	ProductName  string  `json:"product_name"`
	Barcode      *string `json:"barcode"`
	SellingPrice float64 `json:"selling_price"`
	QuantityLeft int     `json:"quantity_left"`
	ExpiryDate   *int64  `json:"expiring_date"`
}

// GetProductNotificationsHandler handles notification queries for products
// Route: /api/notifications/products?type=in_stock|low_stock|expired&page=1&per_page=20
func GetProductNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	notifType := r.URL.Query().Get("type")
	search := r.URL.Query().Get("search")
	page := 1
	perPage := 20
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			perPage = n
		}
	}
	offset := (page - 1) * perPage

	var products []*domain.Product
	var err error
	now := time.Now().Unix()

	switch {
	case search != "":
		// Paginated fuzzy search by product name
		products, err = ProductUC.ProductRepo.SearchProductsPaginated(businessID, search, perPage, offset)
	case notifType == "in_stock":
		products, err = ProductUC.ProductRepo.QueryProductsNotification(businessID, ">", 10, 0, 0, perPage, offset, false)
	case notifType == "low_stock":
		products, err = ProductUC.ProductRepo.QueryProductsNotification(businessID, "<", 10, 0, 0, perPage, offset, false)
	case notifType == "expired":
		products, err = ProductUC.ProductRepo.QueryProductsNotification(businessID, "", 0, now, 0, perPage, offset, true)
	case notifType == "":
		// No filter: return all products paginated
		products, err = ProductUC.ProductRepo.GetAllProductsPaginated(businessID, perPage, offset)
	default:
		http.Error(w, "invalid type param (must be in_stock, low_stock, expired, or empty)", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []NotificationProductResponse
	for _, p := range products {
		resp = append(resp, NotificationProductResponse{
			ProductName:  p.ProductName,
			Barcode:      p.BarcodeValue,
			SellingPrice: p.SellingPrice,
			QuantityLeft: p.QuantityInStock,
			ExpiryDate:   p.ExpiryDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": resp,
		"page":     page,
		"per_page": perPage,
		"count":    len(resp),
	})
}

// NotificationListResponse for notification API
type NotificationListResponse struct {
	Notifications []*domain.Notification `json:"notifications"`
	Page          int                    `json:"page"`
	PerPage       int                    `json:"per_page"`
	Count         int                    `json:"count"`
}

// ListNotificationsHandler returns notifications (paginated, unread filter)
func ListNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}
	page := 1
	perPage := 20
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			perPage = n
		}
	}
	offset := (page - 1) * perPage
	unread := false
	if v := r.URL.Query().Get("unread"); v == "true" {
		unread = true
	}
	notifications, err := NotificationUC.GetNotifications(businessID, unread, perPage, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := NotificationListResponse{
		Notifications: notifications,
		Page:          page,
		PerPage:       perPage,
		Count:         len(notifications),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// MarkNotificationReadHandler sets is_read=true for a notification
func MarkNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}
	notificationID := chi.URLParam(r, "id")
	if notificationID == "" {
		http.Error(w, "notification id required", http.StatusBadRequest)
		return
	}
	err := NotificationUC.MarkNotificationRead(notificationID, businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetDashboardStatsHandler returns business/branch dashboard stats
func GetDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	businessID, ok := middleware.GetBusinessIDFromContext(r.Context())
	if !ok || businessID == "" {
		http.Error(w, "missing or invalid business_id in token", http.StatusUnauthorized)
		return
	}

	// Determine if user is staff or owner
	var branchID string
	var realBusinessID string
	var businessName string
	staff, _ := StaffUC.StaffRepo.GetStaffByID(businessID)
	if staff != nil {
		// Staff user
		realBusinessID = staff.BusinessID
		branchID = staff.BranchID
		biz, err := BusinessUC.BusinessRepo.GetBusinessByID(realBusinessID)
		if err != nil {
			http.Error(w, "business not found", http.StatusNotFound)
			return
		}
		businessName = biz.Name
	} else {
		// Owner user
		realBusinessID = businessID
		branchID = ""
		biz, err := BusinessUC.BusinessRepo.GetBusinessByID(realBusinessID)
		if err != nil {
			http.Error(w, "business not found", http.StatusNotFound)
			return
		}
		businessName = biz.Name
	}

	// Get total sales today
	totalSalesToday, err := SaleUC.SaleRepo.GetTotalSalesToday(realBusinessID, branchID)
	if err != nil {
		http.Error(w, "failed to get total sales today", http.StatusInternalServerError)
		return
	}

	// Get total revenue
	totalRevenue, err := SaleUC.SaleRepo.GetTotalRevenue(realBusinessID, branchID)
	if err != nil {
		http.Error(w, "failed to get total revenue", http.StatusInternalServerError)
		return
	}

	// Get low stock count
	lowStockCount, err := ProductUC.ProductRepo.GetLowStockCount(realBusinessID, branchID)
	if err != nil {
		http.Error(w, "failed to get low stock count", http.StatusInternalServerError)
		return
	}

	// Get 5 recent transactions
	recentSales, err := SaleUC.SaleRepo.GetRecentSales(realBusinessID, branchID, 5)
	if err != nil {
		http.Error(w, "failed to get recent transactions", http.StatusInternalServerError)
		return
	}
	var recentTxs []RecentTransaction
	for _, s := range recentSales {
		recentTxs = append(recentTxs, RecentTransaction{
			SaleID:    s.ID,
			Amount:    s.TotalAmount,
			BranchID:  s.BranchID,
			CashierID: s.CashierID,
			CreatedAt: s.CreatedAt,
		})
	}

	resp := DashboardStatsResponse{
		BusinessName:       businessName,
		TotalSalesToday:    totalSalesToday,
		TotalRevenue:       totalRevenue,
		LowStockCount:      lowStockCount,
		RecentTransactions: recentTxs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DashboardStatsResponse represents the dashboard stats response
type DashboardStatsResponse struct {
	BusinessName       string              `json:"business_name"`
	TotalSalesToday    int                 `json:"total_sales_today"`
	TotalRevenue       float64             `json:"total_revenue"`
	LowStockCount      int                 `json:"low_stock_count"`
	RecentTransactions []RecentTransaction `json:"recent_transactions"`
}

// RecentTransaction represents a recent sale transaction for dashboard
type RecentTransaction struct {
	SaleID    string  `json:"sale_id"`
	Amount    float64 `json:"amount"`
	BranchID  string  `json:"branch_id"`
	CashierID string  `json:"cashier_id"`
	CreatedAt int64   `json:"created_at"`
}
