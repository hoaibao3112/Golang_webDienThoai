package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Banner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Image       string             `bson:"image" json:"image"`
	Link        string             `bson:"link" json:"link"`
	Order       int                `bson:"order" json:"order"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
