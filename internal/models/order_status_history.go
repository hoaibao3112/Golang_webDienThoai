package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatusHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   primitive.ObjectID `bson:"orderId" json:"orderId"`
	Status    string             `bson:"status" json:"status"`
	Note      string             `bson:"note" json:"note"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"` // User ID who made the change
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
