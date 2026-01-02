package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"phone-store-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*OrderResponse, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user's cart
	cart, err := s.repo.FindCartByUserID(ctx, uid)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate stock and calculate total
	var orderItems []*models.OrderItem
	var subTotal float64

	for _, cartItem := range cart.Items {
		variant, err := s.repo.FindVariantByID(ctx, cartItem.VariantID)
		if err != nil {
			return nil, fmt.Errorf("variant %s not found", cartItem.VariantID.Hex())
		}

		// Check stock availability
		if variant.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s", variant.SKU)
		}

		product, err := s.repo.FindProductByID(ctx, variant.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found for variant %s", variant.SKU)
		}

		// Create order item with snapshot price
		orderItem := &models.OrderItem{
			ID:        primitive.NewObjectID(),
			VariantID: variant.ID,
			ProductID: product.ID,
			Name:      product.Name,
			SKU:       variant.SKU,
			Color:     variant.Color,
			Storage:   variant.Storage,
			Price:     variant.Price,
			Quantity:  cartItem.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		orderItems = append(orderItems, orderItem)
		subTotal += variant.Price * float64(cartItem.Quantity)
	}

	// Apply voucher if provided
	discount := 0.0
	if req.VoucherCode != "" {
		voucher, err := s.repo.FindVoucherByCode(ctx, req.VoucherCode)
		if err == nil && voucher.IsActive && time.Now().Before(voucher.ExpiredAt) {
			// Calculate discount
			if subTotal >= voucher.MinOrderValue {
				discount = subTotal * (voucher.DiscountPercent / 100)
				if discount > voucher.MaxDiscount {
					discount = voucher.MaxDiscount
				}
			}
		}
	}

	total := subTotal - discount

	// Create order
	order := &models.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: s.generateOrderNumber(),
		UserID:      uid,
		ShippingAddress: models.OrderShippingAddress{
			FullName: req.ShippingAddress.FullName,
			Phone:    req.ShippingAddress.Phone,
			Address:  req.ShippingAddress.Address,
			City:     req.ShippingAddress.City,
			District: req.ShippingAddress.District,
			Ward:     req.ShippingAddress.Ward,
		},
		VoucherCode: req.VoucherCode,
		SubTotal:    subTotal,
		Discount:    discount,
		Total:       total,
		Status:      models.OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save order
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// Set order ID for items
	for _, item := range orderItems {
		item.OrderID = order.ID
	}

	// Save order items
	if err := s.repo.CreateOrderItems(ctx, orderItems); err != nil {
		return nil, err
	}

	// Decrease stock for each variant
	for _, item := range orderItems {
		if err := s.repo.DecreaseStock(ctx, item.VariantID, item.Quantity); err != nil {
			// Log error but don't fail the order
			// In production, you might want to use transactions
			fmt.Printf("Warning: Failed to decrease stock for variant %s: %v\n", item.VariantID.Hex(), err)
		}
	}

	// Clear cart
	if err := s.repo.ClearCart(ctx, uid); err != nil {
		// Log error but don't fail the order
		fmt.Printf("Warning: Failed to clear cart: %v\n", err)
	}

	// Return order response
	return s.transformOrder(order, orderItems), nil
}

func (s *Service) GetMyOrders(ctx context.Context, userID string) ([]*OrderResponse, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	orders, err := s.repo.FindOrdersByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	var response []*OrderResponse
	for _, order := range orders {
		items, err := s.repo.FindOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			continue
		}
		response = append(response, s.transformOrder(order, items))
	}

	return response, nil
}

func (s *Service) GetOrderByID(ctx context.Context, userID, orderID string, isAdmin bool) (*OrderResponse, error) {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.repo.FindOrderByID(ctx, oid)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check ownership
	if !isAdmin && order.UserID.Hex() != userID {
		return nil, errors.New("access denied")
	}

	items, err := s.repo.FindOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return s.transformOrder(order, items), nil
}

func (s *Service) generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano()/1000000)
}

func (s *Service) transformOrder(order *models.Order, items []*models.OrderItem) *OrderResponse {
	var itemResponses []OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, OrderItemResponse{
			ProductID:  item.ProductID.Hex(),
			VariantID:  item.VariantID.Hex(),
			Name:       item.Name,
			SKU:        item.SKU,
			Color:      item.Color,
			Storage:    item.Storage,
			Price:      item.Price,
			Quantity:   item.Quantity,
			TotalPrice: item.Price * float64(item.Quantity),
		})
	}

	return &OrderResponse{
		ID:              order.ID.Hex(),
		OrderNumber:     order.OrderNumber,
		ShippingAddress: order.ShippingAddress,
		Items:           itemResponses,
		SubTotal:        order.SubTotal,
		Discount:        order.Discount,
		Total:           order.Total,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
	}
}

// UpdateOrderStatus updates order status and logs history
func (s *Service) UpdateOrderStatus(ctx context.Context, orderID, status, note, updatedBy string) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	userID, err := primitive.ObjectIDFromHex(updatedBy)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Validate status
	validStatuses := []string{"PENDING", "CONFIRMED", "PROCESSING", "SHIPPING", "DELIVERED", "CANCELLED"}
	isValid := false
	for _, s := range validStatuses {
		if s == status {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid order status")
	}

	// Update order status
	if err := s.repo.UpdateOrderStatus(ctx, oid, status); err != nil {
		return err
	}

	// Log status history
	history := &models.OrderStatusHistory{
		ID:        primitive.NewObjectID(),
		OrderID:   oid,
		Status:    status,
		Note:      note,
		UpdatedBy: userID,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateStatusHistory(ctx, history)
}

// GetOrderStatusHistory retrieves order status history
func (s *Service) GetOrderStatusHistory(ctx context.Context, orderID string) ([]StatusHistoryResponse, error) {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	history, err := s.repo.FindStatusHistoryByOrderID(ctx, oid)
	if err != nil {
		return nil, err
	}

	var response []StatusHistoryResponse
	for _, h := range history {
		response = append(response, StatusHistoryResponse{
			Status:    h.Status,
			Note:      h.Note,
			UpdatedBy: h.UpdatedBy.Hex(),
			CreatedAt: h.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

// GetAllOrders returns all orders (admin only) with pagination
func (s *Service) GetAllOrders(ctx context.Context, page, limit int, status string) (*OrdersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	filter := primitive.M{}
	if status != "" {
		filter["status"] = status
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	opts.SetSort(primitive.D{{Key: "createdAt", Value: -1}})

	orders, err := s.repo.FindAllOrders(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountOrders(ctx, filter)
	if err != nil {
		return nil, err
	}

	var orderResponses []*OrderResponse
	for _, order := range orders {
		items, err := s.repo.FindOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			continue
		}
		orderResponses = append(orderResponses, s.transformOrder(order, items))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &OrdersListResponse{
		Data:       orderResponses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
