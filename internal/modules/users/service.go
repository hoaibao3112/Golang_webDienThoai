package users

import (
	"context"
	"errors"
	"math"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetUsers returns paginated list of users (staff/admin only)
func (s *Service) GetUsers(ctx context.Context, page, limit int) (*UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// Only get staff and admin users, not clients
	filter := bson.M{
		"role": bson.M{"$in": []string{"ADMIN", "STAFF"}},
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	users, err := s.repo.FindUsers(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       user.ID.Hex(),
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &UsersListResponse{
		Data:       userResponses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID returns a user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.FindUserByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

// CreateUser creates a new staff or admin user
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) error {
	// Check if email already exists
	existingUser, _ := s.repo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Validate role
	if req.Role != "ADMIN" && req.Role != "STAFF" {
		return errors.New("invalid role, must be ADMIN or STAFF")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		Phone:     req.Phone,
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateUser(ctx, user)
}

// UpdateUser updates a user
func (s *Service) UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	user, err := s.repo.FindUserByID(ctx, objectID)
	if err != nil {
		return errors.New("user not found")
	}

	// Prevent changing CLIENT role
	if user.Role == "CLIENT" {
		return errors.New("cannot update client user through this endpoint")
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
	if req.Role != "" {
		if req.Role != "ADMIN" && req.Role != "STAFF" {
			return errors.New("invalid role, must be ADMIN or STAFF")
		}
		update["role"] = req.Role
	}
	if req.IsActive != nil {
		update["isActive"] = *req.IsActive
	}

	return s.repo.UpdateUser(ctx, objectID, update)
}

// DeleteUser soft deletes a user
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	user, err := s.repo.FindUserByID(ctx, objectID)
	if err != nil {
		return errors.New("user not found")
	}

	// Prevent deleting CLIENT users
	if user.Role == "CLIENT" {
		return errors.New("cannot delete client user through this endpoint")
	}

	return s.repo.DeleteUser(ctx, objectID)
}
