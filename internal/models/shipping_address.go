package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShippingAddress struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	FullName  string             `bson:"fullName" json:"fullName"`
	Phone     string             `bson:"phone" json:"phone"`
	Address   string             `bson:"address" json:"address"`
	City      string             `bson:"city" json:"city"`
	District  string             `bson:"district" json:"district"`
	Ward      string             `bson:"ward" json:"ward"`
	IsDefault bool               `bson:"isDefault" json:"isDefault"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
