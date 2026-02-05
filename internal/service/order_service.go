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

type IOrderService interface {
	CreateOrder(userID uuid.UUID, param model.CreateOrderParam) (*model.CreateOrderResponse, error)
	GetUserOrders(userID uuid.UUID) (*model.GetOrderListResponse, error)
}

type OrderService struct {
	db              *gorm.DB
	orderRepository repository.IOrderRepository
	foodRepository  repository.IFoodRepository
	adminRepository repository.IAdminRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, foodRepository repository.IFoodRepository, adminRepository repository.IAdminRepository) IOrderService {
	return &OrderService{
		db:              mariadb.Connection,
		orderRepository: orderRepository,
		foodRepository:  foodRepository,
		adminRepository: adminRepository,
	}
}

func (s *OrderService) CreateOrder(userID uuid.UUID, param model.CreateOrderParam) (*model.CreateOrderResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if len(param.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	var canteenID uuid.UUID
	totalPrice := 0
	orderItems := make([]*entity.OrderItem, 0)
	orderItemResponses := make([]model.OrderItemResponse, 0)

	for i, item := range param.Items {
		food, err := s.foodRepository.GetFoodByID(item.FoodID)
		if err != nil {
			return nil, errors.New("food not found")
		}

		if !food.IsAvailable {
			return nil, errors.New("food is not available")
		}

		if food.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for " + food.FoodName)
		}

		if i == 0 {
			canteenID = food.CanteenID
		} else {
			if food.CanteenID != canteenID {
				return nil, errors.New("all items must be from the same canteen")
			}
		}

		subtotal := food.Price * item.Quantity
		totalPrice += subtotal

		orderItemID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		orderItem := &entity.OrderItem{
			OrderItemID: orderItemID,
			FoodID:      item.FoodID,
			Quantity:    item.Quantity,
			Price:       food.Price,
			Subtotal:    subtotal,
		}

		orderItems = append(orderItems, orderItem)

		orderItemResponses = append(orderItemResponses, model.OrderItemResponse{
			OrderItemID: orderItemID,
			FoodID:      item.FoodID,
			FoodName:    food.FoodName,
			Quantity:    item.Quantity,
			Price:       food.Price,
			Subtotal:    subtotal,
		})
	}

	orderID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	order := &entity.Order{
		OrderID:    orderID,
		UserID:     userID,
		CanteenID:  canteenID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	err = s.orderRepository.CreateOrder(tx, order)
	if err != nil {
		return nil, err
	}

	for i, orderItem := range orderItems {
		orderItem.OrderID = orderID

		err = s.orderRepository.CreateOrderItem(tx, orderItem)
		if err != nil {
			return nil, err
		}

		newStock := 0
		food, _ := s.foodRepository.GetFoodByID(orderItem.FoodID)
		newStock = food.Stock - param.Items[i].Quantity

		err = s.foodRepository.UpdateStock(tx, orderItem.FoodID, newStock)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.CreateOrderResponse{
		OrderID:    order.OrderID,
		CanteenID:  order.CanteenID,
		Items:      orderItemResponses,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}

	return response, nil
}

func (s *OrderService) GetUserOrders(userID uuid.UUID) (*model.GetOrderListResponse, error) {
	orders, err := s.orderRepository.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]model.GetOrderResponse, 0)

	for _, order := range orders {
		orderItems, err := s.orderRepository.GetOrderItemsByOrderID(order.OrderID)
		if err != nil {
			return nil, err
		}

		orderItemResponses := make([]model.OrderItemResponse, 0)
		for _, item := range orderItems {
			food, err := s.foodRepository.GetFoodByID(item.FoodID)
			foodName := ""
			if err == nil {
				foodName = food.FoodName
			}

			orderItemResponses = append(orderItemResponses, model.OrderItemResponse{
				OrderItemID: item.OrderItemID,
				FoodID:      item.FoodID,
				FoodName:    foodName,
				Quantity:    item.Quantity,
				Price:       item.Price,
				Subtotal:    item.Subtotal,
			})
		}

		canteen, err := s.adminRepository.GetCanteenByID(order.CanteenID)
		canteenName := ""
		if err == nil {
			canteenName = canteen.CanteenName
		}

		orderResponses = append(orderResponses, model.GetOrderResponse{
			OrderID:     order.OrderID,
			CanteenID:   order.CanteenID,
			CanteenName: canteenName,
			Items:       orderItemResponses,
			TotalPrice:  order.TotalPrice,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response := &model.GetOrderListResponse{
		Orders: orderResponses,
	}

	return response, nil
}
