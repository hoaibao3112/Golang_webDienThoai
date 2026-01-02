package users

import (
	"context"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	userCollection *mongo.Collection
	roleCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		userCollection: db.Collection("users"),
		roleCollection: db.Collection("roles"),
	}
}

// FindUsers returns all users (staff/admin) with pagination
func (r *Repository) FindUsers(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.User, error) {
	cursor, err := r.userCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// CountUsers counts total users matching filter
func (r *Repository) CountUsers(ctx context.Context, filter bson.M) (int64, error) {
	return r.userCollection.CountDocuments(ctx, filter)
}

// FindUserByID finds a user by ID
func (r *Repository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByEmail finds a user by email
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	result, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateUser updates a user
func (r *Repository) UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.userCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

// DeleteUser soft deletes a user (sets isActive to false)
func (r *Repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.userCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isActive": false}})
	return err
}

// FindRoleByName finds a role by name
func (r *Repository) FindRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.roleCollection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
