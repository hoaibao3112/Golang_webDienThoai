package shipping

import (
	"context"
	"errors"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Address methods
func (s *Service) GetUserAddresses(ctx context.Context, userID string) ([]AddressResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	addresses, err := s.repo.FindAddressesByUserID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	var response []AddressResponse
	for _, addr := range addresses {
		response = append(response, AddressResponse{
			ID:        addr.ID.Hex(),
			FullName:  addr.FullName,
			Phone:     addr.Phone,
			Address:   addr.Address,
			City:      addr.City,
			District:  addr.District,
			Ward:      addr.Ward,
			IsDefault: addr.IsDefault,
		})
	}

	return response, nil
}

func (s *Service) CreateAddress(ctx context.Context, userID string, req *CreateAddressRequest) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// If this is default, unset other default addresses
	if req.IsDefault {
		s.repo.UnsetDefaultAddresses(ctx, objectID)
	}

	address := &models.ShippingAddress{
		UserID:    objectID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		District:  req.District,
		Ward:      req.Ward,
		IsDefault: req.IsDefault,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateAddress(ctx, address)
}

func (s *Service) UpdateAddress(ctx context.Context, userID, addressID string, req *UpdateAddressRequest) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	addressObjectID, err := primitive.ObjectIDFromHex(addressID)
	if err != nil {
		return errors.New("invalid address ID")
	}

	// Verify address belongs to user
	address, err := s.repo.FindAddressByID(ctx, addressObjectID)
	if err != nil {
		return errors.New("address not found")
	}

	if address.UserID != userObjectID {
		return errors.New("unauthorized")
	}

	update := bson.M{"updatedAt": time.Now()}

	if req.FullName != "" {
		update["fullName"] = req.FullName
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.City != "" {
		update["city"] = req.City
	}
	if req.District != "" {
		update["district"] = req.District
	}
	if req.Ward != "" {
		update["ward"] = req.Ward
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			s.repo.UnsetDefaultAddresses(ctx, userObjectID)
		}
		update["isDefault"] = *req.IsDefault
	}

	return s.repo.UpdateAddress(ctx, addressObjectID, update)
}

// Shipping Method methods
func (s *Service) GetShippingMethods(ctx context.Context) ([]ShippingMethodResponse, error) {
	methods, err := s.repo.FindAllShippingMethods(ctx)
	if err != nil {
		return nil, err
	}

	var response []ShippingMethodResponse
	for _, method := range methods {
		response = append(response, ShippingMethodResponse{
			ID:          method.ID.Hex(),
			Name:        method.Name,
			Description: method.Description,
			Cost:        method.Cost,
			EstDays:     method.EstDays,
		})
	}

	return response, nil
}

// Shipment methods
func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) error {
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	shippingMethodID, err := primitive.ObjectIDFromHex(req.ShippingMethodID)
	if err != nil {
		return errors.New("invalid shipping method ID")
	}

	shipment := &models.Shipment{
		OrderID:          orderID,
		ShippingMethodID: shippingMethodID,
		TrackingNumber:   req.TrackingNumber,
		Status:           "PENDING",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return s.repo.CreateShipment(ctx, shipment)
}

func (s *Service) UpdateShipmentStatus(ctx context.Context, orderID, status string) error {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	shipment, err := s.repo.FindShipmentByOrderID(ctx, objectID)
	if err != nil {
		return errors.New("shipment not found")
	}

	update := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}

	if status == "SHIPPING" {
		now := time.Now()
		update["shippedAt"] = now
	} else if status == "DELIVERED" {
		now := time.Now()
		update["deliveredAt"] = now
	}

	return s.repo.UpdateShipment(ctx, shipment.ID, update)
}
