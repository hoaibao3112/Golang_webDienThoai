package products

type ProductQuery struct {
	Search   string  `form:"search"`
	Brand    string  `form:"brand"`
	Category string  `form:"category"`
	MinPrice float64 `form:"minPrice"`
	MaxPrice float64 `form:"maxPrice"`
	Page     int     `form:"page"`
	Limit    int     `form:"limit"`
	Sort     string  `form:"sort"` // price_asc, price_desc, name_asc, name_desc, newest
}

type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Brand       *Brand   `json:"brand"`
	Category    *Category `json:"category"`
	Images      []string `json:"images"`
	MinPrice    float64  `json:"minPrice"`
	MaxPrice    float64  `json:"maxPrice"`
	IsActive    bool     `json:"isActive"`
	IsFeatured  bool     `json:"isFeatured"`
}

type ProductDetailResponse struct {
	ProductResponse
	Variants []VariantResponse `json:"variants"`
}

type VariantResponse struct {
	ID       string  `json:"id"`
	SKU      string  `json:"sku"`
	Color    string  `json:"color"`
	Storage  string  `json:"storage"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	IsActive bool    `json:"isActive"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Logo string `json:"logo"`
}

type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"totalPages"`
}

type CreateProductRequest struct {
	Name        string               `json:"name" binding:"required"`
	Slug        string               `json:"slug" binding:"required"`
	Description string               `json:"description"`
	BrandID     string               `json:"brandId" binding:"required"`
	CategoryID  string               `json:"categoryId" binding:"required"`
	Images      []string             `json:"images"`
	IsFeatured  bool                 `json:"isFeatured"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BrandID     string   `json:"brandId"`
	CategoryID  string   `json:"categoryId"`
	Images      []string `json:"images"`
	IsFeatured  bool     `json:"isFeatured"`
	IsActive    bool     `json:"isActive"`
}

type CreateVariantRequest struct {
	ProductID string  `json:"productId" binding:"required"`
	SKU       string  `json:"sku" binding:"required"`
	Color     string  `json:"color" binding:"required"`
	Storage   string  `json:"storage" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Stock     int     `json:"stock" binding:"required"`
}

type UpdateVariantRequest struct {
	Color    string  `json:"color"`
	Storage  string  `json:"storage"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	IsActive bool    `json:"isActive"`
}

// Brand DTOs
type CreateBrandRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Logo string `json:"logo"`
}

type UpdateBrandRequest struct {
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	IsActive *bool  `json:"isActive"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
	Image string `json:"image"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	IsActive *bool  `json:"isActive"`
}

