package dto

type BranchRequest struct {
	BranchName    string `json:"branch_name"`
	BranchAddress string `json:"branch_address"`
	IsMainBranch  bool   `json:"is_main_branch"`
}

type BranchResponse struct {
	BranchID      string `json:"branch_id"`
	BusinessID    string `json:"business_id"`
	BranchName    string `json:"branch_name"`
	BranchAddress string `json:"branch_address"`
	IsMainBranch  bool   `json:"is_main_branch"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	Message       string `json:"message,omitempty"`
}

type BranchListResponse struct {
	Branches []BranchResponse `json:"branches"`
	Count    int              `json:"count"`
}
