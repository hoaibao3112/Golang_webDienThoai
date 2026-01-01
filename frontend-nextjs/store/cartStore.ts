import { create } from 'zustand'
import { CartItem } from '@/types'

interface CartState {
  items: CartItem[]
  addItem: (item: CartItem) => void
  updateQuantity: (variantId: string, quantity: number) => void
  removeItem: (variantId: string) => void
  clearCart: () => void
  getTotalItems: () => number
  getTotalPrice: () => number
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  addItem: (item) => {
    const items = get().items
    const existingItem = items.find((i) => i.variantId === item.variantId)
    
    if (existingItem) {
      set({
        items: items.map((i) =>
          i.variantId === item.variantId
            ? { ...i, quantity: i.quantity + item.quantity }
            : i
        ),
      })
    } else {
      set({ items: [...items, item] })
    }
  },
  updateQuantity: (variantId, quantity) => {
    set({
      items: get().items.map((item) =>
        item.variantId === variantId ? { ...item, quantity } : item
      ),
    })
  },
  removeItem: (variantId) => {
    set({
      items: get().items.filter((item) => item.variantId !== variantId),
    })
  },
  clearCart: () => {
    set({ items: [] })
  },
  getTotalItems: () => {
    return get().items.reduce((total, item) => total + item.quantity, 0)
  },
  getTotalPrice: () => {
    return get().items.reduce(
      (total, item) => total + item.price * item.quantity,
      0
    )
  },
}))
