package model

import "github.com/google/uuid"

type CreateCanteenOwnerParam struct {
	FullName string `json:"full_name" binding:"required"`
	Email 	string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
	CanteenName string `json:"canteen_name" binding:"required"`
}

type CreateCanteenOwnerResponse struct {
	UserID uuid.UUID `json:"user_id"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	CanteenName string `json:"canteen_name"`
}

type UpdateCanteenOwnerParam struct {
	UserId uuid.UUID `json:"user_id" binding:"required"`
	FullName string `json:"full_name"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
	CanteenName string `json:"canteen_name"`
}

type UpdateCanteenOwnerResponse struct {
	UserID uuid.UUID `json:"user_id"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	CanteenName string `json:"canteen_name"`
}

type DeleteUserParam struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type GetAllUsersResponse struct {
	Users []UserInfo `json:"users"`
}

type UserInfo struct {
	UserID uuid.UUID `json:"user_id"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	RoleID int `json:"role_id"`
	RoleName string `json:"role_name"`
}