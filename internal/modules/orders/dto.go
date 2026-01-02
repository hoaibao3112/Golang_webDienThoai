package orders

import "phone-store-backend/internal/models"

type CreateOrderRequest struct {
	ShippingAddress ShippingAddressRequest `json:"shippingAddress" binding:"required"`
	VoucherCode     string                 `json:"voucherCode"`
}

type ShippingAddressRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	District string `json:"district" binding:"required"`
	Ward     string `json:"ward" binding:"required"`
}

type OrderResponse struct {
	ID              string                      `json:"id"`
	OrderNumber     string                      `json:"orderNumber"`
	ShippingAddress models.OrderShippingAddress `json:"shippingAddress"`
	Items           []OrderItemResponse         `json:"items"`
	SubTotal        float64                     `json:"subTotal"`
	Discount        float64                     `json:"discount"`
	Total           float64                     `json:"total"`
	Status          string                      `json:"status"`
	CreatedAt       string                      `json:"createdAt"`
}

type OrderItemResponse struct {
	ProductID  string  `json:"productId"`
	VariantID  string  `json:"variantId"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Color      string  `json:"color"`
	Storage    string  `json:"storage"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}

// UpdateOrderStatusRequest DTO for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

// StatusHistoryResponse DTO for order status history
type StatusHistoryResponse struct {
	Status    string `json:"status"`
	Note      string `json:"note"`
	UpdatedBy string `json:"updatedBy"`
	CreatedAt string `json:"createdAt"`
}

// OrdersListResponse DTO for paginated orders list
type OrdersListResponse struct {
	Data       []*OrderResponse `json:"data"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
}
