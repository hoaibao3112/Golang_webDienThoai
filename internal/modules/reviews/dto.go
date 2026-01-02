package reviews

import "time"

// CreateReviewRequest DTO for creating a review
type CreateReviewRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

// ReviewResponse DTO for review response
type ReviewResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	UserID    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
