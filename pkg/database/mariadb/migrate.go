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
    )

    if err != nil {
        return err
    }

    return nil
}