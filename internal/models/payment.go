package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID       primitive.ObjectID `bson:"orderId" json:"orderId"`
	Method        string             `bson:"method" json:"method"` // COD, CARD, MOMO, etc.
	Amount        float64            `bson:"amount" json:"amount"`
	Status        PaymentStatus      `bson:"status" json:"status"`
	TransactionID string             `bson:"transactionId,omitempty" json:"transactionId,omitempty"`
	PaidAt        *time.Time         `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
