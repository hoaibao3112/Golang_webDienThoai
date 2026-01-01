package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Voucher struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code            string             `bson:"code" json:"code"`
	Description     string             `bson:"description" json:"description"`
	DiscountPercent float64            `bson:"discountPercent" json:"discountPercent"` // 0-100
	MaxDiscount     float64            `bson:"maxDiscount" json:"maxDiscount"`
	MinOrderValue   float64            `bson:"minOrderValue" json:"minOrderValue"`
	ExpiredAt       time.Time          `bson:"expiredAt" json:"expiredAt"`
	IsActive        bool               `bson:"isActive" json:"isActive"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
