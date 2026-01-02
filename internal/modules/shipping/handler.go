package shipping

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

// Address handlers
// GetAddresses godoc
// @Summary Get user's shipping addresses
// @Tags Shipping
// @Security BearerAuth
// @Success 200 {array} AddressResponse
// @Router /api/shipping-addresses [get]
func (h *Handler) GetAddresses(c *gin.Context) {
	userID := c.GetString("userID")

	addresses, err := h.service.GetUserAddresses(c.Request.Context(), userID)
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
		"message": "Addresses retrieved successfully",
		"data":    addresses,
	})
}

// CreateAddress godoc
// @Summary Create a shipping address
// @Tags Shipping
// @Security BearerAuth
// @Param request body CreateAddressRequest true "Address data"
// @Success 201
// @Router /api/shipping-addresses [post]
func (h *Handler) CreateAddress(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateAddress(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Address created successfully",
		"data":    nil,
	})
}

// UpdateAddress godoc
// @Summary Update a shipping address
// @Tags Shipping
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body UpdateAddressRequest true "Address data"
// @Success 200
// @Router /api/shipping-addresses/{id} [put]
func (h *Handler) UpdateAddress(c *gin.Context) {
	userID := c.GetString("userID")
	addressID := c.Param("id")

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateAddress(c.Request.Context(), userID, addressID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address updated successfully",
		"data":    nil,
	})
}

// Shipping Method handlers
// GetShippingMethods godoc
// @Summary Get all shipping methods
// @Tags Shipping
// @Success 200 {array} ShippingMethodResponse
// @Router /api/shipping-methods [get]
func (h *Handler) GetShippingMethods(c *gin.Context) {
	methods, err := h.service.GetShippingMethods(c.Request.Context())
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
		"message": "Shipping methods retrieved successfully",
		"data":    methods,
	})
}

// Shipment handlers (admin only)
// CreateShipment godoc
// @Summary Create a shipment (admin only)
// @Tags Shipping
// @Security BearerAuth
// @Param request body CreateShipmentRequest true "Shipment data"
// @Success 201
// @Router /api/shipments [post]
func (h *Handler) CreateShipment(c *gin.Context) {
	var req CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateShipment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Shipment created successfully",
		"data":    nil,
	})
}
