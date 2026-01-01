# Phone Store Frontend - Next.js

Frontend cho ứng dụng Phone Store được xây dựng với Next.js 14, TypeScript và Tailwind CSS.

## 🚀 Tech Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **State Management:** Zustand
- **HTTP Client:** Axios
- **Form Handling:** React Hook Form
- **UI Icons:** Lucide React
- **Notifications:** React Hot Toast

## 📁 Cấu trúc thư mục

```
frontend-nextjs/
├── app/                        # App Router (Next.js 14+)
│   ├── layout.tsx              # Root layout
│   ├── page.tsx                # Homepage
│   ├── globals.css             # Global styles
│   ├── products/               # Products pages
│   │   ├── page.tsx            # Product list
│   │   └── [slug]/             # Product detail
│   │       └── page.tsx
│   ├── cart/                   # Shopping cart
│   │   └── page.tsx
│   ├── checkout/               # Checkout process
│   │   └── page.tsx
│   ├── orders/                 # Order history
│   │   ├── page.tsx
│   │   └── [id]/
│   │       └── page.tsx
│   ├── auth/                   # Authentication
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   └── admin/                  # Admin dashboard
│       └── page.tsx
│
├── components/                 # Reusable components
│   ├── layout/                 # Layout components
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   ├── home/                   # Home page components
│   │   ├── HeroSection.tsx
│   │   ├── FeaturedProducts.tsx
│   │   ├── BrandList.tsx
│   │   └── CategoryList.tsx
│   ├── products/               # Product components
│   │   ├── ProductCard.tsx
│   │   ├── ProductList.tsx
│   │   ├── ProductFilter.tsx
│   │   └── ProductDetail.tsx
│   ├── cart/                   # Cart components
│   │   └── CartItem.tsx
│   └── common/                 # Common components
│       ├── Button.tsx
│       ├── Input.tsx
│       └── Loading.tsx
│
├── services/                   # API services
│   ├── authService.ts          # Authentication API
│   ├── productService.ts       # Products API
│   ├── cartService.ts          # Cart API
│   └── orderService.ts         # Orders API
│
├── store/                      # State management
│   ├── authStore.ts            # Auth state (Zustand)
│   └── cartStore.ts            # Cart state (Zustand)
│
├── lib/                        # Utilities
│   ├── api.ts                  # Axios instance & interceptors
│   └── utils.ts                # Helper functions
│
├── types/                      # TypeScript types
│   └── index.ts                # All type definitions
│
├── public/                     # Static assets
│   └── images/
│
└── config files
    ├── package.json
    ├── tsconfig.json
    ├── tailwind.config.ts
    ├── next.config.js
    └── .env.local.example
```

## 🛠️ Cài đặt và chạy

### 1. Cài đặt dependencies:

```bash
cd frontend-nextjs
npm install
```

### 2. Tạo file .env.local:

```bash
cp .env.local.example .env.local
```

Nội dung file `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### 3. Chạy development server:

```bash
npm run dev
```

Ứng dụng sẽ chạy tại: **http://localhost:3000**

### 4. Build for production:

```bash
npm run build
npm start
```

## 📱 Các trang chính

### Public Pages:
- **/** - Trang chủ (Hero, Featured Products, Categories, Brands)
- **/products** - Danh sách sản phẩm (Filter, Sort, Pagination)
- **/products/[slug]** - Chi tiết sản phẩm (Variants, Add to cart)
- **/auth/login** - Đăng nhập
- **/auth/register** - Đăng ký

### Protected Pages (Cần đăng nhập):
- **/cart** - Giỏ hàng
- **/checkout** - Thanh toán
- **/orders** - Lịch sử đơn hàng
- **/orders/[id]** - Chi tiết đơn hàng
- **/profile** - Thông tin tài khoản

### Admin Pages (Cần role ADMIN):
- **/admin** - Dashboard
- **/admin/products** - Quản lý sản phẩm
- **/admin/orders** - Quản lý đơn hàng
- **/admin/users** - Quản lý người dùng

## 🔐 Authentication Flow

1. User đăng ký/đăng nhập
2. Backend trả về JWT token + user info
3. Token được lưu trong localStorage và Zustand store
4. Axios interceptor tự động thêm token vào header
5. Nếu token expired (401), tự động logout và redirect về login

## 🛒 Shopping Flow

1. **Browse Products** → Xem sản phẩm, filter, search
2. **Product Detail** → Chọn variant (màu, dung lượng)
3. **Add to Cart** → Thêm vào giỏ hàng (local state)
4. **Cart** → Xem giỏ, cập nhật số lượng, xóa item
5. **Checkout** → Nhập thông tin giao hàng, mã voucher
6. **Order Confirmation** → Tạo đơn hàng, trừ stock
7. **Order Tracking** → Theo dõi trạng thái đơn hàng

## 🎨 Components chính

### Layout Components:
- **Header** - Logo, Search, Cart icon, User menu, Navigation
- **Footer** - Company info, Links, Social media

### Home Components:
- **HeroSection** - Banner chính
- **FeaturedProducts** - Sản phẩm nổi bật
- **CategoryList** - Danh mục
- **BrandList** - Thương hiệu

### Product Components:
- **ProductCard** - Card hiển thị sản phẩm
- **ProductList** - Danh sách sản phẩm có pagination
- **ProductFilter** - Bộ lọc (brand, category, price)
- **ProductDetail** - Chi tiết sản phẩm với variants

### Cart Components:
- **CartItem** - Item trong giỏ hàng
- **CartSummary** - Tổng tiền, voucher

## 📦 State Management (Zustand)

### Auth Store:
```typescript
- user: User | null
- token: string | null
- setAuth(user, token)
- clearAuth()
- isAuthenticated()
- isAdmin()
```

### Cart Store:
```typescript
- items: CartItem[]
- addItem(item)
- updateQuantity(variantId, quantity)
- removeItem(variantId)
- clearCart()
- getTotalItems()
- getTotalPrice()
```

## 🌐 API Integration

Tất cả API calls đều thông qua Axios instance với:
- Base URL từ env
- Automatic token injection
- Error handling
- Response/Request interceptors

## 🎯 Features

✅ Server-side rendering (SSR) với Next.js App Router
✅ TypeScript cho type safety
✅ Responsive design (Mobile-first)
✅ Authentication & Authorization
✅ Shopping cart với local state
✅ Product filtering & search
✅ Image optimization với Next/Image
✅ Form validation với React Hook Form
✅ Toast notifications
✅ Loading states & skeletons
✅ Error handling
✅ SEO friendly

## 🚧 Todo (Tính năng mở rộng)

- [ ] Product reviews & ratings
- [ ] Wishlist
- [ ] Product comparison
- [ ] Payment integration
- [ ] Real-time order tracking
- [ ] Push notifications
- [ ] Admin analytics dashboard
- [ ] Multi-language support
- [ ] Dark mode
- [ ] PWA support

## 📱 Responsive Breakpoints

- Mobile: < 768px
- Tablet: 768px - 1024px
- Desktop: > 1024px

## 🎨 Design System

### Colors:
- Primary: Blue (#3b82f6)
- Success: Green
- Error: Red
- Warning: Yellow

### Typography:
- Font: Inter (Google Fonts)
- Headings: Bold
- Body: Regular

---

**Happy Coding! 🚀**
