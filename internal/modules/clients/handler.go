package clients

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

// GetProfile godoc
// @Summary Get client profile
// @Tags Clients
// @Security BearerAuth
// @Success 200 {object} ClientProfileResponse
// @Router /api/clients/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	resp, err := h.service.GetProfile(c.Request.Context(), userID)
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
		"message": "Profile retrieved successfully",
		"data":    resp,
	})
}

// UpdateProfile godoc
// @Summary Update client profile
// @Tags Clients
// @Security BearerAuth
// @Param request body UpdateClientProfileRequest true "Profile data"
// @Success 200
// @Router /api/clients/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req UpdateClientProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    nil,
	})
}

// GetAllClients godoc
// @Summary Get all clients (admin only)
// @Tags Clients
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} ClientsListResponse
// @Router /api/admin/clients [get]
func (h *Handler) GetAllClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.service.GetAllClients(c.Request.Context(), page, limit)
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
		"message": "Clients retrieved successfully",
		"data":    resp,
	})
}
