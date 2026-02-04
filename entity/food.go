package entity

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	FoodID      uuid.UUID `json:"food_id" gorm:"type:uuid;primaryKey"`
	CanteenID   uuid.UUID `json:"canteen_id" gorm:"type:uuid;not null"`
	FoodName    string    `json:"food_name" gorm:"type:varchar(255);not null"`
	Description *string   `json:"description" gorm:"type:text"`
	Price       int       `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null;default:0"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
