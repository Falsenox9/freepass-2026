package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID       uuid.UUID  `json:"order_id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	CanteenID     uuid.UUID  `json:"canteen_id" gorm:"type:uuid;not null"`
	TotalPrice    int        `json:"total_price" gorm:"not null"`
	Status        string     `json:"status" gorm:"type:varchar(50);default:'pending'"`
	PaymentStatus string     `json:"payment_status" gorm:"type:varchar(50);default:'unpaid'"`
	PaymentMethod *string    `json:"payment_method" gorm:"type:varchar(50)"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
