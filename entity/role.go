package entity

type Role struct {
	RoleID int    `gorm:"primaryKey;autoIncrement" json:"role_id"`
	Name   string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
}
