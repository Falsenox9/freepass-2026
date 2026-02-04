package entity

import (
    "time"

    "github.com/google/uuid"
)

type User struct {
    UserID    uuid.UUID  `gorm:"type:char(36);primaryKey" json:"user_id"`
    RoleID    int        `gorm:"not null" json:"role_id"`
    FullName  *string    `gorm:"type:varchar(100)" json:"full_name"`
    Email     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Password  string     `gorm:"type:varchar(255);not null" json:"-"`
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

    Role *Role `gorm:"foreignKey:RoleID;references:RoleID" json:"role,omitempty"`
}