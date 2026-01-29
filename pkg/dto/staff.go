package dto

type StaffRequest struct {
	StaffID     string  `json:"staff_id"`
	FullName    string  `json:"full_name"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	BranchID    string  `json:"branch_id"`
	Status      string  `json:"status"`
	PhotoURL    *string `json:"photo_url,omitempty"`
}

type StaffResponse struct {
	ID          string  `json:"id"`
	StaffID     string  `json:"staff_id"`
	FullName    string  `json:"full_name"`
	PhoneNumber string  `json:"phone_number"`
	Role        string  `json:"role"`
	BranchID    string  `json:"branch_id"`
	BusinessID  string  `json:"business_id"`
	Status      string  `json:"status"`
	PhotoURL    *string `json:"photo_url,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	Message     string  `json:"message,omitempty"`
}

type StaffListResponse struct {
	Staff []StaffResponse `json:"staff"`
	Count int             `json:"count"`
}

type StaffLoginRequest struct {
	BusinessIdentifyer string `json:"business_identifyer"`
	StaffID            string `json:"staff_id"`
	Password           string `json:"password"`
}

type StaffLoginResponse struct {
	StaffID      string `json:"staff_id"`
	StaffName    string `json:"staff_name"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message,omitempty"`
}
