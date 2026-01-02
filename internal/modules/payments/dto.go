package payments

import "time"

// PaymentMethodResponse DTO for payment method response
type PaymentMethodResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CreatePaymentRequest DTO for creating a payment
type CreatePaymentRequest struct {
	OrderID           string  `json:"orderId" binding:"required"`
	PaymentMethodCode string  `json:"paymentMethodCode" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
}

// PaymentResponse DTO for payment response
type PaymentResponse struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"orderId"`
	PaymentMethod string    `json:"paymentMethod"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paidAt,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
