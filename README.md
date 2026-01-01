# Phone Store Backend API

Backend API cho ứng dụng Phone Store được xây dựng với Golang, Gin framework và MongoDB.

## 🛠 Tech Stack

- **Language:** Golang 1.22+
- **Framework:** Gin
- **Database:** MongoDB
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **Environment Management:** godotenv

## 📁 Project Structure

```
phone-store-backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── db/
│   │   └── mongo.go             # MongoDB connection & indexes
│   ├── middlewares/
│   │   ├── auth.go              # JWT authentication
│   │   ├── cors.go              # CORS configuration
│   │   ├── error.go             # Error handling
│   │   └── logger.go            # Request logging
│   ├── models/
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── product.go
│   │   ├── product_variant.go
│   │   ├── brand.go
│   │   ├── category.go
│   │   ├── cart.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   ├── payment.go
│   │   ├── review.go
│   │   ├── voucher.go
│   │   └── banner.go
│   └── modules/
│       ├── auth/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   └── dto.go
│       ├── products/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   └── dto.go
│       ├── cart/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   └── dto.go
│       └── orders/
│           ├── handler.go
│           ├── service.go
│           ├── repository.go
│           └── dto.go
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- MongoDB 5.0 or higher
- Git

### Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd phone-store-backend
```

2. **Install Go dependencies:**
```bash
go mod tidy
```

3. **Setup environment variables:**
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:
```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=phone_store
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRATION=24h
CORS_ORIGIN=http://localhost:3000
```

4. **Make sure MongoDB is running:**
```bash
# Check if MongoDB is running
mongosh
```

5. **Run the server:**
```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication

#### Register
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "fullName": "John Doe",
  "phone": "0123456789"
}
```

**Response:**
```json
{
  "message": "User registered successfully"
}
```

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "fullName": "John Doe",
    "phone": "0123456789",
    "role": "USER"
  }
}
```

#### Get Current User
```bash
GET /api/auth/me
Authorization: Bearer <token>
```

### Products (Public)

#### Get All Products
```bash
GET /api/products?search=iphone&brand=apple&category=smartphone&page=1&limit=20&sort=newest
```

**Query Parameters:**
- `search` - Search by product name
- `brand` - Filter by brand slug
- `category` - Filter by category slug
- `minPrice` - Minimum price
- `maxPrice` - Maximum price
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)
- `sort` - Sort by: `newest`, `name_asc`, `name_desc`, `price_asc`, `price_desc`

**Response:**
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "iPhone 15 Pro",
      "slug": "iphone-15-pro",
      "description": "Latest iPhone model",
      "brand": {
        "id": "507f1f77bcf86cd799439012",
        "name": "Apple",
        "slug": "apple",
        "logo": "https://example.com/apple-logo.png"
      },
      "category": {
        "id": "507f1f77bcf86cd799439013",
        "name": "Smartphone",
        "slug": "smartphone",
        "image": "https://example.com/smartphone.jpg"
      },
      "images": ["image1.jpg", "image2.jpg"],
      "minPrice": 999.99,
      "maxPrice": 1299.99,
      "isActive": true,
      "isFeatured": true
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 100,
  "totalPages": 5
}
```

#### Get Product by Slug
```bash
GET /api/products/iphone-15-pro
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "iPhone 15 Pro",
  "slug": "iphone-15-pro",
  "description": "Latest iPhone model",
  "brand": {...},
  "category": {...},
  "images": ["image1.jpg", "image2.jpg"],
  "variants": [
    {
      "id": "507f1f77bcf86cd799439014",
      "sku": "IP15P-BLK-128",
      "color": "Black",
      "storage": "128GB",
      "price": 999.99,
      "stock": 50,
      "isActive": true
    }
  ],
  "isActive": true,
  "isFeatured": true
}
```

#### Get Brands
```bash
GET /api/brands
```

#### Get Categories
```bash
GET /api/categories
```

### Cart (Authenticated)

#### Get Cart
```bash
GET /api/cart
Authorization: Bearer <token>
```

#### Add Item to Cart
```bash
POST /api/cart/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "variantId": "507f1f77bcf86cd799439014",
  "quantity": 2
}
```

#### Update Cart Item
```bash
PUT /api/cart/items/:variantId
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 3
}
```

#### Remove Cart Item
```bash
DELETE /api/cart/items/:variantId
Authorization: Bearer <token>
```

### Orders (Authenticated)

#### Create Order
```bash
POST /api/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "shippingAddress": {
    "fullName": "John Doe",
    "phone": "0123456789",
    "address": "123 Main Street",
    "city": "Ho Chi Minh City",
    "district": "District 1",
    "ward": "Ward 1"
  },
  "voucherCode": "SUMMER2024"
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439015",
  "orderNumber": "ORD-1704123456789",
  "shippingAddress": {...},
  "items": [
    {
      "productId": "507f1f77bcf86cd799439011",
      "variantId": "507f1f77bcf86cd799439014",
      "name": "iPhone 15 Pro",
      "sku": "IP15P-BLK-128",
      "color": "Black",
      "storage": "128GB",
      "price": 999.99,
      "quantity": 2,
      "totalPrice": 1999.98
    }
  ],
  "subTotal": 1999.98,
  "discount": 199.99,
  "total": 1800.00,
  "status": "PENDING",
  "createdAt": "2024-01-01T10:00:00Z"
}
```

#### Get My Orders
```bash
GET /api/orders/me
Authorization: Bearer <token>
```

#### Get Order by ID
```bash
GET /api/orders/:id
Authorization: Bearer <token>
```

### Admin Routes (Admin Only)

#### Create Product
```bash
POST /api/admin/products
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro",
  "slug": "iphone-15-pro",
  "description": "Latest iPhone model",
  "brandId": "507f1f77bcf86cd799439012",
  "categoryId": "507f1f77bcf86cd799439013",
  "images": ["image1.jpg", "image2.jpg"],
  "isFeatured": true
}
```

#### Update Product
```bash
PUT /api/admin/products/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max",
  "description": "Updated description",
  "isActive": true
}
```

#### Delete Product (Soft Delete)
```bash
DELETE /api/admin/products/:id
Authorization: Bearer <admin-token>
```

#### Create Variant
```bash
POST /api/admin/variants
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "productId": "507f1f77bcf86cd799439011",
  "sku": "IP15P-WHT-256",
  "color": "White",
  "storage": "256GB",
  "price": 1099.99,
  "stock": 30
}
```

#### Update Variant
```bash
PUT /api/admin/variants/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "price": 1049.99,
  "stock": 25,
  "isActive": true
}
```

#### Delete Variant (Soft Delete)
```bash
DELETE /api/admin/variants/:id
Authorization: Bearer <admin-token>
```

## 🗄️ Database Collections

### Collections:
- `users` - User accounts
- `roles` - User roles (USER, ADMIN)
- `categories` - Product categories
- `brands` - Product brands
- `products` - Products
- `product_variants` - Product variants (color, storage, price, stock)
- `carts` - Shopping carts
- `orders` - Orders
- `order_items` - Order line items
- `payments` - Payment transactions
- `reviews` - Product reviews
- `vouchers` - Discount vouchers
- `banners` - Homepage banners

### Indexes (Auto-created on startup):
- `users.email` (unique)
- `products.slug` (unique)
- `brands.slug` (unique)
- `categories.slug` (unique)
- `product_variants.sku` (unique)
- `product_variants.productId`
- `reviews.productId`
- `orders.userId`

## 🔒 Security Features

- JWT token-based authentication
- bcrypt password hashing
- CORS protection
- Input validation
- Role-based access control (RBAC)

## 📝 Business Logic

### Order Creation Flow:
1. Validate cart is not empty
2. Check stock availability for all variants
3. Calculate subtotal from cart items
4. Apply voucher discount (if valid)
5. Create order with snapshot prices
6. Create order items
7. **Decrease stock for each variant**
8. Clear user's cart
9. Return order details

### Voucher Application:
- Check if voucher is active
- Check if voucher is not expired
- Check minimum order value
- Calculate discount percentage
- Apply max discount limit

### Stock Management:
- Stock is checked before adding to cart
- Stock is validated again before order creation
- Stock is decreased atomically when order is created
- Soft delete for products/variants (sets `isActive: false`)

## 🧪 Testing

### Sample cURL Commands

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","fullName":"Test User","phone":"0123456789"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**Get Products:**
```bash
curl http://localhost:8080/api/products?page=1&limit=10
```

**Get Product Detail:**
```bash
curl http://localhost:8080/api/products/iphone-15-pro
```

**Add to Cart (with token):**
```bash
curl -X POST http://localhost:8080/api/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"variantId":"507f1f77bcf86cd799439014","quantity":1}'
```

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGO_DB` | Database name | `phone_store` |
| `JWT_SECRET` | JWT signing secret | `default-secret-change-in-production` |
| `JWT_EXPIRATION` | Token expiration | `24h` |
| `CORS_ORIGIN` | Allowed CORS origin | `http://localhost:3000` |

## 📄 Error Response Format

All errors follow this format:
```json
{
  "message": "Error description",
  "code": "ERROR_CODE",
  "details": null
}
```

Common error codes:
- `BAD_REQUEST` - Invalid request data
- `UNAUTHORIZED` - Missing or invalid token
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `INTERNAL_ERROR` - Server error

## 🚀 Production Deployment

1. Set strong `JWT_SECRET` in production
2. Use MongoDB Atlas or production MongoDB
3. Enable MongoDB authentication
4. Set appropriate CORS origins
5. Use environment-specific `.env` files
6. Enable HTTPS/TLS
7. Consider using MongoDB transactions for critical operations (order creation)
8. Implement rate limiting
9. Add logging and monitoring
10. Set up backup strategies

## 👥 Development Team

Backend API developed for Phone Store project.

## 📞 Support

For issues and questions, please create an issue in the repository.

---

**Happy Coding! 🎉**
