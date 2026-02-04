package service

import (
    "errors"
    "freepass-2026/entity"
    "freepass-2026/internal/repository"
    "freepass-2026/model"
    "freepass-2026/pkg/bcrypt"
    "freepass-2026/pkg/database/mariadb"
    "freepass-2026/pkg/jwt"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type IUserService interface {
    RegisterUser(param model.UserRegisterParam) (*model.UserRegisterResponse, error)
    LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error)
    GetUser(param model.UserParam) (*entity.User, error)
}

type UserService struct {
    db             *gorm.DB
    userRepository repository.IUserRepository
    bcrypt         bcrypt.Interface
    jwtAuth        jwt.Interface
}

func NewUserService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
    return &UserService{
        db:             mariadb.Connection,
        userRepository: userRepository,
        bcrypt:         bcrypt,
        jwtAuth:        jwtAuth,
    }
}

func (u *UserService) RegisterUser(param model.UserRegisterParam) (*model.UserRegisterResponse, error) {
    tx := u.db.Begin()
    defer tx.Rollback()

    _, err := u.userRepository.GetUser(model.UserParam{
        Email: param.Email,
    })
    if err == nil {
        return nil, errors.New("email already exists")
    }

    userID, err := uuid.NewUUID()
    if err != nil {
        return nil, err
    }

    if param.Password != param.ConfirmPassword {
        return nil, errors.New("password not match")
    }

    hashPassword, err := u.bcrypt.GenerateFromPassword(param.Password)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        UserID:   userID,
        RoleID:   2, // Regular user role
        Email:    param.Email,
        Password: hashPassword,
    }

    err = u.userRepository.CreateUser(tx, user)
    if err != nil {
        return nil, err
    }

    err = tx.Commit().Error
    if err != nil {
        return nil, err
    }

    response := &model.UserRegisterResponse{
        Email: user.Email,
    }

    return response, nil
}

func (u *UserService) LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error) {
    tx := u.db.Begin()
    defer tx.Rollback()

    user, err := u.userRepository.GetUser(model.UserParam{
        Email: param.Email,
    })
    if err != nil {
        return nil, errors.New("email or password is wrong")
    }

    err = u.bcrypt.CompareAndHashPassword(user.Password, param.Password)
    if err != nil {
        return nil, errors.New("email or password is wrong")
    }

    token, err := u.jwtAuth.CreateJWTToken(user.UserID, false)
    if err != nil {
        return nil, err
    }

    response := &model.UserLoginResponse{
        Token: token,
    }

    return response, nil
}

func (u *UserService) GetUser(param model.UserParam) (*entity.User, error) {
    return u.userRepository.GetUser(param)
}

func (u *UserService) GetUserProfile(userId uuid.UUID) (*model.UserProfile, error) {
    user, err := u.userRepository.GetUser(model.UserParam{
        UserID: userId,
    })
    if err != nil {
        return nil, err
    }

    response := &model.UserProfile{
        FullName: user.FullName,
        Email:    user.Email,
    }

    return response, nil
}