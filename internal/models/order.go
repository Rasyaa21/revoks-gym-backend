package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"order_id"`
	Channel     string      `gorm:"type:varchar(50);not null" json:"channel"`
	BuyerUserID *uuid.UUID  `gorm:"type:uuid" json:"buyer_user_id"`
	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	TotalAmount float64     `gorm:"type:decimal;not null" json:"total_amount"`
	Currency    string      `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	CreatedAt   time.Time   `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	ExpiredAt   *time.Time  `gorm:"type:timestamp" json:"expired_at"`

	BuyerUser *User `gorm:"foreignKey:BuyerUserID;references:UserID" json:"buyer_user,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	OrderItemID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"order_item_id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ItemType    string    `gorm:"type:varchar(50);not null" json:"item_type"`
	RefID       uuid.UUID `gorm:"type:uuid;not null" json:"ref_id"`
	Qty         int       `gorm:"type:int;not null;default:1" json:"qty"`
	UnitPrice   float64   `gorm:"type:decimal;not null" json:"unit_price"`
	Subtotal    float64   `gorm:"type:decimal;not null" json:"subtotal"`

	Order *Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
