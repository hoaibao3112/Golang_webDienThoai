package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   primitive.ObjectID `bson:"orderId" json:"orderId"`
	VariantID primitive.ObjectID `bson:"variantId" json:"variantId"`
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
	Name      string             `bson:"name" json:"name"`
	SKU       string             `bson:"sku" json:"sku"`
	Color     string             `bson:"color" json:"color"`
	Storage   string             `bson:"storage" json:"storage"`
	Price     float64            `bson:"price" json:"price"` // Snapshot price at order time
	Quantity  int                `bson:"quantity" json:"quantity"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
