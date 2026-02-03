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

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Register request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /users/profile [put]
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

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	// Public routes
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
