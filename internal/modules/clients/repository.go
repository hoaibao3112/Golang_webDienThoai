package clients

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
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		userCollection: db.Collection("users"),
	}
}

// FindClients returns all clients with pagination
func (r *Repository) FindClients(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.User, error) {
	cursor, err := r.userCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clients []*models.User
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, err
	}

	return clients, nil
}

// CountClients counts total clients matching filter
func (r *Repository) CountClients(ctx context.Context, filter bson.M) (int64, error) {
	return r.userCollection.CountDocuments(ctx, filter)
}

// FindClientByID finds a client by ID
func (r *Repository) FindClientByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var client models.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": id, "role": "CLIENT"}).Decode(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// UpdateClient updates a client
func (r *Repository) UpdateClient(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.userCollection.UpdateOne(ctx, bson.M{"_id": id, "role": "CLIENT"}, bson.M{"$set": update})
	return err
}
