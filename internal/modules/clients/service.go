package clients

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetProfile returns the profile of the current client
func (s *Service) GetProfile(ctx context.Context, userID string) (*ClientProfileResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	client, err := s.repo.FindClientByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	return &ClientProfileResponse{
		ID:       client.ID.Hex(),
		Email:    client.Email,
		FullName: client.FullName,
		Phone:    client.Phone,
		IsActive: client.IsActive,
	}, nil
}

// UpdateProfile updates the profile of the current client
func (s *Service) UpdateProfile(ctx context.Context, userID string, req *UpdateClientProfileRequest) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"updatedAt": time.Now(),
	}

	if req.FullName != "" {
		update["fullName"] = req.FullName
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}

	return s.repo.UpdateClient(ctx, objectID, update)
}

// GetAllClients returns paginated list of all clients (admin only)
func (s *Service) GetAllClients(ctx context.Context, page, limit int) (*ClientsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	filter := bson.M{
		"role": "CLIENT",
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	clients, err := s.repo.FindClients(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountClients(ctx, filter)
	if err != nil {
		return nil, err
	}

	var clientResponses []ClientProfileResponse
	for _, client := range clients {
		clientResponses = append(clientResponses, ClientProfileResponse{
			ID:       client.ID.Hex(),
			Email:    client.Email,
			FullName: client.FullName,
			Phone:    client.Phone,
			IsActive: client.IsActive,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &ClientsListResponse{
		Data:       clientResponses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
