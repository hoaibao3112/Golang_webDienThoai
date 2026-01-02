package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShippingMethod struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Cost        float64            `bson:"cost" json:"cost"`
	EstDays     int                `bson:"estDays" json:"estDays"` // Estimated delivery days
	IsActive    bool               `bson:"isActive" json:"isActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
