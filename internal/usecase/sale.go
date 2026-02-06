package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/joshuaolumoye/pos-backend/internal/domain"
)

type SaleUsecase struct {
	SaleRepo domain.SaleRepository
}

type SaleItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateSaleRequest struct {
	BranchID      string            `json:"branch_id"`
	PaymentMethod string            `json:"payment_method"`
	Items         []SaleItemRequest `json:"items"`
}

type CreateSaleResponse struct {
	Success     bool    `json:"success"`
	SaleID      string  `json:"sale_id"`
	TotalAmount float64 `json:"total_amount"`
}

func (u *SaleUsecase) CreateSale(req *CreateSaleRequest, businessID, cashierID string) (*CreateSaleResponse, error) {
	if businessID == "" || cashierID == "" {
		return nil, errors.New("unauthorized")
	}
	if req.BranchID == "" || req.PaymentMethod == "" || len(req.Items) == 0 {
		return nil, errors.New("invalid input")
	}
	var items []domain.SaleItem
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		items = append(items, domain.SaleItem{
			ID:        uuid.NewString(),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	sale := &domain.Sale{
		ID:            uuid.NewString(),
		BusinessID:    businessID,
		BranchID:      req.BranchID,
		CashierID:     cashierID,
		TotalAmount:   0,
		PaymentMethod: req.PaymentMethod,
		Status:        "completed",
		CreatedAt:     time.Now().Unix(),
	}
	saleID, total, err := u.SaleRepo.CreateSale(sale, items)
	if err != nil {
		return nil, err
	}
	return &CreateSaleResponse{
		Success:     true,
		SaleID:      saleID,
		TotalAmount: total,
	}, nil
}
