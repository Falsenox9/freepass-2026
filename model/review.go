package model

import "github.com/google/uuid"

type CreateReviewParam struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required,max=500"`
}

type ReviewResponse struct {
	ReviewID  uuid.UUID `json:"review_id"`
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt string    `json:"created_at"`
}

type GetReviewListResponse struct {
	OverallRating float64          `json:"overall_rating"`
	TotalReviews  int              `json:"total_reviews"`
	Reviews       []ReviewResponse `json:"reviews"`
}

// PublicReviewResponse - Used for public endpoints (hides internal IDs)
type PublicReviewResponse struct {
	FullName  string `json:"full_name"` // User's full name from FullName field
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type GetPublicReviewListResponse struct {
	OverallRating float64                `json:"overall_rating"`
	TotalReviews  int                    `json:"total_reviews"`
	Reviews       []PublicReviewResponse `json:"reviews"`
}
