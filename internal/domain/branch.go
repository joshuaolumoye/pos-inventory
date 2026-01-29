package domain

type Branch struct {
	ID            string `json:"id"`
	BusinessID    string `json:"business_id"`
	BranchName    string `json:"branch_name"`
	BranchAddress string `json:"branch_address"`
	IsMainBranch  bool   `json:"is_main_branch"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	DeletedAt     *int64 `json:"deleted_at,omitempty"`
}

type BranchRepository interface {
	CreateBranch(b *Branch) error
	GetBranchByID(id string) (*Branch, error)
	GetBranchesByBusinessID(businessID string) ([]*Branch, error)
	DeleteBranch(id string) error
	UpdateBranch(b *Branch) error
}
