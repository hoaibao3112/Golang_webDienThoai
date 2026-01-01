package auth

import (
	"context"

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

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindRoleByID(ctx context.Context, roleID primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.db.Collection("roles").FindOne(ctx, bson.M{"_id": roleID}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) FindRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Collection("roles").FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
