import type { StateCreator } from 'zustand'

export interface CartLine {
  productId: string
  quantity: number
  /** Server cart item id, present after a sync. */
  itemId?: string
  unitPrice?: number
}

/** Adapter over the generated API client (TASK-007 wires the real one). */
export interface CartApi {
  fetchCart: () => Promise<CartLine[]>
  addItem: (productId: string, quantity: number) => Promise<CartLine[]>
  updateItem: (itemId: string, quantity: number) => Promise<CartLine[]>
  removeItem: (itemId: string) => Promise<CartLine[]>
  clear: () => Promise<void>
}

export interface CartSlice {
  items: CartLine[]
  cartApi: CartApi | null
  setCartApi: (api: CartApi | null) => void
  addItem: (productId: string, quantity?: number) => void
  removeItem: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clearCart: () => void
  /** Pushes local lines to the server, then replaces local state with the server cart. */
  syncWithServer: () => Promise<void>
}

export const createCartSlice: StateCreator<CartSlice, [], [], CartSlice> = (set, get) => ({
  items: [],
  cartApi: null,
  setCartApi: (api) => set({ cartApi: api }),
  addItem: (productId, quantity = 1) =>
    set((s) => {
      const existing = s.items.find((i) => i.productId === productId)
      return {
        items: existing
          ? s.items.map((i) =>
              i.productId === productId ? { ...i, quantity: i.quantity + quantity } : i
            )
          : [...s.items, { productId, quantity }],
      }
    }),
  removeItem: (productId) =>
    set((s) => ({ items: s.items.filter((i) => i.productId !== productId) })),
  updateQuantity: (productId, quantity) =>
    set((s) => ({
      items:
        quantity <= 0
          ? s.items.filter((i) => i.productId !== productId)
          : s.items.map((i) => (i.productId === productId ? { ...i, quantity } : i)),
    })),
  clearCart: () => set({ items: [] }),
  syncWithServer: async () => {
    const api = get().cartApi
    if (!api) return
    const server = await api.fetchCart()
    for (const line of get().items) {
      if (!server.some((s) => s.productId === line.productId)) {
        await api.addItem(line.productId, line.quantity)
      }
    }
    set({ items: await api.fetchCart() })
  },
})
