package entity

import (
    "time"

    "github.com/google/uuid"
)

type Food struct {
    FoodID      uuid.UUID  `gorm:"type:char(36);primaryKey" json:"food_id"`
    CanteenID   uuid.UUID  `gorm:"type:char(36);not null" json:"canteen_id"`
    FoodName    string     `gorm:"type:varchar(100);not null" json:"food_name"`
    Description *string    `gorm:"type:text" json:"description"`
    Price       int        `gorm:"not null" json:"price"`
    Stock       int        `gorm:"not null;default:0" json:"stock"`
    ImageURL    *string    `gorm:"type:varchar(255)" json:"image_url"`
    IsAvailable bool       `gorm:"default:true" json:"is_available"`
    CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

    Canteen *Canteen `gorm:"foreignKey:CanteenID;references:CanteenID" json:"canteen,omitempty"`
}