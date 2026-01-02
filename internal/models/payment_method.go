package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethod struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Code        string             `bson:"code" json:"code"` // COD, BANK_TRANSFER, MOMO, VNPAY
	Description string             `bson:"description" json:"description"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
