package model

import "github.com/google/uuid"

type OrderItemParam struct {
	FoodID   uuid.UUID `json:"food_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

type CreateOrderParam struct {
	Items []OrderItemParam `json:"items" binding:"required,min=1,dive"`
}

type OrderItemResponse struct {
	OrderItemID uuid.UUID `json:"order_item_id"`
	FoodID      uuid.UUID `json:"food_id"`
	FoodName    string    `json:"food_name"`
	Quantity    int       `json:"quantity"`
	Price       int       `json:"price"`
	Subtotal    int       `json:"subtotal"`
}

type CreateOrderResponse struct {
	OrderID    uuid.UUID           `json:"order_id"`
	CanteenID  uuid.UUID           `json:"canteen_id"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice int                 `json:"total_price"`
	Status     string              `json:"status"`
}

type GetOrderResponse struct {
	OrderID     uuid.UUID           `json:"order_id"`
	CanteenID   uuid.UUID           `json:"canteen_id"`
	CanteenName string              `json:"canteen_name"`
	Items       []OrderItemResponse `json:"items"`
	TotalPrice  int                 `json:"total_price"`
	Status      string              `json:"status"`
	CreatedAt   string              `json:"created_at"`
}

type GetOrderListResponse struct {
	Orders []GetOrderResponse `json:"orders"`
}
