package model

import (
    "github.com/google/uuid"
)

type UserParam struct {
    UserID   uuid.UUID `json:"-"`
    FullName string    `json:"-"`
    Email    string    `json:"-"`
}

type UserRegisterParam struct {
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type UserRegisterResponse struct {
    Email string `json:"email"`
}

type UserLoginParam struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
    Token string `json:"token"`
}

type UserProfile struct {
    FullName *string `json:"full_name,omitempty"`
    Email    string  `json:"email"`
}