package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"freepass-2026/pkg/database/mariadb"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFoodService interface {
	CreateFood(ownerID uuid.UUID, param model.CreateFoodParam) (*model.CreateFoodResponse, error)
	BulkCreateFood(ownerID uuid.UUID, param model.BulkCreateFoodParam) (*model.BulkCreateFoodResponse, error)
	UpdateFood(ownerID uuid.UUID, param model.UpdateFoodParam) (*model.UpdateFoodResponse, error)
	DeleteFood(ownerID uuid.UUID, foodID uuid.UUID) (*model.DeleteFoodResponse, error)
	GetFoodByID(foodID uuid.UUID) (*model.GetFoodResponse, error)
	GetFoodsByCanteenOwner(ownerID uuid.UUID) (*model.GetFoodListResponse, error)
	GetAllFoods() (*model.GetFoodListResponse, error)
	GetAllFoodsGroupedByCanteen() (*model.GetFoodsGroupedByCanteenResponse, error)
	UpdateStock(ownerID uuid.UUID, param model.UpdateStockParam) (*model.UpdateStockResponse, error)
	UpdateCanteenStatus(ownerID uuid.UUID, param model.UpdateCanteenStatusParam) (*model.UpdateCanteenStatusResponse, error)
}

type FoodService struct {
	db              *gorm.DB
	foodRepository  repository.IFoodRepository
	adminRepository repository.IAdminRepository
}

func NewFoodService(foodRepository repository.IFoodRepository, adminRepository repository.IAdminRepository) IFoodService {
	return &FoodService{
		db:              mariadb.Connection,
		foodRepository:  foodRepository,
		adminRepository: adminRepository,
	}
}

func (s *FoodService) CreateFood(ownerID uuid.UUID, param model.CreateFoodParam) (*model.CreateFoodResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	foodID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	food := &entity.Food{
		FoodID:      foodID,
		CanteenID:   canteen.CanteenID,
		FoodName:    param.FoodName,
		Description: param.Description,
		Price:       param.Price,
		Stock:       param.Stock,
		IsAvailable: true,
	}

	err = s.foodRepository.CreateFood(tx, food)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.CreateFoodResponse{
		FoodID:      food.FoodID,
		CanteenID:   food.CanteenID,
		FoodName:    food.FoodName,
		Description: food.Description,
		Price:       food.Price,
		Stock:       food.Stock,
		IsAvailable: food.IsAvailable,
	}

	return response, nil
}

func (s *FoodService) BulkCreateFood(ownerID uuid.UUID, param model.BulkCreateFoodParam) (*model.BulkCreateFoodResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	createdFoods := make([]model.CreateFoodResponse, 0)

	for _, foodParam := range param.Foods {
		foodID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		food := &entity.Food{
			FoodID:      foodID,
			CanteenID:   canteen.CanteenID,
			FoodName:    foodParam.FoodName,
			Description: foodParam.Description,
			Price:       foodParam.Price,
			Stock:       foodParam.Stock,
			IsAvailable: true,
		}

		err = s.foodRepository.CreateFood(tx, food)
		if err != nil {
			return nil, err
		}

		createdFoods = append(createdFoods, model.CreateFoodResponse{
			FoodID:      food.FoodID,
			CanteenID:   food.CanteenID,
			FoodName:    food.FoodName,
			Description: food.Description,
			Price:       food.Price,
			Stock:       food.Stock,
			IsAvailable: food.IsAvailable,
		})
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.BulkCreateFoodResponse{
		Foods: createdFoods,
		Count: len(createdFoods),
	}

	return response, nil
}

func (s *FoodService) UpdateFood(ownerID uuid.UUID, param model.UpdateFoodParam) (*model.UpdateFoodResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	isOwner, err := s.foodRepository.CheckFoodOwnership(param.FoodID, canteen.CanteenID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("you don't have permission to update this food")
	}

	existingFood, err := s.foodRepository.GetFoodByID(param.FoodID)
	if err != nil {
		return nil, errors.New("food not found")
	}

	if param.FoodName != nil {
		existingFood.FoodName = *param.FoodName
	}
	if param.Description != nil {
		existingFood.Description = param.Description
	}
	if param.Price != nil {
		existingFood.Price = *param.Price
	}
	if param.Stock != nil {
		existingFood.Stock = *param.Stock
	}
	if param.IsAvailable != nil {
		existingFood.IsAvailable = *param.IsAvailable
	}

	err = s.foodRepository.UpdateFood(tx, existingFood)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.UpdateFoodResponse{
		FoodID:      existingFood.FoodID,
		CanteenID:   existingFood.CanteenID,
		FoodName:    existingFood.FoodName,
		Description: existingFood.Description,
		Price:       existingFood.Price,
		Stock:       existingFood.Stock,
		IsAvailable: existingFood.IsAvailable,
	}

	return response, nil
}

func (s *FoodService) DeleteFood(ownerID uuid.UUID, foodID uuid.UUID) (*model.DeleteFoodResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	isOwner, err := s.foodRepository.CheckFoodOwnership(foodID, canteen.CanteenID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("you don't have permission to delete this food")
	}

	err = s.foodRepository.DeleteFood(tx, foodID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.DeleteFoodResponse{
		Message: "food deleted successfully",
	}

	return response, nil
}

func (s *FoodService) GetFoodByID(foodID uuid.UUID) (*model.GetFoodResponse, error) {
	food, err := s.foodRepository.GetFoodByID(foodID)
	if err != nil {
		return nil, errors.New("food not found")
	}

	canteen, err := s.adminRepository.GetCanteenByID(food.CanteenID)
	canteenName := ""
	if err == nil && canteen != nil {
		canteenName = canteen.CanteenName
	}

	response := &model.GetFoodResponse{
		FoodID:      food.FoodID,
		CanteenID:   food.CanteenID,
		CanteenName: canteenName,
		FoodName:    food.FoodName,
		Description: food.Description,
		Price:       food.Price,
		Stock:       food.Stock,
		IsAvailable: food.IsAvailable,
	}

	return response, nil
}

func (s *FoodService) GetFoodsByCanteenOwner(ownerID uuid.UUID) (*model.GetFoodListResponse, error) {
	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	foods, err := s.foodRepository.GetFoodsByCanteenID(canteen.CanteenID)
	if err != nil {
		return nil, err
	}

	foodInfos := make([]model.FoodInfo, 0)
	for _, food := range foods {
		foodInfos = append(foodInfos, model.FoodInfo{
			FoodID:      food.FoodID,
			FoodName:    food.FoodName,
			Description: food.Description,
			Price:       food.Price,
			Stock:       food.Stock,
			IsAvailable: food.IsAvailable,
		})
	}

	response := &model.GetFoodListResponse{
		Foods: foodInfos,
	}

	return response, nil
}

func (s *FoodService) GetAllFoods() (*model.GetFoodListResponse, error) {
	foods, err := s.foodRepository.GetAllFoods()
	if err != nil {
		return nil, err
	}

	foodInfos := make([]model.FoodInfo, 0)
	for _, food := range foods {
		foodInfos = append(foodInfos, model.FoodInfo{
			FoodID:      food.FoodID,
			FoodName:    food.FoodName,
			Description: food.Description,
			Price:       food.Price,
			Stock:       food.Stock,
			IsAvailable: food.IsAvailable,
		})
	}

	response := &model.GetFoodListResponse{
		Foods: foodInfos,
	}

	return response, nil
}

func (s *FoodService) GetAllFoodsGroupedByCanteen() (*model.GetFoodsGroupedByCanteenResponse, error) {
	foods, err := s.foodRepository.GetAllFoods()
	if err != nil {
		return nil, err
	}

	canteenMap := make(map[uuid.UUID]*model.CanteenWithFoods)

	for _, food := range foods {
		if _, exists := canteenMap[food.CanteenID]; !exists {
			canteen, err := s.adminRepository.GetCanteenByID(food.CanteenID)
			canteenName := ""
			if err == nil {
				canteenName = canteen.CanteenName
			}

			canteenMap[food.CanteenID] = &model.CanteenWithFoods{
				CanteenID:   food.CanteenID,
				CanteenName: canteenName,
				Foods:       make([]model.FoodInfo, 0),
			}
		}

		canteenMap[food.CanteenID].Foods = append(canteenMap[food.CanteenID].Foods, model.FoodInfo{
			FoodID:      food.FoodID,
			FoodName:    food.FoodName,
			Description: food.Description,
			Price:       food.Price,
			Stock:       food.Stock,
			IsAvailable: food.IsAvailable,
		})
	}

	canteens := make([]model.CanteenWithFoods, 0)
	for _, canteen := range canteenMap {
		canteens = append(canteens, *canteen)
	}

	response := &model.GetFoodsGroupedByCanteenResponse{
		Canteens: canteens,
	}

	return response, nil
}

func (s *FoodService) UpdateStock(ownerID uuid.UUID, param model.UpdateStockParam) (*model.UpdateStockResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	isOwner, err := s.foodRepository.CheckFoodOwnership(param.FoodID, canteen.CanteenID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("you don't have permission to update this food stock")
	}

	err = s.foodRepository.UpdateStock(tx, param.FoodID, param.Stock)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.UpdateStockResponse{
		FoodID: param.FoodID,
		Stock:  param.Stock,
	}

	return response, nil
}

func (s *FoodService) UpdateCanteenStatus(ownerID uuid.UUID, param model.UpdateCanteenStatusParam) (*model.UpdateCanteenStatusResponse, error) {
	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("canteen not found")
		}
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = s.adminRepository.UpdateCanteenStatus(tx, canteen.CanteenID, param.IsOpen)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.UpdateCanteenStatusResponse{
		CanteenID:   canteen.CanteenID,
		CanteenName: canteen.CanteenName,
		IsOpen:      param.IsOpen,
	}

	return response, nil
}
