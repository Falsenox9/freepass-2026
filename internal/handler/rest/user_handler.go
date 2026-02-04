package rest

import (
    "fmt"
    "freepass-2026/model"
    "freepass-2026/pkg/response"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

func (r *Rest) RegisterUser(c *gin.Context) {
    var param model.UserRegisterParam

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

    resp, err := r.service.UserService.RegisterUser(param)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "failed to register user", err)
        return
    }

    response.Success(c, http.StatusOK, "user registered successfully", resp)
}

func (r *Rest) LoginUser(c *gin.Context) {
    var param model.UserLoginParam

    err := c.ShouldBindJSON(&param)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "failed to bind json", err)
        return
    }

    resp, err := r.service.UserService.LoginUser(param)
    if err != nil {
        if err.Error() == "email or password is wrong" {
            response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
            return
        }
        response.Error(c, http.StatusInternalServerError, "failed to login user", err)
        return
    }

    response.Success(c, http.StatusOK, "user logged in successfully", resp)
}