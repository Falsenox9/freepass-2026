package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFoodRepository interface {
	CreateFood(tx *gorm.DB, food *entity.Food) error
	UpdateFood(tx *gorm.DB, food *entity.Food) error
	DeleteFood(tx *gorm.DB, foodID uuid.UUID) error
	GetFoodByID(foodID uuid.UUID) (*entity.Food, error)
	GetFoodsByCanteenID(canteenID uuid.UUID) ([]*entity.Food, error)
	GetAllFoods() ([]*entity.Food, error)
	UpdateStock(tx *gorm.DB, foodID uuid.UUID, stock int) error
	CheckFoodOwnership(foodID uuid.UUID, canteenID uuid.UUID) (bool, error)
}

type FoodRepository struct {
	db *gorm.DB
}

func NewFoodRepository(db *gorm.DB) IFoodRepository {
	return &FoodRepository{db: db}
}

func (r *FoodRepository) CreateFood(tx *gorm.DB, food *entity.Food) error {
	err := tx.Debug().Create(&food).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) UpdateFood(tx *gorm.DB, food *entity.Food) error {
	err := tx.Debug().Model(&entity.Food{}).Where("food_id = ?", food.FoodID).Updates(food).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) DeleteFood(tx *gorm.DB, foodID uuid.UUID) error {
	err := tx.Debug().Where("food_id = ?", foodID).Delete(&entity.Food{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) GetFoodByID(foodID uuid.UUID) (*entity.Food, error) {
	var food *entity.Food
	err := r.db.Debug().Where("food_id = ?", foodID).First(&food).Error
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (r *FoodRepository) GetFoodsByCanteenID(canteenID uuid.UUID) ([]*entity.Food, error) {
	var foods []*entity.Food
	err := r.db.Debug().Where("canteen_id = ?", canteenID).Find(&foods).Error
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (r *FoodRepository) GetAllFoods() ([]*entity.Food, error) {
	var foods []*entity.Food
	err := r.db.Debug().Find(&foods).Error
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (r *FoodRepository) UpdateStock(tx *gorm.DB, foodID uuid.UUID, stock int) error {
	err := tx.Debug().Model(&entity.Food{}).Where("food_id = ?", foodID).Update("stock", stock).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) CheckFoodOwnership(foodID uuid.UUID, canteenID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Debug().Model(&entity.Food{}).Where("food_id = ? AND canteen_id = ?", foodID, canteenID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
