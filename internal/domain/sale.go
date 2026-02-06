package domain

type Sale struct {
	ID            string     `json:"id"`
	BusinessID    string     `json:"business_id"`
	BranchID      string     `json:"branch_id"`
	CashierID     string     `json:"cashier_id"`
	TotalAmount   float64    `json:"total_amount"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	CreatedAt     int64      `json:"created_at"`
	Items         []SaleItem `json:"items"`
}

type SaleItem struct {
	ID        string  `json:"id"`
	SaleID    string  `json:"sale_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
	CreatedAt int64   `json:"created_at"`
}

type SaleRepository interface {
	CreateSale(sale *Sale, items []SaleItem) (string, float64, error)
	GetTotalSalesToday(businessID, branchID string) (int, error)
	GetTotalRevenue(businessID, branchID string) (float64, error)
	GetRecentSales(businessID, branchID string, limit int) ([]*Sale, error)
}
