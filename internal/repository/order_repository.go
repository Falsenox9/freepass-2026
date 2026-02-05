package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	CreateOrder(tx *gorm.DB, order *entity.Order) error
	CreateOrderItem(tx *gorm.DB, orderItem *entity.OrderItem) error
	GetOrdersByUserID(userID uuid.UUID) ([]*entity.Order, error)
	GetOrderItemsByOrderID(orderID uuid.UUID) ([]*entity.OrderItem, error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *entity.Order) error {
	err := tx.Debug().Create(&order).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) CreateOrderItem(tx *gorm.DB, orderItem *entity.OrderItem) error {
	err := tx.Debug().Create(&orderItem).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uuid.UUID) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := r.db.Debug().Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderItemsByOrderID(orderID uuid.UUID) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem

	err := r.db.Debug().Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}

	return orderItems, nil
}
