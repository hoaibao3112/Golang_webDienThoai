package products

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

// Public endpoints
func (h *Handler) GetProducts(c *gin.Context) {
	var query ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid query parameters",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.service.GetProducts(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	resp, err := h.service.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"code":    "NOT_FOUND",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetBrands(c *gin.Context) {
	brands, err := h.service.GetBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, brands)
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Admin endpoints
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.CreateProduct(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "CREATE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
	})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.UpdateProduct(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "UPDATE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
	})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "DELETE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

func (h *Handler) CreateVariant(c *gin.Context) {
	var req CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.CreateVariant(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "CREATE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Variant created successfully",
	})
}

func (h *Handler) UpdateVariant(c *gin.Context) {
	id := c.Param("id")

	var req UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    "BAD_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.UpdateVariant(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "UPDATE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Variant updated successfully",
	})
}

func (h *Handler) DeleteVariant(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteVariant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    "DELETE_FAILED",
			"details": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Variant deleted successfully",
	})
}

// Brand admin handlers
func (h *Handler) CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateBrand(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Brand created successfully",
		"data":    nil,
	})
}

func (h *Handler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")

	var req UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateBrand(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brand updated successfully",
		"data":    nil,
	})
}

func (h *Handler) DeleteBrand(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteBrand(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Brand deleted successfully",
		"data":    nil,
	})
}

// Category admin handlers
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.CreateCategory(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Category created successfully",
		"data":    nil,
	})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"data":    err.Error(),
		})
		return
	}

	if err := h.service.UpdateCategory(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
		"data":    nil,
	})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
		"data":    nil,
	})
}

