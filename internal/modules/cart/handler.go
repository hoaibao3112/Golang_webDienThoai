package cart

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

func (h *Handler) GetCart(c *gin.Context) {
	userID, _ := c.Get("userID")

	cart, err := h.service.GetCart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddItem(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.AddItem(c.Request.Context(), userID.(string), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "ADD_ITEM_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
	})
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	variantID := c.Param("variantId")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.UpdateItem(c.Request.Context(), userID.(string), variantID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "UPDATE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item updated",
	})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	variantID := c.Param("variantId")

	if err := h.service.RemoveItem(c.Request.Context(), userID.(string), variantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "REMOVE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart",
	})
}
