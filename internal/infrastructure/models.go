package infrastructure

type Business struct {
	ID               string  `gorm:"primaryKey;type:char(36)" json:"id"`
	CreatedAt        int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        int64   `gorm:"autoUpdateTime" json:"updated_at"`
	Name             string  `gorm:"not null" json:"name"`
	OwnerFullName    string  `gorm:"not null" json:"owner_full_name"`
	Email            string  `gorm:"uniqueIndex;size:191;not null" json:"email"`
	PhoneNumber      string  `gorm:"not null" json:"phone_number"`
	PasswordHash     string  `gorm:"not null" json:"-"`
	StoreAddress     string  `gorm:"not null" json:"store_address"`
	BusinessCategory string  `gorm:"not null" json:"business_category"`
	Currency         string  `gorm:"not null" json:"currency"`
	StoreIcon        *string `gorm:"type:text" json:"store_icon,omitempty"`
	Identifyer       string  `gorm:"uniqueIndex;size:20;not null" json:"identifyer"`

	// Relationships
	Branches []Branch  `gorm:"foreignKey:BusinessID" json:"branches,omitempty"`
	Staff    []Staff   `gorm:"foreignKey:BusinessID" json:"staff,omitempty"`
	Products []Product `gorm:"foreignKey:BusinessID" json:"products,omitempty"`
}

type Branch struct {
	ID            string `gorm:"primaryKey;type:char(36)" json:"id"`
	CreatedAt     int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     int64  `gorm:"autoUpdateTime" json:"updated_at"`
	BusinessID    string `gorm:"index;not null;type:char(36)" json:"business_id"`
	BranchName    string `gorm:"not null" json:"branch_name"`
	BranchAddress string `gorm:"not null" json:"branch_address"`
	IsMainBranch  bool   `gorm:"default:false" json:"is_main_branch"`

	// Relationships
	Business Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	Staff    []Staff   `gorm:"foreignKey:BranchID" json:"staff,omitempty"`
	Products []Product `gorm:"foreignKey:BranchID" json:"products,omitempty"`
}

type StaffRole string

const (
	RoleOwner     StaffRole = "owner"
	RoleManager   StaffRole = "manager"
	RoleCashier   StaffRole = "cashier"
	RoleInventory StaffRole = "inventory_staff"
)

type Staff struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	CreatedAt    int64     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    int64     `gorm:"autoUpdateTime" json:"updated_at"`
	StaffID      string    `gorm:"uniqueIndex;size:191;not null" json:"staff_id"`
	FullName     string    `gorm:"not null" json:"full_name"`
	PhoneNumber  string    `gorm:"not null" json:"phone_number"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         StaffRole `gorm:"type:varchar(32);not null" json:"role"`
	BranchID     string    `gorm:"index;not null;type:char(36)" json:"branch_id"`
	BusinessID   string    `gorm:"index;not null;type:char(36)" json:"business_id"`
	Status       string    `gorm:"type:varchar(16);not null" json:"status"`
	PhotoURL     *string   `gorm:"type:text" json:"photo_url,omitempty"`

	// Relationships
	Branch   Branch   `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Business Business `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
}

type Product struct {
	ID                string  `gorm:"primaryKey;type:char(36)" json:"id"`
	ProductName       string  `gorm:"not null" json:"product_name"`
	ProductCategory   string  `gorm:"not null" json:"product_category"`
	BusinessID        string  `gorm:"index;not null;type:char(36)" json:"business_id"`
	BranchID          string  `gorm:"index;not null;type:char(36)" json:"branch_id"`
	BarcodeValue      *string `gorm:"uniqueIndex" json:"barcode_value,omitempty"`
	NAFDACRegNumber   *string `json:"nafdac_reg_number,omitempty"`
	SellingPrice      float64 `gorm:"not null" json:"selling_price"`
	CostPrice         float64 `gorm:"not null" json:"cost_price"`
	QuantityInStock   int     `gorm:"not null" json:"quantity_in_stock"`
	LowStockThreshold int     `gorm:"not null" json:"low_stock_threshold"`
	ExpiryDate        *int64  `json:"expiry_date,omitempty"`
	ProductImageURL   *string `gorm:"type:text" json:"product_image_url,omitempty"`
	CreatedAt         int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         int64   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         *int64  `json:"deleted_at,omitempty"`
	CreatedBy         string  `gorm:"type:char(36);not null" json:"created_by"`
	UpdatedBy         *string `gorm:"type:char(36)" json:"updated_by,omitempty"`

	// Relationships
	Branch   Branch   `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Business Business `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
}
