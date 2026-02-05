package entity

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	OrderItemID uuid.UUID `json:"order_item_id" gorm:"type:uuid;primaryKey"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	FoodID      uuid.UUID `json:"food_id" gorm:"type:uuid;not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Subtotal    int       `json:"subtotal" gorm:"not null"`
}
