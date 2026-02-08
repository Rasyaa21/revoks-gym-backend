package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Payment struct {
	PaymentID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"payment_id"`
	OrderID                uuid.UUID     `gorm:"type:uuid;not null" json:"order_id"`
	Provider               string        `gorm:"type:varchar(50);not null;default:'MIDTRANS'" json:"provider"`
	MidtransOrderID        string        `gorm:"type:varchar(100);not null;uniqueIndex:uq_midtrans_order" json:"midtrans_order_id"`
	MidtransTransactionID  *string       `gorm:"type:varchar(100);uniqueIndex:uq_midtrans_trx" json:"midtrans_transaction_id"`
	PaymentType            *string       `gorm:"type:varchar(50)" json:"payment_type"`
	TransactionStatus      *string       `gorm:"type:varchar(50)" json:"transaction_status"`
	FraudStatus            *string       `gorm:"type:varchar(50)" json:"fraud_status"`
	GrossAmount            float64       `gorm:"type:decimal;not null" json:"gross_amount"`
	Currency               string        `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	ActionsJSON            datatypes.JSON `gorm:"type:jsonb" json:"actions_json"`
	VANumbersJSON          datatypes.JSON `gorm:"type:jsonb" json:"va_numbers_json"`
	MetadataJSON           datatypes.JSON `gorm:"type:jsonb" json:"metadata_json"`
	Status                 PaymentStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	PaidAt                 *time.Time    `gorm:"type:timestamp" json:"paid_at"`
	ExpiresAt              *time.Time    `gorm:"type:timestamp" json:"expires_at"`
	RawResponseJSON        datatypes.JSON `gorm:"type:jsonb" json:"raw_response_json"`
	CreatedAt              time.Time     `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	// Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentWebhookEvent struct {
	EventID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"event_id"`
	Provider               string         `gorm:"type:varchar(50);not null;default:'MIDTRANS'" json:"provider"`
	ProviderEventID        *string        `gorm:"type:varchar(100);uniqueIndex" json:"provider_event_id"`
	OrderID                *uuid.UUID     `gorm:"type:uuid" json:"order_id"`
	MidtransOrderID        *string        `gorm:"type:varchar(100)" json:"midtrans_order_id"`
	MidtransTransactionID  *string        `gorm:"type:varchar(100)" json:"midtrans_transaction_id"`
	RawPayloadJSON         datatypes.JSON `gorm:"type:jsonb;not null" json:"raw_payload_json"`
	ReceivedAt             time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"received_at"`
	ProcessedAt            *time.Time     `gorm:"type:timestamp" json:"processed_at"`
	ProcessStatus          string         `gorm:"type:varchar(20);not null;default:'RECEIVED'" json:"process_status"`

	// Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
}

func (PaymentWebhookEvent) TableName() string {
	return "payment_webhook_events"
}
