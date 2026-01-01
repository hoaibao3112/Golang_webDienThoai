package cart

type AddItemRequest struct {
	VariantID string `json:"variantId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID    string     `json:"id"`
	Items []CartItem `json:"items"`
}

type CartItem struct {
	VariantID   string  `json:"variantId"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	SKU         string  `json:"sku"`
	Color       string  `json:"color"`
	Storage     string  `json:"storage"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Quantity    int     `json:"quantity"`
	Image       string  `json:"image"`
}
