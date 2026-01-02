package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shipment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID          primitive.ObjectID `bson:"orderId" json:"orderId"`
	ShippingMethodID primitive.ObjectID `bson:"shippingMethodId" json:"shippingMethodId"`
	TrackingNumber   string             `bson:"trackingNumber" json:"trackingNumber"`
	Status           string             `bson:"status" json:"status"` // PENDING, SHIPPING, DELIVERED
	ShippedAt        *time.Time         `bson:"shippedAt,omitempty" json:"shippedAt,omitempty"`
	DeliveredAt      *time.Time         `bson:"deliveredAt,omitempty" json:"deliveredAt,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
