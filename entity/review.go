package entity

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ReviewID  uuid.UUID `json:"review_id" gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
