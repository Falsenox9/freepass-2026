package repository 

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAdminRepository interface {
	CreateCanteenOwner(tx *gorm.DB, user *entity.User, canteen *entity.Canteen)
	UpdateCanteenOwner(tx *gorm.DB, user *entity.User, canteen *entity.Canteen) error
	DeleteUser(tx *gorm.DB, userID uuid.UUID) error
	GetUserByID(userID uuid.UUID) (*entity.User, error)
	GetCanteenByOwbnerID(ownerID uuid.UUID) (*entity.Canteen, error)
	GetAllUsers() ([]entity.User, error)
	UpdateCanteenInfo(tx *gorm.DB, canteen *entity.Canteen) error
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) IAdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateCanteenOwner(tx *gorm.DB, user *entity.User, canteen *entity.Canteen) error {
    err := tx.Debug().Create(&user).Error
    if err != nil {
        return err
    }

    err = tx.Debug().Create(&canteen).Error
    if err != nil {
        return err
    }

    return nil
}

func (r *AdminRepository) UpdateCanteenOwner(tx *gorm.DB, user *entity.User, canteen *entity.Canteen) error {
    err := tx.Debug().Model(&entity.User{}).Where("user_id = ?", user.UserID).Updates(user).Error
    if err != nil {
        return err
    }

    if canteen != nil {
        err = tx.Debug().Model(&entity.Canteen{}).Where("owner_id = ?", user.UserID).Updates(canteen).Error
        if err != nil {
            return err
        }
    }

    return nil
}

func (r *AdminRepository) DeleteUser(tx *gorm.DB, userID uuid.UUID) error {
    err := tx.Debug().Where("user_id = ?", userID).Delete(&entity.User{}).Error
    if err != nil {
        return err
    }

    return nil
}

func (r *AdminRepository) GetUserByID(userID uuid.UUID) (*entity.User, error) {
    var user *entity.User
    err := r.db.Debug().Where("user_id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (r *AdminRepository) GetCanteenByOwnerID(ownerID uuid.UUID) (*entity.Canteen, error) {
    var canteen *entity.Canteen
    err := r.db.Debug().Where("owner_id = ?", ownerID).First(&canteen).Error
    if err != nil {
        return nil, err
    }

    return canteen, nil
}

func (r *AdminRepository) GetAllUsers() ([]*entity.User, error) {
    var users []*entity.User
    err := r.db.Debug().Preload("Role").Find(&users).Error
    if err != nil {
        return nil, err
    }

    return users, nil
}

func (r *AdminRepository) UpdateCanteenInfo(tx *gorm.DB, canteen *entity.Canteen) error {
    err := tx.Debug().Model(&entity.Canteen{}).Where("canteen_id = ?", canteen.CanteenID).Updates(canteen).Error
    if err != nil {
        return err
    }

    return nil
}



