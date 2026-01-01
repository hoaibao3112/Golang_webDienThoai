package auth

import (
	"net/http"

	"phone-store-backend/internal/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Register(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "REGISTRATION_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
			"code":    "LOGIN_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"code":    "UNAUTHORIZED",
			"details": nil,
		})
		return
	}

	profile, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"code":    "USER_NOT_FOUND",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func RegisterRoutes(router *gin.RouterGroup, db interface{}, cfg *config.Config) {
	repo := NewRepository(db.(*mongo.Database))
	service := NewService(repo, cfg)
	handler := NewHandler(service)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
	}
}
