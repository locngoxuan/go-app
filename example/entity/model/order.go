package model

type Order struct {
	ID          int64  `gorm:"column:id"`
	Customer    string `gorm:"column:customer"`
	OrderNumber string `gorm:"column:order_number"`
	Amount      int64  `gorm:"column:amount"`
	Created     int64  `gorm:"column:created_at"`
}

func (Order) TableName() string {
	return "orders"
}
