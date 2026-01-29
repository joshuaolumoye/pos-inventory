package usecase

import (
	"errors"
	"time"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type ProductUsecase struct {
	ProductRepo domain.ProductRepository
}

func (u *ProductUsecase) AddProduct(p *domain.Product) error {
	// Validate required fields
	if p.ProductName == "" || p.BusinessID == "" || p.BranchID == "" {
		return errors.New("missing required fields: product_name, business_id, or branch_id")
	}

	// Sanitize inputs
	p.ProductName = utils.Sanitize(p.ProductName)
	p.ProductCategory = utils.Sanitize(p.ProductCategory)
	if p.BarcodeValue != nil {
		sanitized := utils.Sanitize(*p.BarcodeValue)
		p.BarcodeValue = &sanitized
	}
	if p.NAFDACRegNumber != nil {
		sanitized := utils.Sanitize(*p.NAFDACRegNumber)
		p.NAFDACRegNumber = &sanitized
	}

	// Validate prices
	if p.SellingPrice <= 0 {
		return errors.New("selling price must be greater than 0")
	}
	if p.CostPrice < 0 {
		return errors.New("cost price cannot be negative")
	}

	// Validate stock
	if p.QuantityInStock < 0 {
		return errors.New("quantity in stock cannot be negative")
	}

	// Generate ID and timestamps
	p.ID = utils.GenerateUUID()
	now := time.Now().Unix()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Ensure CreatedBy is set (should come from JWT)
	if p.CreatedBy == "" {
		return errors.New("created_by is required")
	}

	return u.ProductRepo.CreateProduct(p)
}

func (u *ProductUsecase) UpdateProduct(p *domain.Product) error {
	// Validate required fields
	if p.ID == "" || p.BusinessID == "" {
		return errors.New("missing product_id or business_id")
	}

	// Get existing product to verify ownership and preserve certain fields
	existing, err := u.ProductRepo.GetProductByID(p.ID)
	if err != nil {
		return errors.New("product not found")
	}

	// Verify product belongs to the authenticated business
	if existing.BusinessID != p.BusinessID {
		return errors.New("unauthorized: product does not belong to your business")
	}

	// Sanitize inputs
	p.ProductName = utils.Sanitize(p.ProductName)
	p.ProductCategory = utils.Sanitize(p.ProductCategory)
	if p.BarcodeValue != nil {
		sanitized := utils.Sanitize(*p.BarcodeValue)
		p.BarcodeValue = &sanitized
	}
	if p.NAFDACRegNumber != nil {
		sanitized := utils.Sanitize(*p.NAFDACRegNumber)
		p.NAFDACRegNumber = &sanitized
	}

	// Validate prices
	if p.SellingPrice <= 0 {
		return errors.New("selling price must be greater than 0")
	}
	if p.CostPrice < 0 {
		return errors.New("cost price cannot be negative")
	}

	// Validate stock
	if p.QuantityInStock < 0 {
		return errors.New("quantity in stock cannot be negative")
	}

	// Preserve creation info and update timestamps
	p.CreatedAt = existing.CreatedAt
	p.CreatedBy = existing.CreatedBy
	p.UpdatedAt = time.Now().Unix()

	// Ensure UpdatedBy is set (should come from JWT)
	if p.UpdatedBy == nil || *p.UpdatedBy == "" {
		return errors.New("updated_by is required")
	}

	return u.ProductRepo.UpdateProduct(p)
}

func (u *ProductUsecase) GetProductByID(productID string) (*domain.Product, error) {
	if productID == "" {
		return nil, errors.New("missing product_id")
	}
	return u.ProductRepo.GetProductByID(productID)
}

func (u *ProductUsecase) GetProductsByBusinessID(businessID string) ([]*domain.Product, error) {
	if businessID == "" {
		return nil, errors.New("missing business_id")
	}
	return u.ProductRepo.GetProductsByBusinessID(businessID)
}

func (u *ProductUsecase) GetProductsByBranchID(businessID, branchID string) ([]*domain.Product, error) {
	if businessID == "" || branchID == "" {
		return nil, errors.New("missing business_id or branch_id")
	}
	return u.ProductRepo.GetProductsByBranchID(businessID, branchID)
}

func (u *ProductUsecase) DeleteProduct(productID, businessID string) error {
	if productID == "" || businessID == "" {
		return errors.New("missing product_id or business_id")
	}

	// Get existing product to verify ownership
	product, err := u.ProductRepo.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Verify product belongs to the authenticated business
	if product.BusinessID != businessID {
		return errors.New("unauthorized: product does not belong to your business")
	}

	return u.ProductRepo.DeleteProduct(productID)
}

func (u *ProductUsecase) SearchProducts(businessID, searchTerm string) ([]*domain.Product, error) {
	if businessID == "" {
		return nil, errors.New("missing business_id")
	}
	if searchTerm == "" {
		return u.ProductRepo.GetProductsByBusinessID(businessID)
	}
	searchTerm = utils.Sanitize(searchTerm)
	return u.ProductRepo.SearchProducts(businessID, searchTerm)
}

func (u *ProductUsecase) GetLowStockProducts(businessID string) ([]*domain.Product, error) {
	if businessID == "" {
		return nil, errors.New("missing business_id")
	}
	return u.ProductRepo.GetLowStockProducts(businessID)
}
