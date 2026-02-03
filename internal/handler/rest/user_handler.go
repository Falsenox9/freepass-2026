package rest

import (
	"net/http"

	"freepass-2026/internal/service"
	"freepass-2026/model"
	"freepass-2026/pkg/middleware"
	"freepass-2026/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	loginResp, err := h.userService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResp)
}


func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}


func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", user)
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update current user's password
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdatePasswordRequest true "Update password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Error(c, http.StatusUnauthorized, "Current password is incorrect")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password updated successfully", nil)
}


func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	// Protected user routes
	users := r.Group("/users")
	users.Use(authMiddleware.Authenticate())
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.PUT("/password", h.UpdatePassword)
	}
}
