package shipping

import "time"

// Address DTOs
type CreateAddressRequest struct {
	FullName  string `json:"fullName" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Ward      string `json:"ward" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

type UpdateAddressRequest struct {
	FullName  string `json:"fullName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	District  string `json:"district"`
	Ward      string `json:"ward"`
	IsDefault *bool  `json:"isDefault"`
}

type AddressResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"fullName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	District  string `json:"district"`
	Ward      string `json:"ward"`
	IsDefault bool   `json:"isDefault"`
}

// Shipping Method DTOs
type ShippingMethodResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	EstDays     int     `json:"estDays"`
}

// Shipment DTOs
type CreateShipmentRequest struct {
	OrderID          string `json:"orderId" binding:"required"`
	ShippingMethodID string `json:"shippingMethodId" binding:"required"`
	TrackingNumber   string `json:"trackingNumber"`
}

type ShipmentResponse struct {
	ID             string     `json:"id"`
	OrderID        string     `json:"orderId"`
	TrackingNumber string     `json:"trackingNumber"`
	Status         string     `json:"status"`
	ShippedAt      *time.Time `json:"shippedAt,omitempty"`
	DeliveredAt    *time.Time `json:"deliveredAt,omitempty"`
}
