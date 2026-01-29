package dto

// ProductRequest represents the request body for creating/updating a product
type ProductRequest struct {
	ProductName       string  `json:"product_name"`
	ProductCategory   string  `json:"product_category"`
	SellingPrice      float64 `json:"selling_price"`
	CostPrice         float64 `json:"cost_price"`
	QuantityInStock   int     `json:"quantity_in_stock"`
	LowStockThreshold int     `json:"low_stock_threshold"`
	BarcodeValue      string  `json:"barcode_value"`
	NAFDACRegNumber   string  `json:"nafdac_reg_number"`
	ExpiryDate        *int64  `json:"expiry_date,omitempty"`
	ProductImageURL   *string `json:"product_image_url,omitempty"`
	BranchID          string  `json:"branch_id"`
}

// ProductResponse represents the response for a product
type ProductResponse struct {
	ID              string  `json:"id"`
	ProductName     string  `json:"product_name"`
	ProductCategory string  `json:"product_category,omitempty"`
	NAFDACRegNumber string  `json:"nafdac_reg_number,omitempty"`
	SellingPrice    float64 `json:"selling_price"`
	CostPrice       float64 `json:"cost_price,omitempty"`
	QuantityLeft    int     `json:"quantity_left"`
	ProductImageURL *string `json:"product_image_url,omitempty"`
	BranchID        string  `json:"branch_id,omitempty"`
	BusinessID      string  `json:"business_id,omitempty"`
	Message         string  `json:"message,omitempty"`
}

// ProductListResponse represents a list of products
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Count    int               `json:"count"`
}
