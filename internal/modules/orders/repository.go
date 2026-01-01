package orders

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

func (r *Repository) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := r.db.Collection("orders").InsertOne(ctx, order)
	return err
}

func (r *Repository) CreateOrderItems(ctx context.Context, items []*models.OrderItem) error {
	var docs []interface{}
	for _, item := range items {
		docs = append(docs, item)
	}
	_, err := r.db.Collection("order_items").InsertMany(ctx, docs)
	return err
}

func (r *Repository) FindCartByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Collection("carts").FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	return &cart, err
}

func (r *Repository) FindVariantByID(ctx context.Context, id primitive.ObjectID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Collection("product_variants").FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&variant)
	return &variant, err
}

func (r *Repository) FindProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.db.Collection("products").FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&product)
	return &product, err
}

func (r *Repository) FindVoucherByCode(ctx context.Context, code string) (*models.Voucher, error) {
	var voucher models.Voucher
	err := r.db.Collection("vouchers").FindOne(ctx, bson.M{"code": code, "isActive": true}).Decode(&voucher)
	return &voucher, err
}

// Decrease stock when order is created
func (r *Repository) DecreaseStock(ctx context.Context, variantID primitive.ObjectID, quantity int) error {
	_, err := r.db.Collection("product_variants").UpdateOne(
		ctx,
		bson.M{"_id": variantID},
		bson.M{
			"$inc": bson.M{"stock": -quantity},
			"$set": bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (r *Repository) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.db.Collection("carts").UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{
			"$set": bson.M{
				"items":     []models.CartItem{},
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *Repository) FindOrdersByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Order, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) FindOrderByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	var order models.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	return &order, err
}

func (r *Repository) FindOrderItemsByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*models.OrderItem, error) {
	cursor, err := r.db.Collection("order_items").Find(ctx, bson.M{"orderId": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*models.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}
