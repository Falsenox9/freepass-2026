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

func (r *Rest) CreateFood(c *gin.Context) {
	var param model.CreateFoodParam

	ownerID, exists := c.Get("user_id")
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

	resp, err := r.service.FoodService.CreateFood(ownerID.(uuid.UUID), param)
	if err != nil {
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create food", err)
		return
	}

	response.Success(c, http.StatusCreated, "food created successfully", resp)
}

func (r *Rest) BulkCreateFood(c *gin.Context) {
	var param model.BulkCreateFoodParam

	ownerID, exists := c.Get("user_id")
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

	resp, err := r.service.FoodService.BulkCreateFood(ownerID.(uuid.UUID), param)
	if err != nil {
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create foods", err)
		return
	}

	response.Success(c, http.StatusCreated, "foods created successfully", resp)
}

func (r *Rest) UpdateFood(c *gin.Context) {
	var param model.UpdateFoodParam

	ownerID, exists := c.Get("user_id")
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

	resp, err := r.service.FoodService.UpdateFood(ownerID.(uuid.UUID), param)
	if err != nil {
		if err.Error() == "food not found" {
			response.Error(c, http.StatusNotFound, "food not found", err)
			return
		}
		if err.Error() == "you don't have permission to update this food" {
			response.Error(c, http.StatusForbidden, "you don't have permission to update this food", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update food", err)
		return
	}

	response.Success(c, http.StatusOK, "food updated successfully", resp)
}

func (r *Rest) DeleteFood(c *gin.Context) {
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	foodIDStr := c.Param("food_id")
	foodID, err := uuid.Parse(foodIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid food id format", err)
		return
	}

	resp, err := r.service.FoodService.DeleteFood(ownerID.(uuid.UUID), foodID)
	if err != nil {
		if err.Error() == "food not found" {
			response.Error(c, http.StatusNotFound, "food not found", err)
			return
		}
		if err.Error() == "you don't have permission to delete this food" {
			response.Error(c, http.StatusForbidden, "you don't have permission to delete this food", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete food", err)
		return
	}

	response.Success(c, http.StatusOK, "food deleted successfully", resp)
}

func (r *Rest) GetFoodByID(c *gin.Context) {
	foodIDStr := c.Param("food_id")
	foodID, err := uuid.Parse(foodIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid food id format", err)
		return
	}

	resp, err := r.service.FoodService.GetFoodByID(foodID)
	if err != nil {
		if err.Error() == "food not found" {
			response.Error(c, http.StatusNotFound, "food not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get food", err)
		return
	}

	response.Success(c, http.StatusOK, "food retrieved successfully", resp)
}

func (r *Rest) GetMyFoods(c *gin.Context) {
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	resp, err := r.service.FoodService.GetFoodsByCanteenOwner(ownerID.(uuid.UUID))
	if err != nil {
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get foods", err)
		return
	}

	response.Success(c, http.StatusOK, "foods retrieved successfully", resp)
}

func (r *Rest) GetAllFoods(c *gin.Context) {
	resp, err := r.service.FoodService.GetAllFoods()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get all foods", err)
		return
	}

	response.Success(c, http.StatusOK, "foods retrieved successfully", resp)
}

func (r *Rest) UpdateStock(c *gin.Context) {
	var param model.UpdateStockParam

	ownerID, exists := c.Get("user_id")
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

	resp, err := r.service.FoodService.UpdateStock(ownerID.(uuid.UUID), param)
	if err != nil {
		if err.Error() == "you don't have permission to update this food stock" {
			response.Error(c, http.StatusForbidden, "you don't have permission to update this food stock", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update stock", err)
		return
	}

	response.Success(c, http.StatusOK, "stock updated successfully", resp)
}
