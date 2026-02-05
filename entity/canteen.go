package entity

import (
	"time"

	"github.com/google/uuid"
)

type Canteen struct {
	CanteenID   uuid.UUID `json:"canteen_id" gorm:"type:uuid;primaryKey"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;unique"`
	CanteenName string    `json:"canteen_name" gorm:"type:varchar(255);not null"`
	Description *string   `json:"description" gorm:"type:text"`
	Location    *string   `json:"location" gorm:"type:varchar(255)"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20)"`
	IsOpen      bool      `json:"is_open" gorm:"type:boolean;not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Foods []Food `json:"foods" gorm:"foreignKey:CanteenID"`
}
