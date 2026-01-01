package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description" json:"description"`
	BrandID     primitive.ObjectID `bson:"brandId" json:"brandId"`
	CategoryID  primitive.ObjectID `bson:"categoryId" json:"categoryId"`
	Images      []string           `bson:"images" json:"images"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	IsFeatured  bool               `bson:"isFeatured" json:"isFeatured"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
