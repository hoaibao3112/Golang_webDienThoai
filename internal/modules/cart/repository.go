package cart

import (
	"context"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindCartByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Collection("carts").FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *Repository) CreateCart(ctx context.Context, cart *models.Cart) error {
	_, err := r.db.Collection("carts").InsertOne(ctx, cart)
	return err
}

func (r *Repository) UpdateCart(ctx context.Context, cartID primitive.ObjectID, items []models.CartItem) error {
	_, err := r.db.Collection("carts").UpdateOne(
		ctx,
		bson.M{"_id": cartID},
		bson.M{
			"$set": bson.M{
				"items":     items,
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *Repository) FindVariantByID(ctx context.Context, id primitive.ObjectID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Collection("product_variants").FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&variant)
	return &variant, err
}

func (r *Repository) FindProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.db.Collection("products").FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&product)
	return &product, err
}

func (r *Repository) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.db.Collection("carts").UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{
			"$set": bson.M{
				"items":     []models.CartItem{},
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}
