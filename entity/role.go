package entity

type Role struct {
	RoleID int    `json:"role_id" gorm:"type:int;primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"type:varchar(255);not null"`

	// Relation
	Users []User `json:"users" gorm:"foreignKey:RoleID"`
}
