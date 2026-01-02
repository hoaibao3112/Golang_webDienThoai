package products

import (
	"context"
	"errors"
	"math"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProducts(ctx context.Context, query *ProductQuery) (*PaginatedResponse, error) {
	// Build filter
	filter := bson.M{"isActive": true}

	if query.Search != "" {
		filter["name"] = bson.M{"$regex": query.Search, "$options": "i"}
	}

	if query.Brand != "" {
		filter["brandId"] = query.Brand
	}

	if query.Category != "" {
		filter["categoryId"] = query.Category
	}

	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}

	// Build options
	opts := options.Find()
	opts.SetSkip(int64((query.Page - 1) * query.Limit))
	opts.SetLimit(int64(query.Limit))

	// Sorting
	switch query.Sort {
	case "newest":
		opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	case "name_asc":
		opts.SetSort(bson.D{{Key: "name", Value: 1}})
	case "name_desc":
		opts.SetSort(bson.D{{Key: "name", Value: -1}})
	default:
		opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	}

	products, err := s.repo.FindProducts(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Transform to response
	var productResponses []ProductResponse
	for _, product := range products {
		resp, err := s.transformProduct(ctx, product)
		if err != nil {
			continue
		}

		// Get variants to calculate price range
		variants, _ := s.repo.FindVariantsByProductID(ctx, product.ID)
		if len(variants) > 0 {
			minPrice, maxPrice := s.calculatePriceRange(variants)
			resp.MinPrice = minPrice
			resp.MaxPrice = maxPrice
		}

		productResponses = append(productResponses, *resp)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &PaginatedResponse{
		Data:       productResponses,
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetProductBySlug(ctx context.Context, slug string) (*ProductDetailResponse, error) {
	product, err := s.repo.FindProductBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("product not found")
	}

	resp, err := s.transformProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	// Get variants
	variants, err := s.repo.FindVariantsByProductID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	var variantResponses []VariantResponse
	for _, v := range variants {
		variantResponses = append(variantResponses, VariantResponse{
			ID:       v.ID.Hex(),
			SKU:      v.SKU,
			Color:    v.Color,
			Storage:  v.Storage,
			Price:    v.Price,
			Stock:    v.Stock,
			IsActive: v.IsActive,
		})
	}

	return &ProductDetailResponse{
		ProductResponse: *resp,
		Variants:        variantResponses,
	}, nil
}

func (s *Service) GetBrands(ctx context.Context) ([]*Brand, error) {
	brands, err := s.repo.FindAllBrands(ctx)
	if err != nil {
		return nil, err
	}

	var response []*Brand
	for _, b := range brands {
		response = append(response, &Brand{
			ID:   b.ID.Hex(),
			Name: b.Name,
			Slug: b.Slug,
			Logo: b.Logo,
		})
	}
	return response, nil
}

func (s *Service) GetCategories(ctx context.Context) ([]*Category, error) {
	categories, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	var response []*Category
	for _, c := range categories {
		response = append(response, &Category{
			ID:    c.ID.Hex(),
			Name:  c.Name,
			Slug:  c.Slug,
			Image: c.Image,
		})
	}
	return response, nil
}

// Admin methods
func (s *Service) CreateProduct(ctx context.Context, req *CreateProductRequest) error {
	brandID, err := primitive.ObjectIDFromHex(req.BrandID)
	if err != nil {
		return errors.New("invalid brand ID")
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return errors.New("invalid category ID")
	}

	product := &models.Product{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		BrandID:     brandID,
		CategoryID:  categoryID,
		Images:      req.Images,
		IsActive:    true,
		IsFeatured:  req.IsFeatured,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateProduct(ctx, product)
}

func (s *Service) UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest) error {
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	update := bson.M{"updatedAt": time.Now()}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.BrandID != "" {
		brandID, err := primitive.ObjectIDFromHex(req.BrandID)
		if err != nil {
			return errors.New("invalid brand ID")
		}
		update["brandId"] = brandID
	}
	if req.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return errors.New("invalid category ID")
		}
		update["categoryId"] = categoryID
	}
	if req.Images != nil {
		update["images"] = req.Images
	}
	update["isFeatured"] = req.IsFeatured
	update["isActive"] = req.IsActive

	return s.repo.UpdateProduct(ctx, productID, update)
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}
	return s.repo.DeleteProduct(ctx, productID)
}

func (s *Service) CreateVariant(ctx context.Context, req *CreateVariantRequest) error {
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	variant := &models.ProductVariant{
		ID:        primitive.NewObjectID(),
		ProductID: productID,
		SKU:       req.SKU,
		Color:     req.Color,
		Storage:   req.Storage,
		Price:     req.Price,
		Stock:     req.Stock,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateVariant(ctx, variant)
}

func (s *Service) UpdateVariant(ctx context.Context, id string, req *UpdateVariantRequest) error {
	variantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid variant ID")
	}

	update := bson.M{"updatedAt": time.Now()}

	if req.Color != "" {
		update["color"] = req.Color
	}
	if req.Storage != "" {
		update["storage"] = req.Storage
	}
	if req.Price > 0 {
		update["price"] = req.Price
	}
	if req.Stock >= 0 {
		update["stock"] = req.Stock
	}
	update["isActive"] = req.IsActive

	return s.repo.UpdateVariant(ctx, variantID, update)
}

func (s *Service) DeleteVariant(ctx context.Context, id string) error {
	variantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid variant ID")
	}
	return s.repo.DeleteVariant(ctx, variantID)
}

// Helper methods
func (s *Service) transformProduct(ctx context.Context, product *models.Product) (*ProductResponse, error) {
	brand, _ := s.repo.FindBrandByID(ctx, product.BrandID)
	category, _ := s.repo.FindCategoryByID(ctx, product.CategoryID)

	resp := &ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Images:      product.Images,
		IsActive:    product.IsActive,
		IsFeatured:  product.IsFeatured,
	}

	if brand != nil {
		resp.Brand = &Brand{
			ID:   brand.ID.Hex(),
			Name: brand.Name,
			Slug: brand.Slug,
			Logo: brand.Logo,
		}
	}

	if category != nil {
		resp.Category = &Category{
			ID:    category.ID.Hex(),
			Name:  category.Name,
			Slug:  category.Slug,
			Image: category.Image,
		}
	}

	return resp, nil
}

func (s *Service) calculatePriceRange(variants []*models.ProductVariant) (float64, float64) {
	if len(variants) == 0 {
		return 0, 0
	}

	minPrice := variants[0].Price
	maxPrice := variants[0].Price

	for _, v := range variants {
		if v.Price < minPrice {
			minPrice = v.Price
		}
		if v.Price > maxPrice {
			maxPrice = v.Price
		}
	}

	return minPrice, maxPrice
}

// Brand admin methods
func (s *Service) CreateBrand(ctx context.Context, req *CreateBrandRequest) error {
	// Check if slug already exists
	existing, _ := s.repo.FindBrandBySlug(ctx, req.Slug)
	if existing != nil {
		return errors.New("brand slug already exists")
	}

	brand := &models.Brand{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Slug:      req.Slug,
		Logo:      req.Logo,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateBrand(ctx, brand)
}

func (s *Service) UpdateBrand(ctx context.Context, id string, req *UpdateBrandRequest) error {
	brandID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid brand ID")
	}

	update := bson.M{"updatedAt": time.Now()}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Logo != "" {
		update["logo"] = req.Logo
	}
	if req.IsActive != nil {
		update["isActive"] = *req.IsActive
	}

	return s.repo.UpdateBrand(ctx, brandID, update)
}

func (s *Service) DeleteBrand(ctx context.Context, id string) error {
	brandID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid brand ID")
	}
	return s.repo.DeleteBrand(ctx, brandID)
}

// Category admin methods
func (s *Service) CreateCategory(ctx context.Context, req *CreateCategoryRequest) error {
	// Check if slug already exists
	existing, _ := s.repo.FindCategoryBySlug(ctx, req.Slug)
	if existing != nil {
		return errors.New("category slug already exists")
	}

	category := &models.Category{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Slug:      req.Slug,
		Image:     req.Image,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateCategory(ctx, category)
}

func (s *Service) UpdateCategory(ctx context.Context, id string, req *UpdateCategoryRequest) error {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID")
	}

	update := bson.M{"updatedAt": time.Now()}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Image != "" {
		update["image"] = req.Image
	}
	if req.IsActive != nil {
		update["isActive"] = *req.IsActive
	}

	return s.repo.UpdateCategory(ctx, categoryID, update)
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID")
	}
	return s.repo.DeleteCategory(ctx, categoryID)
}

