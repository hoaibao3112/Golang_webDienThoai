package orders

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "CREATE_ORDER_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("userID")

	orders, err := h.service.GetMyOrders(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrderByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	orderID := c.Param("id")

	isAdmin := role == "ADMIN"

	order, err := h.service.GetOrderByID(c.Request.Context(), userID.(string), orderID, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"code":    "NOT_FOUND",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus godoc
// @Summary Update order status (admin only)
// @Tags Orders
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param request body UpdateOrderStatusRequest true "Status data"
// @Success 200
// @Router /api/orders/{id}/status [put]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	userID := c.GetString("userID")
	orderID := c.Param("id")

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateOrderStatus(c.Request.Context(), orderID, req.Status, req.Note, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order status updated successfully",
		"data":    nil,
	})
}

// GetOrderHistory godoc
// @Summary Get order status history
// @Tags Orders
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {array} StatusHistoryResponse
// @Router /api/orders/{id}/history [get]
func (h *Handler) GetOrderHistory(c *gin.Context) {
	orderID := c.Param("id")

	history, err := h.service.GetOrderStatusHistory(c.Request.Context(), orderID)
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
		"message": "Order history retrieved successfully",
		"data":    history,
	})
}

// GetAllOrders godoc
// @Summary Get all orders (admin only)
// @Tags Orders
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param status query string false "Order status"
// @Success 200 {object} OrdersListResponse
// @Router /api/admin/orders [get]
func (h *Handler) GetAllOrders(c *gin.Context) {
	page := 1
	limit := 20
	status := c.Query("status")

	if p, ok := c.GetQuery("page"); ok {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}

	if l, ok := c.GetQuery("limit"); ok {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	orders, err := h.service.GetAllOrders(c.Request.Context(), page, limit, status)
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
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

