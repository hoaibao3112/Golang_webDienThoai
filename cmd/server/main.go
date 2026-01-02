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
	"phone-store-backend/internal/modules/clients"
	"phone-store-backend/internal/modules/orders"
	"phone-store-backend/internal/modules/payments"
	"phone-store-backend/internal/modules/products"
	"phone-store-backend/internal/modules/reviews"
	"phone-store-backend/internal/modules/shipping"
	"phone-store-backend/internal/modules/users"

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
		// Client profile
		clientRepo := clients.NewRepository(mongodb.Database)
		clientService := clients.NewService(clientRepo)
		clientHandler := clients.NewHandler(clientService)

		clientGroup := protected.Group("/clients")
		{
			clientGroup.GET("/profile", clientHandler.GetProfile)
			clientGroup.PUT("/profile", clientHandler.UpdateProfile)
		}

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
			orderGroup.GET("/:id/history", orderHandler.GetOrderHistory)
		}

		// Review routes
		reviewRepo := reviews.NewRepository(mongodb.Database)
		reviewService := reviews.NewService(reviewRepo, mongodb.Database)
		reviewHandler := reviews.NewHandler(reviewService)

		protected.POST("/reviews", reviewHandler.CreateReview)

		// Shipping addresses
		shippingRepo := shipping.NewRepository(mongodb.Database)
		shippingService := shipping.NewService(shippingRepo)
		shippingHandler := shipping.NewHandler(shippingService)

		shippingGroup := protected.Group("/shipping-addresses")
		{
			shippingGroup.GET("", shippingHandler.GetAddresses)
			shippingGroup.POST("", shippingHandler.CreateAddress)
			shippingGroup.PUT("/:id", shippingHandler.UpdateAddress)
		}
	}

	// Public routes (no auth required)
	api.GET("/products/:id/reviews", func(c *gin.Context) {
		reviewRepo := reviews.NewRepository(mongodb.Database)
		reviewService := reviews.NewService(reviewRepo, mongodb.Database)
		reviewHandler := reviews.NewHandler(reviewService)
		reviewHandler.GetProductReviews(c)
	})

	// Payment methods (public)
	paymentRepo := payments.NewRepository(mongodb.Database)
	paymentService := payments.NewService(paymentRepo)
	paymentHandler := payments.NewHandler(paymentService)

	api.GET("/payment-methods", paymentHandler.GetPaymentMethods)
	api.GET("/shipping-methods", func(c *gin.Context) {
		shippingRepo := shipping.NewRepository(mongodb.Database)
		shippingService := shipping.NewService(shippingRepo)
		shippingHandler := shipping.NewHandler(shippingService)
		shippingHandler.GetShippingMethods(c)
	})

	protected.POST("/payments", paymentHandler.CreatePayment)
	protected.GET("/payments/:orderId", paymentHandler.GetPaymentByOrderID)

	// Admin routes (require authentication + admin role)
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(cfg))
	admin.Use(middlewares.AdminOnly())
	{
		// User management (staff/admin)
		userRepo := users.NewRepository(mongodb.Database)
		userService := users.NewService(userRepo)
		userHandler := users.NewHandler(userService)

		adminUsers := admin.Group("/users")
		{
			adminUsers.GET("", userHandler.GetUsers)
			adminUsers.GET("/:id", userHandler.GetUserByID)
			adminUsers.POST("", userHandler.CreateUser)
			adminUsers.PUT("/:id", userHandler.UpdateUser)
			adminUsers.DELETE("/:id", userHandler.DeleteUser)
		}

		// Client management
		clientRepo := clients.NewRepository(mongodb.Database)
		clientService := clients.NewService(clientRepo)
		clientHandler := clients.NewHandler(clientService)

		admin.GET("/clients", clientHandler.GetAllClients)

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

		// Brand management
		adminBrands := admin.Group("/brands")
		{
			adminBrands.POST("", productHandler.CreateBrand)
			adminBrands.PUT("/:id", productHandler.UpdateBrand)
			adminBrands.DELETE("/:id", productHandler.DeleteBrand)
		}

		// Category management
		adminCategories := admin.Group("/categories")
		{
			adminCategories.POST("", productHandler.CreateCategory)
			adminCategories.PUT("/:id", productHandler.UpdateCategory)
			adminCategories.DELETE("/:id", productHandler.DeleteCategory)
		}

		// Review management
		reviewRepo := reviews.NewRepository(mongodb.Database)
		reviewService := reviews.NewService(reviewRepo, mongodb.Database)
		reviewHandler := reviews.NewHandler(reviewService)

		admin.DELETE("/reviews/:id", reviewHandler.DeleteReview)

		// Order management
		orderRepo := orders.NewRepository(mongodb.Database)
		orderService := orders.NewService(orderRepo)
		orderHandler := orders.NewHandler(orderService)

		admin.GET("/orders", orderHandler.GetAllOrders)
		admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)

		// Shipment management
		shippingRepo := shipping.NewRepository(mongodb.Database)
		shippingService := shipping.NewService(shippingRepo)
		shippingHandler := shipping.NewHandler(shippingService)

		admin.POST("/shipments", shippingHandler.CreateShipment)
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
