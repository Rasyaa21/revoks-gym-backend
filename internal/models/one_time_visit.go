package models

import (
	"time"

	"github.com/google/uuid"
)

type OneTimeProduct struct {
	OTVProductID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"otv_product_id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	PriceAmount     float64   `gorm:"type:decimal;not null" json:"price_amount"`
	ValidityMinutes int       `gorm:"type:int;not null" json:"validity_minutes"`
	IsActive        bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

func (OneTimeProduct) TableName() string {
	return "one_time_products"
}

type OneTimePass struct {
	OneTimePassID uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"one_time_pass_id"`
	OrderID       uuid.UUID    `gorm:"type:uuid;not null" json:"order_id"`
	OTVProductID  uuid.UUID    `gorm:"type:uuid;not null" json:"otv_product_id"`
	TokenHash     string       `gorm:"type:varchar(255);not null;uniqueIndex" json:"token_hash"`
	Status        OTVPassStatus `gorm:"type:varchar(20);not null;default:'UNUSED'" json:"status"`
	IssuedAt      time.Time    `gorm:"type:timestamptz;not null;autoCreateTime" json:"issued_at"`
	ExpiresAt     *time.Time   `gorm:"type:timestamp" json:"expires_at"`
	UsedAt        *time.Time   `gorm:"type:timestamp" json:"used_at"`

	Order   *Order          `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Product *OneTimeProduct `gorm:"foreignKey:OTVProductID;references:OTVProductID" json:"product,omitempty"`
}

func (OneTimePass) TableName() string {
	return "one_time_passes"
}
