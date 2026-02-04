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

func (r *Rest) CreateCanteenOwner(c *gin.Context) {
    var param model.CreateCanteenOwnerParam

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

    resp, err := r.service.AdminService.CreateCanteenOwner(param)
    if err != nil {
        if err.Error() == "email already exists" {
            response.Error(c, http.StatusConflict, "email already exists", err)
            return
        }
        response.Error(c, http.StatusInternalServerError, "failed to create canteen owner", err)
        return
    }

    response.Success(c, http.StatusCreated, "canteen owner created successfully", resp)
}

func (r *Rest) UpdateCanteenOwner(c *gin.Context) {
    var param model.UpdateCanteenOwnerParam

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

    resp, err := r.service.AdminService.UpdateCanteenOwner(param)
    if err != nil {
        if err.Error() == "user not found" {
            response.Error(c, http.StatusNotFound, "user not found", err)
            return
        }
        if err.Error() == "user is not a canteen owner" {
            response.Error(c, http.StatusBadRequest, "user is not a canteen owner", err)
            return
        }
        if err.Error() == "email already exists" {
            response.Error(c, http.StatusConflict, "email already exists", err)
            return
        }
        response.Error(c, http.StatusInternalServerError, "failed to update canteen owner", err)
        return
    }

    response.Success(c, http.StatusOK, "canteen owner updated successfully", resp)
}

func (r *Rest) DeleteUser(c *gin.Context) {
    userIDStr := c.Param("user_id")
    
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "invalid user id format", err)
        return
    }

    resp, err := r.service.AdminService.DeleteUser(userID)
    if err != nil {
        if err.Error() == "user not found" {
            response.Error(c, http.StatusNotFound, "user not found", err)
            return
        }
        response.Error(c, http.StatusInternalServerError, "failed to delete user", err)
        return
    }

    response.Success(c, http.StatusOK, "user deleted successfully", resp)
}

func (r *Rest) GetAllUsers(c *gin.Context) {
    resp, err := r.service.AdminService.GetAllUsers()
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "failed to get all users", err)
        return
    }

    response.Success(c, http.StatusOK, "users retrieved successfully", resp)
}