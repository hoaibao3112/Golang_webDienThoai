package reviews

import (
	"context"
	"errors"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo           *Repository
	userCollection *mongo.Collection
}

func NewService(repo *Repository, db *mongo.Database) *Service {
	return &Service{
		repo:           repo,
		userCollection: db.Collection("users"),
	}
}

// GetProductReviews returns all reviews for a product
func (s *Service) GetProductReviews(ctx context.Context, productID string) ([]ReviewResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	reviews, err := s.repo.FindReviewsByProductID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	var response []ReviewResponse
	for _, review := range reviews {
		// Get user name
		var user models.User
		userName := "Anonymous"
		if err := s.userCollection.FindOne(ctx, primitive.M{"_id": review.UserID}).Decode(&user); err == nil {
			userName = user.FullName
		}

		response = append(response, ReviewResponse{
			ID:        review.ID.Hex(),
			ProductID: review.ProductID.Hex(),
			UserID:    review.UserID.Hex(),
			UserName:  userName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
		})
	}

	return response, nil
}

// CreateReview creates a new review
func (s *Service) CreateReview(ctx context.Context, userID string, req *CreateReviewRequest) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	productObjectID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Check if user already reviewed this product
	existing, _ := s.repo.FindReviewByUserAndProduct(ctx, userObjectID, productObjectID)
	if existing != nil {
		return errors.New("you have already reviewed this product")
	}

	review := &models.Review{
		ProductID: productObjectID,
		UserID:    userObjectID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateReview(ctx, review)
}

// DeleteReview deletes a review (admin only)
func (s *Service) DeleteReview(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid review ID")
	}

	// Check if review exists
	_, err = s.repo.FindReviewByID(ctx, objectID)
	if err != nil {
		return errors.New("review not found")
	}

	return s.repo.DeleteReview(ctx, objectID)
}
