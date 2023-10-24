package message

type CreateOrder struct {
	Customer    string `json:"customer,omitempty"`
	OrderNumber string `json:"orderNumber,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
}

type Order struct {
	ID          int64  `json:"id"`
	Customer    string `json:"customer,omitempty"`
	OrderNumber string `json:"orderNumber,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	Created     int64  `json:"createdAt,omitempty"`
}
