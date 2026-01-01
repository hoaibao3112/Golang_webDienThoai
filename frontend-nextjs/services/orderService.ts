import api from '@/lib/api'
import { Order, ShippingAddress } from '@/types'

export interface CreateOrderData {
  shippingAddress: ShippingAddress
  voucherCode?: string
}

export const orderService = {
  createOrder: async (data: CreateOrderData): Promise<Order> => {
    const response = await api.post<Order>('/orders', data)
    return response.data
  },

  getMyOrders: async (): Promise<Order[]> => {
    const response = await api.get<Order[]>('/orders/me')
    return response.data
  },

  getOrderById: async (id: string): Promise<Order> => {
    const response = await api.get<Order>(`/orders/${id}`)
    return response.data
  },
}
