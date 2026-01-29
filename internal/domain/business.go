package domain

type Business struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	OwnerFullName    string  `json:"owner_full_name"`
	Email            string  `json:"email"`
	PhoneNumber      string  `json:"phone_number"`
	PasswordHash     string  `json:"-"`
	StoreAddress     string  `json:"store_address"`
	BusinessCategory string  `json:"business_category"`
	Currency         string  `json:"currency"`
	StoreIcon        *string `json:"store_icon,omitempty"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
	Identifyer       string  `json:"identifyer"`
}

type BusinessRepository interface {
	CreateBusiness(b *Business) error
	GetBusinessByEmail(email string) (*Business, error)
	GetBusinessByID(id string) (*Business, error)
	GetBusinessByIdentifyer(identifyer string) (*Business, error)
}
