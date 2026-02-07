package repository

import (
	"errors"
	"fmt"
	"time"

	"database/sql"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/infrastructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleRepo struct {
	DB *gorm.DB
}

// GetTotalSalesToday returns the number of sales for today (optionally filtered by branch)
func (r *SaleRepo) GetTotalSalesToday(businessID, branchID string) (int, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour).Unix()
	query := r.DB.Model(&infrastructure.Sale{}).Where("business_id = ? AND created_at >= ?", businessID, today)
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	err := query.Count(&count).Error
	return int(count), err
}

// GetTotalRevenue returns the total revenue (optionally filtered by branch)
func (r *SaleRepo) GetTotalRevenue(businessID, branchID string) (float64, error) {
	// var total float64
	query := r.DB.Model(&infrastructure.Sale{}).Select("SUM(total_amount)").Where("business_id = ?", businessID)
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	// Handle NULL SUM(total_amount) gracefully

	var result sql.NullFloat64
	err := query.Row().Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if !result.Valid {
		return 0, nil
	}
	return result.Float64, nil
}

// GetRecentSales returns the 5 most recent sales (optionally filtered by branch)
func (r *SaleRepo) GetRecentSales(businessID, branchID string, limit int) ([]*domain.Sale, error) {
	var sales []*infrastructure.Sale
	query := r.DB.Where("business_id = ?", businessID)
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&sales).Error
	if err != nil {
		return nil, err
	}
	var result []*domain.Sale
	for _, s := range sales {
		result = append(result, &domain.Sale{
			ID:            s.ID,
			BusinessID:    s.BusinessID,
			BranchID:      s.BranchID,
			CashierID:     s.CashierID,
			TotalAmount:   s.TotalAmount,
			PaymentMethod: s.PaymentMethod,
			Status:        s.Status,
			CreatedAt:     s.CreatedAt,
		})
	}
	return result, nil
}

func NewSaleRepo(db *gorm.DB) *SaleRepo {
	return &SaleRepo{DB: db}
}

func (r *SaleRepo) CreateSale(sale *domain.Sale, items []domain.SaleItem) (string, float64, error) {
	var total float64
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Insert sale with total_amount = 0
		saleModel := infrastructure.Sale{
			ID:            sale.ID,
			BusinessID:    sale.BusinessID,
			BranchID:      sale.BranchID,
			CashierID:     sale.CashierID,
			TotalAmount:   0,
			PaymentMethod: sale.PaymentMethod,
			Status:        sale.Status,
			CreatedAt:     sale.CreatedAt,
		}
		if err := tx.Create(&saleModel).Error; err != nil {
			return err
		}

		for i := range items {
			// Lock product row FOR UPDATE
			var product infrastructure.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ?", items[i].ProductID).Error; err != nil {
				return errors.New("product not found")
			}
			if product.BusinessID != sale.BusinessID {
				return errors.New("product does not belong to business")
			}
			if product.QuantityInStock < items[i].Quantity {
				return fmt.Errorf("insufficient stock for product %s", product.ID)
			}
			if items[i].Quantity <= 0 {
				return errors.New("quantity must be greater than 0")
			}
			unitPrice := product.SellingPrice
			subtotal := unitPrice * float64(items[i].Quantity)
			// Insert sale item
			saleItemModel := infrastructure.SaleItem{
				ID:        items[i].ID,
				SaleID:    sale.ID,
				ProductID: items[i].ProductID,
				Quantity:  items[i].Quantity,
				UnitPrice: unitPrice,
				Subtotal:  subtotal,
				CreatedAt: time.Now().Unix(),
			}
			if err := tx.Create(&saleItemModel).Error; err != nil {
				return err
			}
			// Update product stock
			newStock := product.QuantityInStock - items[i].Quantity
			if newStock < 0 {
				return errors.New("stock would become negative")
			}
			if err := tx.Model(&product).Update("quantity_in_stock", newStock).Error; err != nil {
				return err
			}
			total += subtotal
		}
		// Update sale total_amount
		if err := tx.Model(&saleModel).Update("total_amount", total).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", 0, err
	}
	return sale.ID, total, nil
}
