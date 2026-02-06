package domain

type Product struct {
	ID                string  `db:"id" json:"id"`
	ProductName       string  `db:"product_name" json:"product_name"`
	ProductCategory   string  `db:"product_category" json:"product_category"`
	BusinessID        string  `db:"business_id" json:"business_id"`
	BranchID          string  `db:"branch_id" json:"branch_id"`
	BarcodeValue      *string `db:"barcode_value" json:"barcode_value"`
	NAFDACRegNumber   *string `db:"nafdac_reg_number" json:"nafdac_reg_number"`
	SellingPrice      float64 `db:"selling_price" json:"selling_price"`
	CostPrice         float64 `db:"cost_price" json:"cost_price"`
	QuantityInStock   int     `db:"quantity_in_stock" json:"quantity_in_stock"`
	LowStockThreshold int     `db:"low_stock_threshold" json:"low_stock_threshold"`
	ExpiryDate        *int64  `db:"expiry_date" json:"expiry_date,omitempty"`
	ProductImageURL   *string `db:"product_image_url" json:"product_image_url,omitempty"`
	CreatedAt         int64   `db:"created_at" json:"created_at"`
	UpdatedAt         int64   `db:"updated_at" json:"updated_at"`
	DeletedAt         *int64  `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedBy         string  `db:"created_by" json:"created_by"`
	UpdatedBy         *string `db:"updated_by" json:"updated_by,omitempty"`
}

type ProductRepository interface {
	CreateProduct(product *Product) error
	GetProductByID(productID string) (*Product, error)
	GetProductsByBusinessID(businessID string) ([]*Product, error)
	GetProductsByBranchID(businessID, branchID string) ([]*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(productID string) error
	SearchProducts(businessID, searchTerm string) ([]*Product, error)
	GetLowStockProducts(businessID string) ([]*Product, error)
	UpdateProductStock(productID string, quantity int) error
	GetProductsByBranch(branchID string) ([]*Product, error) // Keep for backward compatibility
	QueryProductsNotification(businessID, op string, stock int, expiry int64, lowStock int, limit, offset int, expired bool) ([]*Product, error)
	GetAllProductsPaginated(businessID string, limit, offset int) ([]*Product, error)
	SearchProductsPaginated(businessID, search string, limit, offset int) ([]*Product, error)
	GetLowStockCount(businessID, branchID string) (int, error)
}

type Notification struct {
	ID               string `json:"id"`
	BusinessID       string `json:"business_id"`
	ProductID        string `json:"product_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
	IsRead           bool   `json:"is_read"`
	CreatedAt        int64  `json:"created_at"`
}

type NotificationRepository interface {
	CreateNotification(n *Notification) error
	GetNotifications(businessID string, unreadOnly bool, limit, offset int) ([]*Notification, error)
	MarkNotificationRead(notificationID, businessID string) error
	ExistsUnreadLowStockNotification(businessID, productID string) (bool, error)
}
