package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository  IUserRepository
	AdminRepository IAdminRepository
	FoodRepository  IFoodRepository
	OrderRepository IOrderRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:  NewUserRepository(db),
		AdminRepository: NewAdminRepository(db),
		FoodRepository:  NewFoodRepository(db),
		OrderRepository: NewOrderRepository(db),
	}
}
