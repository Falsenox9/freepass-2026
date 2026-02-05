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

func (r *Rest) CreateReview(c *gin.Context) {
	var param model.CreateReviewParam

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

	err = r.service.ReviewService.CreateReview(userID.(uuid.UUID), orderID, param)
	if err != nil {
		if err.Error() == "order not found" {
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		}
		if err.Error() == "unauthorized to review this order" {
			response.Error(c, http.StatusForbidden, "unauthorized to review this order", err)
			return
		}
		if err.Error() == "can only review paid orders" {
			response.Error(c, http.StatusBadRequest, "can only review paid orders", err)
			return
		}
		if err.Error() == "can only review completed orders" {
			response.Error(c, http.StatusBadRequest, "can only review completed orders", err)
			return
		}
		if err.Error() == "review already exists for this order" {
			response.Error(c, http.StatusBadRequest, "review already exists for this order", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create review", err)
		return
	}

	response.Success(c, http.StatusCreated, "review created successfully", nil)
}

func (r *Rest) GetCanteenReviews(c *gin.Context) {
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	resp, err := r.service.ReviewService.GetCanteenReviews(ownerID.(uuid.UUID))
	if err != nil {
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get reviews", err)
		return
	}

	response.Success(c, http.StatusOK, "reviews retrieved successfully", resp)
}

func (r *Rest) GetPublicCanteenReviews(c *gin.Context) {
	canteenIDStr := c.Param("canteen_id")
	canteenID, err := uuid.Parse(canteenIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid canteen id", err)
		return
	}

	resp, err := r.service.ReviewService.GetPublicCanteenReviews(canteenID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get reviews", err)
		return
	}

	response.Success(c, http.StatusOK, "reviews retrieved successfully", resp)
}

func (r *Rest) DeleteReview(c *gin.Context) {
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", fmt.Errorf("user not authenticated"))
		return
	}

	reviewIDStr := c.Param("review_id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid review id", err)
		return
	}

	err = r.service.ReviewService.DeleteReview(ownerID.(uuid.UUID), reviewID)
	if err != nil {
		if err.Error() == "review not found" {
			response.Error(c, http.StatusNotFound, "review not found", err)
			return
		}
		if err.Error() == "canteen not found for this owner" {
			response.Error(c, http.StatusNotFound, "canteen not found for this owner", err)
			return
		}
		if err.Error() == "unauthorized to delete this review" {
			response.Error(c, http.StatusForbidden, "unauthorized to delete this review", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete review", err)
		return
	}

	response.Success(c, http.StatusOK, "review deleted successfully", nil)
}
