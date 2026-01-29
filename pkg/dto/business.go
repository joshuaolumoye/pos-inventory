package dto

type RegisterBusinessRequest struct {
	BusinessName     string  `json:"business_name"`
	OwnerFullName    string  `json:"owner_full_name"`
	Email            string  `json:"email"`
	PhoneNumber      string  `json:"phone_number"`
	Password         string  `json:"password"`
	StoreAddress     string  `json:"store_address"`
	BusinessCategory string  `json:"business_category"`
	Currency         string  `json:"currency"`
	StoreIcon        *string `json:"store_icon,omitempty"`
}

type RegisterBusinessResponse struct {
	BusinessID string `json:"business_id"`
	Message    string `json:"message"`
}
