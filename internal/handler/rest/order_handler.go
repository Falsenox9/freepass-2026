package rest

import (
	"fmt"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (r *Rest) CreateOrder(c *gin.Context) {
	var param model.CreateOrderParam

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make([]string, 0)
			for _, e := range validationErrors {
				msg := fmt.Sprintf("Field '%s' error: %s", e.Field(), e.Tag())
				errorMessages = append(errorMessages, msg)
			}
			response.Error(c, http.StatusBadRequest, "invalid validation format", fmt.Errorf("%v", errorMessages))
			return
		}
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.OrderService.CreateOrder(userID.(uuid.UUID), param)
	if err != nil {
		if err.Error() == "food not found" {
			response.Error(c, http.StatusNotFound, "food not found", err)
			return
		}
		if err.Error() == "food is not available" {
			response.Error(c, http.StatusBadRequest, "food is not available", err)
			return
		}
		if err.Error() == "all items must be from the same canteen" {
			response.Error(c, http.StatusBadRequest, "all items must be from the same canteen", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create order", err)
		return
	}

	response.Success(c, http.StatusCreated, "order created successfully", resp)
}

func (r *Rest) GetMyOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	resp, err := r.service.OrderService.GetUserOrders(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	response.Success(c, http.StatusOK, "orders retrieved successfully", resp)
}

func (r *Rest) GetCanteenOrders(c *gin.Context) {
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	resp, err := r.service.OrderService.GetCanteenOrders(ownerID.(uuid.UUID))
	if err != nil {
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get canteen orders", err)
		return
	}

	response.Success(c, http.StatusOK, "canteen orders retrieved successfully", resp)
}

func (r *Rest) PayOrder(c *gin.Context) {
	var param model.PayOrderParam

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id", err)
		return
	}

	err = c.ShouldBindJSON(&param)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make([]string, 0)
			for _, e := range validationErrors {
				msg := fmt.Sprintf("Field '%s' error: %s", e.Field(), e.Tag())
				errorMessages = append(errorMessages, msg)
			}
			response.Error(c, http.StatusBadRequest, "invalid validation format", fmt.Errorf("%v", errorMessages))
			return
		}
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.OrderService.PayOrder(userID.(uuid.UUID), orderID, param)
	if err != nil {
		if err.Error() == "order not found" {
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		}
		if err.Error() == "unauthorized to pay this order" {
			response.Error(c, http.StatusForbidden, "unauthorized to pay this order", err)
			return
		}
		if err.Error() == "order already paid" {
			response.Error(c, http.StatusBadRequest, "order already paid", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to process payment", err)
		return
	}

	response.Success(c, http.StatusOK, "payment processed successfully", resp)
}

func (r *Rest) UpdateOrderStatus(c *gin.Context) {
	var param model.UpdateOrderStatusParam

	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id", err)
		return
	}

	err = c.ShouldBindJSON(&param)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make([]string, 0)
			for _, e := range validationErrors {
				msg := fmt.Sprintf("Field '%s' error: %s", e.Field(), e.Tag())
				errorMessages = append(errorMessages, msg)
			}
			response.Error(c, http.StatusBadRequest, "invalid validation format", fmt.Errorf("%v", errorMessages))
			return
		}
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	err = r.service.OrderService.UpdateOrderStatus(ownerID.(uuid.UUID), orderID, param)
	if err != nil {
		if err.Error() == "order not found" {
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		}
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		if err.Error() == "unauthorized to update this order" {
			response.Error(c, http.StatusForbidden, "unauthorized to update this order", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update order status", err)
		return
	}

	response.Success(c, http.StatusOK, "order status updated successfully", nil)
}
