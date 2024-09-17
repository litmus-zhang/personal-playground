import { computed, map } from "nanostores";


export const $cart = map<Record<number, CartItem>>({});


export function addItemToCart(item: ShopItem) {
    const CartItem = $cart.get()[item.id];
    const Cartquantity = CartItem ? CartItem.quantity : 0;

    $cart.setKey(item.id, {
        item,
        quantity: Cartquantity + 1,
    });
}


export function removeItemFromCart(itemId: number) {
    $cart.setKey(itemId, undefined);
}

export const subTotal = computed($cart, (cart) => {
    let subtotal = 0
    Object.values(cart).forEach((cartItem) => {
        if (!cartItem) return;
        subtotal += cartItem.item.price * cartItem.quantity;
    });
    return subtotal;
});

