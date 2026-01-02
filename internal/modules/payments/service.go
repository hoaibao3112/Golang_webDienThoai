package payments

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

// GetPaymentMethods returns all active payment methods
func (s *Service) GetPaymentMethods(ctx context.Context) ([]PaymentMethodResponse, error) {
	methods, err := s.repo.FindAllPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	var response []PaymentMethodResponse
	for _, method := range methods {
		response = append(response, PaymentMethodResponse{
			ID:          method.ID.Hex(),
			Name:        method.Name,
			Code:        method.Code,
			Description: method.Description,
		})
	}

	return response, nil
}

// CreatePayment creates a payment for an order
func (s *Service) CreatePayment(ctx context.Context, req *CreatePaymentRequest) error {
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	// Get payment method
	paymentMethod, err := s.repo.FindPaymentMethodByCode(ctx, req.PaymentMethodCode)
	if err != nil {
		return errors.New("payment method not found")
	}

	// Create payment
	payment := &models.Payment{
		OrderID:         orderID,
		PaymentMethodID: paymentMethod.ID,
		Amount:          req.Amount,
		Status:          "PENDING", // PENDING, COMPLETED, FAILED
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// If COD, payment is pending until delivered
	// Other methods would need payment gateway integration

	return s.repo.CreatePayment(ctx, payment)
}

// GetPaymentByOrderID returns payment information for an order
func (s *Service) GetPaymentByOrderID(ctx context.Context, orderID string) (*PaymentResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	payment, err := s.repo.FindPaymentByOrderID(ctx, objectID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Get payment method name
	var paymentMethodName string
	if paymentMethod, err := s.repo.FindPaymentMethodByCode(ctx, ""); err == nil {
		paymentMethodName = paymentMethod.Name
	}

	response := &PaymentResponse{
		ID:            payment.ID.Hex(),
		OrderID:       payment.OrderID.Hex(),
		PaymentMethod: paymentMethodName,
		Amount:        payment.Amount,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt,
	}

	if payment.PaidAt != nil {
		response.PaidAt = *payment.PaidAt
	}

	return response, nil
}

// UpdatePaymentStatus updates payment status (for admin/system)
func (s *Service) UpdatePaymentStatus(ctx context.Context, orderID, status string) error {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	payment, err := s.repo.FindPaymentByOrderID(ctx, objectID)
	if err != nil {
		return errors.New("payment not found")
	}

	update := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}

	if status == "COMPLETED" {
		now := time.Now()
		update["paidAt"] = now
	}

	return s.repo.UpdatePayment(ctx, payment.ID, update)
}
