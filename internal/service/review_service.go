package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"

	"github.com/google/uuid"
)

type IReviewService interface {
	CreateReview(userID uuid.UUID, orderID uuid.UUID, param model.CreateReviewParam) error
	GetCanteenReviews(ownerID uuid.UUID) (*model.GetReviewListResponse, error)
	GetPublicCanteenReviews(canteenID uuid.UUID) (*model.GetPublicReviewListResponse, error)
	DeleteReview(ownerID uuid.UUID, reviewID uuid.UUID) error
}

type ReviewService struct {
	reviewRepository repository.IReviewRepository
	orderRepository  repository.IOrderRepository
	userRepository   repository.IUserRepository
	adminRepository  repository.IAdminRepository
}

func NewReviewService(reviewRepository repository.IReviewRepository, orderRepository repository.IOrderRepository, userRepository repository.IUserRepository, adminRepository repository.IAdminRepository) IReviewService {
	return &ReviewService{
		reviewRepository: reviewRepository,
		orderRepository:  orderRepository,
		userRepository:   userRepository,
		adminRepository:  adminRepository,
	}
}

func (s *ReviewService) CreateReview(userID uuid.UUID, orderID uuid.UUID, param model.CreateReviewParam) error {
	order, err := s.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized to review this order")
	}

	if order.PaymentStatus != "paid" {
		return errors.New("can only review paid orders")
	}

	if order.Status != "completed" {
		return errors.New("can only review completed orders")
	}

	exists, err := s.reviewRepository.CheckReviewExists(orderID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("review already exists for this order")
	}

	reviewID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	review := &entity.Review{
		ReviewID: reviewID,
		OrderID:  orderID,
		UserID:   userID,
		Rating:   param.Rating,
		Comment:  param.Comment,
	}

	err = s.reviewRepository.CreateReview(review)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) GetCanteenReviews(ownerID uuid.UUID) (*model.GetReviewListResponse, error) {
	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("canteen not found for this owner")
	}

	reviews, err := s.reviewRepository.GetReviewsByCanteenID(canteen.CanteenID)
	if err != nil {
		return nil, err
	}

	reviewResponses := make([]model.ReviewResponse, 0)
	totalRating := 0
	for _, review := range reviews {
		userParam := model.UserParam{
			UserID: review.UserID,
		}
		user, err := s.userRepository.GetUser(userParam)
		userName := "Unknown"
		if err == nil && user.FullName != nil {
			userName = *user.FullName
		}

		reviewResponses = append(reviewResponses, model.ReviewResponse{
			ReviewID:  review.ReviewID,
			OrderID:   review.OrderID,
			UserID:    review.UserID,
			UserName:  userName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		totalRating += review.Rating
	}

	overallRating := 0.0
	if len(reviews) > 0 {
		overallRating = float64(totalRating) / float64(len(reviews))
	}

	response := &model.GetReviewListResponse{
		OverallRating: overallRating,
		TotalReviews:  len(reviews),
		Reviews:       reviewResponses,
	}

	return response, nil
}

func (s *ReviewService) DeleteReview(ownerID uuid.UUID, reviewID uuid.UUID) error {
	review, err := s.reviewRepository.GetReviewByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	order, err := s.orderRepository.GetOrderByID(review.OrderID)
	if err != nil {
		return errors.New("order not found")
	}

	canteen, err := s.adminRepository.GetCanteenByOwnerID(ownerID)
	if err != nil {
		return errors.New("canteen not found for this owner")
	}

	if order.CanteenID != canteen.CanteenID {
		return errors.New("unauthorized to delete this review")
	}

	err = s.reviewRepository.DeleteReview(reviewID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) GetPublicCanteenReviews(canteenID uuid.UUID) (*model.GetPublicReviewListResponse, error) {
	reviews, err := s.reviewRepository.GetReviewsByCanteenID(canteenID)
	if err != nil {
		return nil, err
	}

	reviewResponses := make([]model.PublicReviewResponse, 0)
	totalRating := 0
	for _, review := range reviews {
		userParam := model.UserParam{
			UserID: review.UserID,
		}
		user, err := s.userRepository.GetUser(userParam)
		fullName := "Anonymous"
		if err == nil && user.FullName != nil {
			fullName = *user.FullName
		}

		reviewResponses = append(reviewResponses, model.PublicReviewResponse{
			FullName:  fullName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		totalRating += review.Rating
	}

	overallRating := 0.0
	if len(reviews) > 0 {
		overallRating = float64(totalRating) / float64(len(reviews))
	}

	response := &model.GetPublicReviewListResponse{
		OverallRating: overallRating,
		TotalReviews:  len(reviews),
		Reviews:       reviewResponses,
	}

	return response, nil
}
