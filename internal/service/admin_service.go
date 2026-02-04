package service

import (
    "errors"
    "freepass-2026/entity"
    "freepass-2026/internal/repository"
    "freepass-2026/model"
    "freepass-2026/pkg/bcrypt"
    "freepass-2026/pkg/database/mariadb"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type IAdminService interface {
    CreateCanteenOwner(param model.CreateCanteenOwnerParam) (*model.CreateCanteenOwnerResponse, error)
    UpdateCanteenOwner(param model.UpdateCanteenOwnerParam) (*model.UpdateCanteenOwnerResponse, error)
    DeleteUser(userID uuid.UUID) (*model.DeleteUserResponse, error)
    GetAllUsers() (*model.GetAllUsersResponse, error)
}

type AdminService struct {
    db              *gorm.DB
    adminRepository repository.IAdminRepository
    userRepository  repository.IUserRepository
    bcrypt          bcrypt.Interface
}

func NewAdminService(adminRepository repository.IAdminRepository, userRepository repository.IUserRepository, bcrypt bcrypt.Interface) IAdminService {
    return &AdminService{
        db:              mariadb.Connection,
        adminRepository: adminRepository,
        userRepository:  userRepository,
        bcrypt:          bcrypt,
    }
}

func (s *AdminService) CreateCanteenOwner(param model.CreateCanteenOwnerParam) (*model.CreateCanteenOwnerResponse, error) {
    tx := s.db.Begin()
    defer tx.Rollback()

    
    _, err := s.userRepository.GetUser(model.UserParam{
        Email: param.Email,
    })
    if err == nil {
        return nil, errors.New("email already exists")
    }

    
    if param.Password != param.ConfirmPassword {
        return nil, errors.New("password not match")
    }

    
    userID, err := uuid.NewUUID()
    if err != nil {
        return nil, err
    }

   
    hashPassword, err := s.bcrypt.GenerateFromPassword(param.Password)
    if err != nil {
        return nil, err
    }

    
    user := &entity.User{
        UserID:   userID,
        RoleID:   3, 
        FullName: &param.FullName,
        Email:    param.Email,
        Password: hashPassword,
    }

    
    canteenID, err := uuid.NewUUID()
    if err != nil {
        return nil, err
    }

    
    canteen := &entity.Canteen{
        CanteenID:   canteenID,
        OwnerID:     userID,
        CanteenName: param.CanteenName,
    }

    
    err = s.adminRepository.CreateCanteenOwner(tx, user, canteen)
    if err != nil {
        return nil, err
    }

    err = tx.Commit().Error
    if err != nil {
        return nil, err
    }

    response := &model.CreateCanteenOwnerResponse{
        UserID:      user.UserID,
        FullName:    *user.FullName,
        Email:       user.Email,
        CanteenName: canteen.CanteenName,
    }

    return response, nil
}

func (s *AdminService) UpdateCanteenOwner(param model.UpdateCanteenOwnerParam) (*model.UpdateCanteenOwnerResponse, error) {
    tx := s.db.Begin()
    defer tx.Rollback()

    
    existingUser, err := s.adminRepository.GetUserByID(param.UserID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    
    if existingUser.RoleID != 3 {
        return nil, errors.New("user is not a canteen owner")
    }

    
    if param.Email != "" && param.Email != existingUser.Email {
        _, err := s.userRepository.GetUser(model.UserParam{
            Email: param.Email,
        })
        if err == nil {
            return nil, errors.New("email already exists")
        }
    }

    
    if param.FullName != "" {
        existingUser.FullName = &param.FullName
    }
    if param.Email != "" {
        existingUser.Email = param.Email
    }

    
    var canteen *entity.Canteen
    if param.CanteenName != "" {
        canteen, err = s.adminRepository.GetCanteenByOwnerID(param.UserID)
        if err != nil {
            return nil, errors.New("canteen not found")
        }
        canteen.CanteenName = param.CanteenName
    }

    
    err = s.adminRepository.UpdateCanteenOwner(tx, existingUser, canteen)
    if err != nil {
        return nil, err
    }

    err = tx.Commit().Error
    if err != nil {
        return nil, err
    }


    updatedCanteen, _ := s.adminRepository.GetCanteenByOwnerID(param.UserID)

    response := &model.UpdateCanteenOwnerResponse{
        UserID:      existingUser.UserID,
        FullName:    *existingUser.FullName,
        Email:       existingUser.Email,
        CanteenName: updatedCanteen.CanteenName,
    }

    return response, nil
}

func (s *AdminService) DeleteUser(userID uuid.UUID) (*model.DeleteUserResponse, error) {
    tx := s.db.Begin()
    defer tx.Rollback()

    
    _, err := s.adminRepository.GetUserByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    
    err = s.adminRepository.DeleteUser(tx, userID)
    if err != nil {
        return nil, err
    }

    err = tx.Commit().Error
    if err != nil {
        return nil, err
    }

    response := &model.DeleteUserResponse{
        Message: "user deleted successfully",
    }

    return response, nil
}

func (s *AdminService) GetAllUsers() (*model.GetAllUsersResponse, error) {
    users, err := s.adminRepository.GetAllUsers()
    if err != nil {
        return nil, err
    }

    userInfos := make([]model.UserInfo, 0)
    for _, user := range users {
        fullName := ""
        if user.FullName != nil {
            fullName = *user.FullName
        }

        roleName := ""
        if user.Role != nil {
            roleName = user.Role.RoleName
        }

        userInfos = append(userInfos, model.UserInfo{
            UserID:   user.UserID,
            FullName: fullName,
            Email:    user.Email,
            RoleID:   user.RoleID,
            RoleName: roleName,
        })
    }

    response := &model.GetAllUsersResponse{
        Users: userInfos,
    }

    return response, nil
}