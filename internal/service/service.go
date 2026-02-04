package service

import (
    "freepass-2026/internal/repository"
    "freepass-2026/pkg/bcrypt"
    "freepass-2026/pkg/jwt"
)

type Service struct {
    UserService  IUserService
    AdminService IAdminService
    FoodService  IFoodService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
    return &Service{
        UserService:  NewUserService(repository.UserRepository, bcrypt, jwtAuth),
        AdminService: NewAdminService(repository.AdminRepository, repository.UserRepository, bcrypt),
        FoodService:  NewFoodService(repository.FoodRepository, repository.AdminRepository),
    }
}