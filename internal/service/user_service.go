package service

import (
	"errors"

	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"freepass-2026/pkg/bcrypt"
	"freepass-2026/pkg/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(req *model.RegisterRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(userID uuid.UUID) (*model.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req *model.UpdateProfileRequest) (*model.UserResponse, error)
	UpdatePassword(userID uuid.UUID, req *model.UpdatePasswordRequest) error
}

type userService struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
}

func NewUserService(userRepo repository.UserRepository, jwtService *jwt.JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *userService) Register(req *model.RegisterRequest) (*model.UserResponse, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !bcrypt.ComparePassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *toUserResponse(user),
	}, nil
}

func (s *userService) GetProfile(userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req *model.UpdateProfileRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) UpdatePassword(userID uuid.UUID, req *model.UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if !bcrypt.ComparePassword(user.Password, req.OldPassword) {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func toUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:     user.ID.String(),
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		RoleID: user.RoleID,
	}
}
