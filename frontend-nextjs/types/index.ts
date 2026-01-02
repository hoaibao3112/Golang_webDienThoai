// Types for API responses
export interface User {
  id: string
  email: string
  fullName: string
  phone: string
  role: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Brand {
  id: string
  name: string
  slug: string
  logo: string
}

export interface Category {
  id: string
  name: string
  slug: string
  image: string
}

export interface Product {
  id: string
  name: string
  slug: string
  description: string
  brand: Brand
  category: Category
  images: string[]
  minPrice: number
  maxPrice: number
  originalPrice?: number
  discount?: number
  isNew?: boolean
  isActive: boolean
  isFeatured: boolean
  variants?: ProductVariant[]
}

export interface ProductVariant {
  id: string
  sku: string
  color: string
  storage: string
  price: number
  stock: number
  isActive: boolean
}

export interface ProductDetail extends Product {
  variants: ProductVariant[]
}

export interface CartItem {
  variantId: string
  productId: string
  productName: string
  sku: string
  color: string
  storage: string
  price: number
  stock: number
  quantity: number
  image: string
}

export interface Cart {
  id: string
  items: CartItem[]
}

export interface ShippingAddress {
  fullName: string
  phone: string
  address: string
  city: string
  district: string
  ward: string
}

export interface OrderItem {
  productId: string
  variantId: string
  name: string
  sku: string
  color: string
  storage: string
  price: number
  quantity: number
  totalPrice: number
  image?: string
  productName?: string
  variantName?: string
}

export interface StatusHistory {
  status: string
  note: string
  createdAt: string
  created_at?: string
}

export interface Order {
  id: string
  orderNumber: string
  order_number?: string
  userId?: string
  user_id?: string
  shippingAddress: ShippingAddress
  shipping_address?: ShippingAddress
  items: OrderItem[]
  subTotal: number
  discount: number
  total: number
  totalAmount?: number
  total_amount?: number
  status: 'PENDING' | 'PAID' | 'SHIPPING' | 'COMPLETED' | 'CANCELED' | 'pending' | 'confirmed' | 'processing' | 'shipping' | 'delivered' | 'cancelled'
  paymentStatus?: string
  payment_status?: string
  paymentMethod?: string
  payment_method?: string
  statusHistory?: StatusHistory[]
  status_history?: StatusHistory[]
  createdAt: string
  created_at?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  limit: number
  total: number
  totalPages: number
}

export interface ApiError {
  message: string
  code: string
  details?: any
}
