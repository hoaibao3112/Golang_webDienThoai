package reviews

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

// GetProductReviews godoc
// @Summary Get reviews for a product
// @Tags Reviews
// @Param id path string true "Product ID"
// @Success 200 {array} ReviewResponse
// @Router /api/products/{id}/reviews [get]
func (h *Handler) GetProductReviews(c *gin.Context) {
	productID := c.Param("id")

	reviews, err := h.service.GetProductReviews(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reviews retrieved successfully",
		"data":    reviews,
	})
}

// CreateReview godoc
// @Summary Create a review (client only)
// @Tags Reviews
// @Security BearerAuth
// @Param request body CreateReviewRequest true "Review data"
// @Success 201
// @Router /api/reviews [post]
func (h *Handler) CreateReview(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateReview(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Review created successfully",
		"data":    nil,
	})
}

// DeleteReview godoc
// @Summary Delete a review (admin only)
// @Tags Reviews
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 200
// @Router /api/reviews/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteReview(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review deleted successfully",
		"data":    nil,
	})
}
