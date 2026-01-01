import api from '@/lib/api'
import { Product, ProductDetail, Brand, Category, PaginatedResponse } from '@/types'

export interface ProductQuery {
  search?: string
  brand?: string
  category?: string
  minPrice?: number
  maxPrice?: number
  page?: number
  limit?: number
  sort?: string
}

export const productService = {
  getProducts: async (query: ProductQuery = {}): Promise<PaginatedResponse<Product>> => {
    const response = await api.get<PaginatedResponse<Product>>('/products', {
      params: query,
    })
    return response.data
  },

  getProductBySlug: async (slug: string): Promise<ProductDetail> => {
    const response = await api.get<ProductDetail>(`/products/${slug}`)
    return response.data
  },

  getBrands: async (): Promise<Brand[]> => {
    const response = await api.get<Brand[]>('/brands')
    return response.data
  },

  getCategories: async (): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories')
    return response.data
  },

  // Admin endpoints
  createProduct: async (data: any): Promise<void> => {
    const response = await api.post('/admin/products', data)
    return response.data
  },

  updateProduct: async (id: string, data: any): Promise<void> => {
    const response = await api.put(`/admin/products/${id}`, data)
    return response.data
  },

  deleteProduct: async (id: string): Promise<void> => {
    const response = await api.delete(`/admin/products/${id}`)
    return response.data
  },
}
