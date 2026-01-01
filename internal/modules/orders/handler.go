package orders

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
