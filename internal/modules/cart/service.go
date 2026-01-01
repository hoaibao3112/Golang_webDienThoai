package cart

import (
	"context"
	"errors"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCart(ctx context.Context, userID string) (*CartResponse, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	cart, err := s.repo.FindCartByUserID(ctx, uid)
	if err == mongo.ErrNoDocuments {
		// Create new cart if doesn't exist
		cart = &models.Cart{
			ID:        primitive.NewObjectID(),
			UserID:    uid,
			Items:     []models.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.CreateCart(ctx, cart); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Transform to response with product details
	items, err := s.transformCartItems(ctx, cart.Items)
	if err != nil {
		return nil, err
	}

	return &CartResponse{
		ID:    cart.ID.Hex(),
		Items: items,
	}, nil
}

func (s *Service) AddItem(ctx context.Context, userID string, req *AddItemRequest) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	variantID, err := primitive.ObjectIDFromHex(req.VariantID)
	if err != nil {
		return errors.New("invalid variant ID")
	}

	// Check if variant exists and has stock
	variant, err := s.repo.FindVariantByID(ctx, variantID)
	if err != nil {
		return errors.New("variant not found")
	}

	if variant.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	// Get or create cart
	cart, err := s.repo.FindCartByUserID(ctx, uid)
	if err == mongo.ErrNoDocuments {
		cart = &models.Cart{
			ID:        primitive.NewObjectID(),
			UserID:    uid,
			Items:     []models.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.CreateCart(ctx, cart); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Check if item already in cart
	found := false
	for i, item := range cart.Items {
		if item.VariantID == variantID {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, models.CartItem{
			VariantID: variantID,
			Quantity:  req.Quantity,
		})
	}

	return s.repo.UpdateCart(ctx, cart.ID, cart.Items)
}

func (s *Service) UpdateItem(ctx context.Context, userID, variantID string, req *UpdateItemRequest) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	vid, err := primitive.ObjectIDFromHex(variantID)
	if err != nil {
		return errors.New("invalid variant ID")
	}

	// Check stock
	variant, err := s.repo.FindVariantByID(ctx, vid)
	if err != nil {
		return errors.New("variant not found")
	}

	if variant.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.repo.FindCartByUserID(ctx, uid)
	if err != nil {
		return errors.New("cart not found")
	}

	// Update quantity
	found := false
	for i, item := range cart.Items {
		if item.VariantID == vid {
			cart.Items[i].Quantity = req.Quantity
			found = true
			break
		}
	}

	if !found {
		return errors.New("item not in cart")
	}

	return s.repo.UpdateCart(ctx, cart.ID, cart.Items)
}

func (s *Service) RemoveItem(ctx context.Context, userID, variantID string) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	vid, err := primitive.ObjectIDFromHex(variantID)
	if err != nil {
		return errors.New("invalid variant ID")
	}

	cart, err := s.repo.FindCartByUserID(ctx, uid)
	if err != nil {
		return errors.New("cart not found")
	}

	// Remove item
	newItems := []models.CartItem{}
	for _, item := range cart.Items {
		if item.VariantID != vid {
			newItems = append(newItems, item)
		}
	}

	return s.repo.UpdateCart(ctx, cart.ID, newItems)
}

func (s *Service) transformCartItems(ctx context.Context, items []models.CartItem) ([]CartItem, error) {
	var result []CartItem

	for _, item := range items {
		variant, err := s.repo.FindVariantByID(ctx, item.VariantID)
		if err != nil {
			continue // Skip if variant not found
		}

		product, err := s.repo.FindProductByID(ctx, variant.ProductID)
		if err != nil {
			continue
		}

		image := ""
		if len(product.Images) > 0 {
			image = product.Images[0]
		}

		result = append(result, CartItem{
			VariantID:   variant.ID.Hex(),
			ProductID:   product.ID.Hex(),
			ProductName: product.Name,
			SKU:         variant.SKU,
			Color:       variant.Color,
			Storage:     variant.Storage,
			Price:       variant.Price,
			Stock:       variant.Stock,
			Quantity:    item.Quantity,
			Image:       image,
		})
	}

	return result, nil
}
