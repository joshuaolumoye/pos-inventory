package domain

type StaffRole string

const (
	RoleOwner     StaffRole = "owner"
	RoleManager   StaffRole = "manager"
	RoleCashier   StaffRole = "cashier"
	RoleInventory StaffRole = "inventory_staff"
)

type Staff struct {
	ID           string    `json:"id"`
	StaffID      string    `json:"staff_id"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	PasswordHash string    `json:"-"`
	Role         StaffRole `json:"role"`
	BranchID     string    `json:"branch_id"`
	BusinessID   string    `json:"business_id"`
	Status       string    `json:"status"`
	PhotoURL     *string   `json:"photo_url,omitempty"`
	CreatedAt    int64     `json:"created_at"`
	UpdatedAt    int64     `json:"updated_at"`
	DeletedAt    *int64    `json:"deleted_at,omitempty"`
}

type StaffRepository interface {
	CreateStaff(staff *Staff) error
	GetStaffByStaffID(staffID string) (*Staff, error)
	GetStaffByID(id string) (*Staff, error)
	GetStaffByBusinessID(businessID string) ([]*Staff, error)
	GetStaffByBranchID(businessID, branchID string) ([]*Staff, error)
	UpdateStaff(staff *Staff) error
}
