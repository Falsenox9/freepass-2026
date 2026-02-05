package mariadb

import (
	"freepass-2026/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Canteen{},
		&entity.Food{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Review{},
	)

	if err != nil {
		return err
	}

	err = seedRoles(db)
	if err != nil {
		return err
	}

	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []entity.Role{
		{RoleID: 1, Name: "Admin"},
		{RoleID: 2, Name: "User"},
		{RoleID: 3, Name: "Canteen Owner"},
	}

	for _, role := range roles {
		var existingRole entity.Role
		result := db.Where("role_id = ?", role.RoleID).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
