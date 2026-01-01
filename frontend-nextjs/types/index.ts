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
  isActive: boolean
  isFeatured: boolean
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
}

export interface Order {
  id: string
  orderNumber: string
  shippingAddress: ShippingAddress
  items: OrderItem[]
  subTotal: number
  discount: number
  total: number
  status: 'PENDING' | 'PAID' | 'SHIPPING' | 'COMPLETED' | 'CANCELED'
  createdAt: string
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
