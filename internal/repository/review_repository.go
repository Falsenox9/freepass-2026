package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IReviewRepository interface {
	CreateReview(review *entity.Review) error
	GetReviewsByCanteenID(canteenID uuid.UUID) ([]*entity.Review, error)
	GetReviewByID(reviewID uuid.UUID) (*entity.Review, error)
	DeleteReview(reviewID uuid.UUID) error
	CheckReviewExists(orderID uuid.UUID) (bool, error)
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) IReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(review *entity.Review) error {
	err := r.db.Debug().Create(&review).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ReviewRepository) GetReviewsByCanteenID(canteenID uuid.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review

	err := r.db.Debug().
		Joins("JOIN orders ON orders.order_id = reviews.order_id").
		Where("orders.canteen_id = ?", canteenID).
		Order("reviews.created_at desc").
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) GetReviewByID(reviewID uuid.UUID) (*entity.Review, error) {
	var review entity.Review

	err := r.db.Debug().Where("review_id = ?", reviewID).First(&review).Error
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepository) DeleteReview(reviewID uuid.UUID) error {
	err := r.db.Debug().Where("review_id = ?", reviewID).Delete(&entity.Review{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ReviewRepository) CheckReviewExists(orderID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Debug().Model(&entity.Review{}).Where("order_id = ?", orderID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
