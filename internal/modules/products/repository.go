package products

import (
	"context"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindProducts(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.Product, error) {
	cursor, err := r.db.Collection("products").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) CountProducts(ctx context.Context, filter bson.M) (int64, error) {
	return r.db.Collection("products").CountDocuments(ctx, filter)
}

func (r *Repository) FindProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.Collection("products").FindOne(ctx, bson.M{"slug": slug, "isActive": true}).Decode(&product)
	return &product, err
}

func (r *Repository) FindProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.db.Collection("products").FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return &product, err
}

func (r *Repository) FindVariantsByProductID(ctx context.Context, productID primitive.ObjectID) ([]*models.ProductVariant, error) {
	cursor, err := r.db.Collection("product_variants").Find(ctx, bson.M{"productId": productID, "isActive": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var variants []*models.ProductVariant
	if err := cursor.All(ctx, &variants); err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *Repository) FindBrandByID(ctx context.Context, id primitive.ObjectID) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.Collection("brands").FindOne(ctx, bson.M{"_id": id}).Decode(&brand)
	return &brand, err
}

func (r *Repository) FindCategoryByID(ctx context.Context, id primitive.ObjectID) (*models.Category, error) {
	var category models.Category
	err := r.db.Collection("categories").FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	return &category, err
}

func (r *Repository) FindAllBrands(ctx context.Context) ([]*models.Brand, error) {
	cursor, err := r.db.Collection("brands").Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var brands []*models.Brand
	if err := cursor.All(ctx, &brands); err != nil {
		return nil, err
	}
	return brands, nil
}

func (r *Repository) FindAllCategories(ctx context.Context) ([]*models.Category, error) {
	cursor, err := r.db.Collection("categories").Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// Admin methods
func (r *Repository) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.Collection("products").InsertOne(ctx, product)
	return err
}

func (r *Repository) UpdateProduct(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.db.Collection("products").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	// Soft delete
	_, err := r.db.Collection("products").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isActive": false, "updatedAt": time.Now()}})
	return err
}

func (r *Repository) CreateVariant(ctx context.Context, variant *models.ProductVariant) error {
	_, err := r.db.Collection("product_variants").InsertOne(ctx, variant)
	return err
}

func (r *Repository) FindVariantByID(ctx context.Context, id primitive.ObjectID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Collection("product_variants").FindOne(ctx, bson.M{"_id": id}).Decode(&variant)
	return &variant, err
}

func (r *Repository) UpdateVariant(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.db.Collection("product_variants").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteVariant(ctx context.Context, id primitive.ObjectID) error {
	// Soft delete
	_, err := r.db.Collection("product_variants").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isActive": false, "updatedAt": time.Now()}})
	return err
}
