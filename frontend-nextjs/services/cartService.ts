import api from '@/lib/api'
import { Cart } from '@/types'

export interface AddItemData {
  variantId: string
  quantity: number
}

export interface UpdateItemData {
  quantity: number
}

export const cartService = {
  getCart: async (): Promise<Cart> => {
    const response = await api.get<Cart>('/cart')
    return response.data
  },

  addItem: async (data: AddItemData): Promise<void> => {
    const response = await api.post('/cart/items', data)
    return response.data
  },

  updateItem: async (variantId: string, data: UpdateItemData): Promise<void> => {
    const response = await api.put(`/cart/items/${variantId}`, data)
    return response.data
  },

  removeItem: async (variantId: string): Promise<void> => {
    const response = await api.delete(`/cart/items/${variantId}`)
    return response.data
  },
}
