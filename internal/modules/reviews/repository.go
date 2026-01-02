package reviews

import (
	"context"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	reviewCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		reviewCollection: db.Collection("reviews"),
	}
}

// FindReviewsByProductID finds all reviews for a product
func (r *Repository) FindReviewsByProductID(ctx context.Context, productID primitive.ObjectID) ([]*models.Review, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.reviewCollection.Find(ctx, bson.M{"productId": productID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// FindReviewByID finds a review by ID
func (r *Repository) FindReviewByID(ctx context.Context, id primitive.ObjectID) (*models.Review, error) {
	var review models.Review
	err := r.reviewCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// CreateReview creates a new review
func (r *Repository) CreateReview(ctx context.Context, review *models.Review) error {
	result, err := r.reviewCollection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// DeleteReview deletes a review
func (r *Repository) DeleteReview(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.reviewCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindReviewByUserAndProduct checks if user already reviewed this product
func (r *Repository) FindReviewByUserAndProduct(ctx context.Context, userID, productID primitive.ObjectID) (*models.Review, error) {
	var review models.Review
	err := r.reviewCollection.FindOne(ctx, bson.M{
		"userId":    userID,
		"productId": productID,
	}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}
