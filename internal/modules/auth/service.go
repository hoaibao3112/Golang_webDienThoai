package auth

import (
	"context"
	"errors"
	"time"

	"phone-store-backend/internal/config"
	"phone-store-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) error {
	// Check if user already exists
	_, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("email already exists")
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Get default USER role
	role, err := s.repo.FindRoleByName(ctx, "USER")
	if err != nil {
		// If role doesn't exist, create a default ObjectID (in production, ensure roles are seeded)
		role = &models.Role{ID: primitive.NewObjectID()}
	}

	// Create user
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		Phone:     req.Phone,
		RoleID:    role.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Get user role
	role, err := s.repo.FindRoleByID(ctx, user.RoleID)
	roleName := "USER"
	if err == nil {
		roleName = role.Name
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID.Hex(), user.Email, roleName)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User: UserProfile{
			ID:       user.ID.Hex(),
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			Role:     roleName,
		},
	}, nil
}

func (s *Service) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get user role
	role, err := s.repo.FindRoleByID(ctx, user.RoleID)
	roleName := "USER"
	if err == nil {
		roleName = role.Name
	}

	return &UserProfile{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     roleName,
	}, nil
}

func (s *Service) generateToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"role":   role,
		"exp":    time.Now().Add(s.cfg.JWTExpiration).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
