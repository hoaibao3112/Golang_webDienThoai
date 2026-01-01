package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"phone-store-backend/internal/config"
	"phone-store-backend/internal/db"
	"phone-store-backend/internal/middlewares"
	"phone-store-backend/internal/modules/auth"
	"phone-store-backend/internal/modules/cart"
	"phone-store-backend/internal/modules/orders"
	"phone-store-backend/internal/modules/products"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("✅ Configuration loaded")

	// Connect to MongoDB
	mongodb, err := db.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	// Initialize Gin
	router := gin.Default()

	// Apply middlewares
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORSMiddleware(cfg.CORSOrigin))
	router.Use(middlewares.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	{
		authRepo := auth.NewRepository(mongodb.Database)
		authService := auth.NewService(authRepo, cfg)
		authHandler := auth.NewHandler(authService)

		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middlewares.AuthMiddleware(cfg), authHandler.GetMe)
	}

	// Public product routes
	productRepo := products.NewRepository(mongodb.Database)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)

	api.GET("/products", productHandler.GetProducts)
	api.GET("/products/:slug", productHandler.GetProductBySlug)
	api.GET("/brands", productHandler.GetBrands)
	api.GET("/categories", productHandler.GetCategories)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		// Cart routes
		cartRepo := cart.NewRepository(mongodb.Database)
		cartService := cart.NewService(cartRepo)
		cartHandler := cart.NewHandler(cartService)

		cartGroup := protected.Group("/cart")
		{
			cartGroup.GET("", cartHandler.GetCart)
			cartGroup.POST("/items", cartHandler.AddItem)
			cartGroup.PUT("/items/:variantId", cartHandler.UpdateItem)
			cartGroup.DELETE("/items/:variantId", cartHandler.RemoveItem)
		}

		// Order routes
		orderRepo := orders.NewRepository(mongodb.Database)
		orderService := orders.NewService(orderRepo)
		orderHandler := orders.NewHandler(orderService)

		orderGroup := protected.Group("/orders")
		{
			orderGroup.POST("", orderHandler.CreateOrder)
			orderGroup.GET("/me", orderHandler.GetMyOrders)
			orderGroup.GET("/:id", orderHandler.GetOrderByID)
		}
	}

	// Admin routes (require authentication + admin role)
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(cfg))
	admin.Use(middlewares.RequireAdmin())
	{
		// Product management
		adminProducts := admin.Group("/products")
		{
			adminProducts.POST("", productHandler.CreateProduct)
			adminProducts.PUT("/:id", productHandler.UpdateProduct)
			adminProducts.DELETE("/:id", productHandler.DeleteProduct)
		}

		// Variant management
		adminVariants := admin.Group("/variants")
		{
			adminVariants.POST("", productHandler.CreateVariant)
			adminVariants.PUT("/:id", productHandler.UpdateVariant)
			adminVariants.DELETE("/:id", productHandler.DeleteVariant)
		}
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")
}
