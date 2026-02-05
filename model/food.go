package model

import "github.com/google/uuid"

type CreateFoodParam struct {
	FoodName    string  `json:"food_name" binding:"required"`
	Description *string `json:"description"`
	Price       int     `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
}

type CreateFoodResponse struct {
	FoodID      uuid.UUID `json:"food_id"`
	CanteenID   uuid.UUID `json:"canteen_id"`
	FoodName    string    `json:"food_name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
}

type BulkCreateFoodParam struct {
	Foods []CreateFoodParam `json:"foods" binding:"required,min=1,dive"`
}

type BulkCreateFoodResponse struct {
	Foods []CreateFoodResponse `json:"foods"`
	Count int                  `json:"count"`
}

type UpdateFoodParam struct {
	FoodID      uuid.UUID `json:"food_id" binding:"required"`
	FoodName    *string   `json:"food_name" binding:"omitempty"`
	Description *string   `json:"description"`
	Price       *int      `json:"price" binding:"omitempty,min=0"`
	Stock       *int      `json:"stock" binding:"omitempty,min=0"`
	IsAvailable *bool     `json:"is_available"`
}

type UpdateFoodResponse struct {
	FoodID      uuid.UUID `json:"food_id"`
	CanteenID   uuid.UUID `json:"canteen_id"`
	FoodName    string    `json:"food_name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
}

type DeleteFoodResponse struct {
	Message string `json:"message"`
}

type GetFoodListResponse struct {
	Foods []FoodInfo `json:"foods"`
}

type CanteenWithFoods struct {
	CanteenID   uuid.UUID  `json:"canteen_id"`
	CanteenName string     `json:"canteen_name"`
	Foods       []FoodInfo `json:"foods"`
}

type GetFoodsGroupedByCanteenResponse struct {
	Canteens []CanteenWithFoods `json:"canteens"`
}

type FoodInfo struct {
	FoodID      uuid.UUID `json:"food_id"`
	FoodName    string    `json:"food_name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
}

type GetFoodResponse struct {
	FoodID      uuid.UUID `json:"food_id"`
	CanteenID   uuid.UUID `json:"canteen_id"`
	CanteenName string    `json:"canteen_name"`
	FoodName    string    `json:"food_name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
}

type UpdateStockParam struct {
	FoodID uuid.UUID `json:"food_id" binding:"required"`
	Stock  int       `json:"stock" binding:"required,min=0"`
}

type UpdateStockResponse struct {
	FoodID uuid.UUID `json:"food_id"`
	Stock  int       `json:"stock"`
}

type UpdateCanteenStatusParam struct {
	IsOpen bool `json:"is_open" binding:"required"`
}

type UpdateCanteenStatusResponse struct {
	CanteenID   uuid.UUID `json:"canteen_id"`
	CanteenName string    `json:"canteen_name"`
	IsOpen      bool      `json:"is_open"`
}
