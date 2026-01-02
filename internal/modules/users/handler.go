package users

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

// GetUsers godoc
// @Summary Get all users (staff/admin)
// @Tags Users
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} UsersListResponse
// @Router /api/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.service.GetUsers(c.Request.Context(), page, limit)
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
		"message": "Users retrieved successfully",
		"data":    resp,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Router /api/users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.service.GetUserByID(c.Request.Context(), id)
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
		"message": "User retrieved successfully",
		"data":    resp,
	})
}

// CreateUser godoc
// @Summary Create a new user (staff/admin)
// @Tags Users
// @Security BearerAuth
// @Param request body CreateUserRequest true "User data"
// @Success 201
// @Router /api/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateUser(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    nil,
	})
}

// UpdateUser godoc
// @Summary Update a user
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User data"
// @Success 200
// @Router /api/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateUser(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    nil,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200
// @Router /api/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
		"data":    nil,
	})
}
