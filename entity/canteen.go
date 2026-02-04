package entity

import (
    "time"

    "github.com/google/uuid"
)

type Canteen struct {
    CanteenID   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"canteen_id"`
    OwnerID     uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex" json:"owner_id"`
    CanteenName string     `gorm:"type:varchar(100);not null" json:"canteen_name"`
    Description *string    `gorm:"type:text" json:"description"`
    Location    *string    `gorm:"type:varchar(255)" json:"location"`
    PhoneNumber *string    `gorm:"type:varchar(20)" json:"phone_number"`
    CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

    Owner *User   `gorm:"foreignKey:OwnerID;references:UserID" json:"owner,omitempty"`
    Foods []*Food `gorm:"foreignKey:CanteenID;references:CanteenID" json:"foods,omitempty"`
}