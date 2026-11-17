import { addCartItem, clearCart, getCart, removeCartItem, updateCartItem } from './generated/cart/cart'
import type { Cart } from './generated/model'
import type { CartApi, CartLine } from '@/lib/store'

const toLines = (cart: Cart): CartLine[] =>
  cart.items.map((i) => ({
    productId: i.product.id,
    quantity: i.quantity,
    itemId: i.id,
    unitPrice: i.unit_price,
  }))

const unwrap = (res: { status: number; data: unknown }): CartLine[] => {
  if (res.status >= 400) throw new Error(`cart request failed: ${res.status}`)
  return toLines(res.data as Cart)
}

/** Binds the TASK-008 cart slice to the generated cart endpoints. */
export const cartApi: CartApi = {
  fetchCart: async () => unwrap(await getCart()),
  addItem: async (productId, quantity) =>
    unwrap(await addCartItem({ product_id: productId, quantity })),
  updateItem: async (itemId, quantity) => unwrap(await updateCartItem(itemId, { quantity })),
  removeItem: async (itemId) => unwrap(await removeCartItem(itemId)),
  clear: async () => {
    await clearCart()
  },
}
