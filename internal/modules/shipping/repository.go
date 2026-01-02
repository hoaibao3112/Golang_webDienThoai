package shipping

import (
	"context"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	addressCollection *mongo.Collection
	methodCollection  *mongo.Collection
	shipmentCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		addressCollection:  db.Collection("shipping_addresses"),
		methodCollection:   db.Collection("shipping_methods"),
		shipmentCollection: db.Collection("shipments"),
	}
}

// Shipping Address methods
func (r *Repository) FindAddressesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.ShippingAddress, error) {
	cursor, err := r.addressCollection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var addresses []*models.ShippingAddress
	if err := cursor.All(ctx, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *Repository) FindAddressByID(ctx context.Context, id primitive.ObjectID) (*models.ShippingAddress, error) {
	var address models.ShippingAddress
	err := r.addressCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&address)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *Repository) CreateAddress(ctx context.Context, address *models.ShippingAddress) error {
	result, err := r.addressCollection.InsertOne(ctx, address)
	if err != nil {
		return err
	}
	address.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) UpdateAddress(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.addressCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *Repository) UnsetDefaultAddresses(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.addressCollection.UpdateMany(
		ctx,
		bson.M{"userId": userID, "isDefault": true},
		bson.M{"$set": bson.M{"isDefault": false}},
	)
	return err
}

// Shipping Method methods
func (r *Repository) FindAllShippingMethods(ctx context.Context) ([]*models.ShippingMethod, error) {
	cursor, err := r.methodCollection.Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []*models.ShippingMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}

	return methods, nil
}

func (r *Repository) FindShippingMethodByID(ctx context.Context, id primitive.ObjectID) (*models.ShippingMethod, error) {
	var method models.ShippingMethod
	err := r.methodCollection.FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&method)
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// Shipment methods
func (r *Repository) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	result, err := r.shipmentCollection.InsertOne(ctx, shipment)
	if err != nil {
		return err
	}
	shipment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) FindShipmentByOrderID(ctx context.Context, orderID primitive.ObjectID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.shipmentCollection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&shipment)
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *Repository) UpdateShipment(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.shipmentCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
