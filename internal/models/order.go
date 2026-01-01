package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipping  OrderStatus = "SHIPPING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCanceled  OrderStatus = "CANCELED"
)

type ShippingAddress struct {
	FullName string `bson:"fullName" json:"fullName"`
	Phone    string `bson:"phone" json:"phone"`
	Address  string `bson:"address" json:"address"`
	City     string `bson:"city" json:"city"`
	District string `bson:"district" json:"district"`
	Ward     string `bson:"ward" json:"ward"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderNumber     string             `bson:"orderNumber" json:"orderNumber"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	ShippingAddress ShippingAddress    `bson:"shippingAddress" json:"shippingAddress"`
	VoucherCode     string             `bson:"voucherCode,omitempty" json:"voucherCode,omitempty"`
	SubTotal        float64            `bson:"subTotal" json:"subTotal"`
	Discount        float64            `bson:"discount" json:"discount"`
	Total           float64            `bson:"total" json:"total"`
	Status          OrderStatus        `bson:"status" json:"status"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
