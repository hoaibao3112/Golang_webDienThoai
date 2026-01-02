package payments

import (
	"context"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	paymentCollection       *mongo.Collection
	paymentMethodCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		paymentCollection:       db.Collection("payments"),
		paymentMethodCollection: db.Collection("payment_methods"),
	}
}

// FindAllPaymentMethods returns all active payment methods
func (r *Repository) FindAllPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	cursor, err := r.paymentMethodCollection.Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []*models.PaymentMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}

	return methods, nil
}

// FindPaymentByOrderID finds payment by order ID
func (r *Repository) FindPaymentByOrderID(ctx context.Context, orderID primitive.ObjectID) (*models.Payment, error) {
	var payment models.Payment
	err := r.paymentCollection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// CreatePayment creates a new payment
func (r *Repository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	result, err := r.paymentCollection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}
	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdatePayment updates payment status
func (r *Repository) UpdatePayment(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.paymentCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

// FindPaymentMethodByCode finds payment method by code
func (r *Repository) FindPaymentMethodByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.paymentMethodCollection.FindOne(ctx, bson.M{"code": code, "isActive": true}).Decode(&method)
	if err != nil {
		return nil, err
	}
	return &method, nil
}
