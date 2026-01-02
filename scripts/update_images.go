package main

import (
	"context"
	"log"
	"strings"

	"phone-store-backend/internal/config"
	"phone-store-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load config
	cfg := config.Load()
	log.Println("✅ Configuration loaded")

	// Connect to MongoDB
	mongodb, err := db.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()
	log.Println("✅ Connected to MongoDB")

	database := mongodb.Database

	// Update Brands
	if err := updateBrands(database); err != nil {
		log.Printf("❌ Failed to update brands: %v", err)
	} else {
		log.Println("✅ Updated brand logos successfully")
	}

	// Update Products
	if err := updateProducts(database); err != nil {
		log.Printf("❌ Failed to update products: %v", err)
	} else {
		log.Println("✅ Updated product images successfully")
	}

	// Update Categories
	if err := updateCategories(database); err != nil {
		log.Printf("❌ Failed to update categories: %v", err)
	} else {
		log.Println("✅ Updated category images successfully")
	}

	log.Println("🎉 All updates completed!")
}

func updateBrands(db *mongo.Database) error {
	collection := db.Collection("brands")
	ctx := context.Background()

	// Find all brands
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var brand bson.M
		if err := cursor.Decode(&brand); err != nil {
			continue
		}

		logo, ok := brand["logo"].(string)
		if !ok || logo == "" {
			continue
		}

		// Skip if already has prefix
		if strings.HasPrefix(logo, "images/") || strings.HasPrefix(logo, "http") {
			continue
		}

		// Update with prefix
		newLogo := "images/brands/" + logo
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": brand["_id"]},
			bson.M{"$set": bson.M{"logo": newLogo}},
		)
		if err != nil {
			log.Printf("Failed to update brand %s: %v", brand["_id"], err)
			continue
		}
		count++
	}

	log.Printf("Updated %d brand logos", count)
	return nil
}

func updateProducts(db *mongo.Database) error {
	collection := db.Collection("products")
	ctx := context.Background()

	// Find all products
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var product bson.M
		if err := cursor.Decode(&product); err != nil {
			continue
		}

		images, ok := product["images"].(bson.A)
		if !ok || len(images) == 0 {
			continue
		}

		// Update images
		newImages := make([]string, 0)
		hasChanges := false
		for _, img := range images {
			imgStr, ok := img.(string)
			if !ok || imgStr == "" {
				continue
			}

			// Skip if already has prefix
			if strings.HasPrefix(imgStr, "images/") || strings.HasPrefix(imgStr, "http") {
				newImages = append(newImages, imgStr)
				continue
			}

			// Add prefix
			newImages = append(newImages, "images/products/"+imgStr)
			hasChanges = true
		}

		if hasChanges {
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"_id": product["_id"]},
				bson.M{"$set": bson.M{"images": newImages}},
			)
			if err != nil {
				log.Printf("Failed to update product %s: %v", product["_id"], err)
				continue
			}
			count++
		}
	}

	log.Printf("Updated %d product images", count)
	return nil
}

func updateCategories(db *mongo.Database) error {
	collection := db.Collection("categories")
	ctx := context.Background()

	// Find all categories
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var category bson.M
		if err := cursor.Decode(&category); err != nil {
			continue
		}

		image, ok := category["image"].(string)
		if !ok || image == "" {
			continue
		}

		// Skip if already has prefix
		if strings.HasPrefix(image, "images/") || strings.HasPrefix(image, "http") {
			continue
		}

		// Update with prefix
		newImage := "images/categories/" + image
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": category["_id"]},
			bson.M{"$set": bson.M{"image": newImage}},
		)
		if err != nil {
			log.Printf("Failed to update category %s: %v", category["_id"], err)
			continue
		}
		count++
	}

	log.Printf("Updated %d category images", count)
	return nil
}
