package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductVariant struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
	SKU       string             `bson:"sku" json:"sku"`
	Color     string             `bson:"color" json:"color"`
	Storage   string             `bson:"storage" json:"storage"`
	Price     float64            `bson:"price" json:"price"`
	Stock     int                `bson:"stock" json:"stock"`
	IsActive  bool               `bson:"isActive" json:"isActive"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
