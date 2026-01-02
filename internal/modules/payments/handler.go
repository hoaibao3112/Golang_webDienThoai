package payments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetPaymentMethods godoc
// @Summary Get all payment methods
// @Tags Payments
// @Success 200 {array} PaymentMethodResponse
// @Router /api/payment-methods [get]
func (h *Handler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.service.GetPaymentMethods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment methods retrieved successfully",
		"data":    methods,
	})
}

// CreatePayment godoc
// @Summary Create a payment
// @Tags Payments
// @Security BearerAuth
// @Param request body CreatePaymentRequest true "Payment data"
// @Success 201
// @Router /api/payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreatePayment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Payment created successfully",
		"data":    nil,
	})
}

// GetPaymentByOrderID godoc
// @Summary Get payment by order ID
// @Tags Payments
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Success 200 {object} PaymentResponse
// @Router /api/payments/{orderId} [get]
func (h *Handler) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")

	payment, err := h.service.GetPaymentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment retrieved successfully",
		"data":    payment,
	})
}
