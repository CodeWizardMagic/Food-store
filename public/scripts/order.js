const state = {
    products: [],
    cart: [],
    currentPage: 1,
    perPage: 6,
    total: 0,
    filters: {
        name: '',
        minPrice: '',
        maxPrice: ''
    },
    userId: '',
    orders: []
};

const DOM = {
    productsList: document.getElementById('products-list'),
    pagination: document.getElementById('pagination'),
    cartItems: document.getElementById('cart-items'),
    cartEmpty: document.getElementById('cart-empty'),
    cartCount: document.getElementById('cart-count'),
    cartTotal: document.getElementById('cart-total'),
    cartSection: document.getElementById('cart-section'),
    cartIcon: document.getElementById('cart-icon'),
    checkoutButton: document.getElementById('checkout-button'),
    orderForm: document.getElementById('order-form'),
    userIdInput: document.getElementById('user-id'),
    nameFilter: document.getElementById('name-filter'),
    minPriceFilter: document.getElementById('min-price'),
    maxPriceFilter: document.getElementById('max-price'),
    applyFiltersButton: document.getElementById('apply-filters'),
    resetFiltersButton: document.getElementById('reset-filters'),
    userIdFilter: document.getElementById('user-id-filter'),
    viewOrdersButton: document.getElementById('view-orders'),
    ordersList: document.getElementById('orders-list')
};

document.addEventListener('DOMContentLoaded', () => {
    fetchProducts();

    DOM.applyFiltersButton.addEventListener('click', applyFilters);
    DOM.resetFiltersButton.addEventListener('click', resetFilters);
    DOM.orderForm.addEventListener('submit', placeOrder);
    DOM.viewOrdersButton.addEventListener('click', fetchOrders);
    DOM.userIdInput.addEventListener('input', updateCheckoutButton);

    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
        state.cart = JSON.parse(savedCart);
        updateCart();
    }

    const savedUserId = localStorage.getItem('userId');
    if (savedUserId) {
        DOM.userIdInput.value = savedUserId;
        DOM.userIdFilter.value = savedUserId;
        state.userId = savedUserId;
        updateCheckoutButton();
        fetchOrders();
    }
});

async function fetchProducts() {
    try {
        const url = new URL('/api/products', window.location.origin);
        url.searchParams.append('page', state.currentPage);
        url.searchParams.append('per_page', state.perPage);

        if (state.filters.name) {
            url.searchParams.append('name', state.filters.name);
        }
        if (state.filters.minPrice) {
            url.searchParams.append('min_price', state.filters.minPrice);
        }
        if (state.filters.maxPrice) {
            url.searchParams.append('max_price', state.filters.maxPrice);
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Error fetching products: ${response.statusText}`);
        }

        const data = await response.json();
        state.products = data.products;
        state.total = data.total;
        
        renderProducts();
        renderPagination();
    } catch (error) {
        console.error('Error fetching products:', error);
    }
}

function renderProducts() {
    DOM.productsList.innerHTML = '';

    if (state.products.length === 0) {
        DOM.productsList.innerHTML = '<div class="main__products-empty">Products not found</div>';
        return;
    }

    state.products.forEach(product => {
        const productElement = document.createElement('div');
        productElement.className = 'main__product-item';
        
        const isInCart = state.cart.some(item => item.id === product.ID);
        const isOutOfStock = product.Stock <= 0;
        
        productElement.innerHTML = `
            <div class="main__product-item-name">${product.Name}</div>
            <div class="main__product-item-price">${product.Price.toFixed(2)} ₸</div>
            <div class="main__product-item-stock">Stock: ${product.Stock}</div>
            <button class="main__product-item-button" 
                    data-id="${product.ID}" 
                    data-name="${product.Name}" 
                    data-price="${product.Price}" 
                    data-stock="${product.Stock}"
                    ${isInCart || isOutOfStock ? 'disabled' : ''}
            >
                ${isInCart ? 'In cart' : isOutOfStock ? 'Out of stock' : 'Add to cart'}
            </button>
        `;
        
        const addButton = productElement.querySelector('button');
        if (!isInCart && !isOutOfStock) {
            addButton.addEventListener('click', addToCart);
        }
        
        DOM.productsList.appendChild(productElement);
    });
}

function renderPagination() {
    DOM.pagination.innerHTML = '';
    
    const totalPages = Math.ceil(state.total / state.perPage);
    if (totalPages <= 1) return;
    
    if (state.currentPage > 1) {
        const prevButton = document.createElement('button');
        prevButton.textContent = 'prev';
        prevButton.addEventListener('click', () => {
            state.currentPage--;
            fetchProducts();
        });
        DOM.pagination.appendChild(prevButton);
    }
    
    for (let i = 1; i <= totalPages; i++) {
        const pageButton = document.createElement('button');
        pageButton.textContent = i;
        pageButton.disabled = i === state.currentPage;
        pageButton.addEventListener('click', () => {
            state.currentPage = i;
            fetchProducts();
        });
        DOM.pagination.appendChild(pageButton);
    }
    
    if (state.currentPage < totalPages) {
        const nextButton = document.createElement('button');
        nextButton.textContent = 'next';
        nextButton.addEventListener('click', () => {
            state.currentPage++;
            fetchProducts();
        });
        DOM.pagination.appendChild(nextButton);
    }
}

function addToCart(event) {
    const button = event.target;
    const id = button.dataset.id;
    const name = button.dataset.name;
    const price = parseFloat(button.dataset.price);
    const stock = parseInt(button.dataset.stock);
    
    const existingItem = state.cart.find(item => item.id === id);
    
    if (existingItem) {
        existingItem.quantity += 1;
    } else {
        state.cart.push({
            id,
            name,
            price,
            stock,
            quantity: 1
        });
    }
    
    button.disabled = true;
    button.textContent = 'Cart';
    
    localStorage.setItem('cart', JSON.stringify(state.cart));
    
    updateCart();
}

function updateCart() {
    const totalItems = state.cart.reduce((sum, item) => sum + item.quantity, 0);
    DOM.cartCount.textContent = totalItems;
    
    if (totalItems === 0) {
        DOM.cartEmpty.style.display = 'block';
        DOM.cartItems.innerHTML = '';
        DOM.cartTotal.textContent = '0 ₸';
        DOM.checkoutButton.disabled = true;
        return;
    }
    
    DOM.cartEmpty.style.display = 'none';
    
    DOM.cartItems.innerHTML = '';
    
    let totalPrice = 0;
    
    state.cart.forEach(item => {
        const itemTotal = item.price * item.quantity;
        totalPrice += itemTotal;
        
        const cartItem = document.createElement('div');
        cartItem.className = 'main__cart-item';
        cartItem.innerHTML = `
            <div class="main__cart-item-details">
                <div class="main__cart-item-name">${item.name}</div>
                <div class="main__cart-item-price">${item.price.toFixed(2)} ₸ x ${item.quantity}</div>
            </div>
            <div class="main__cart-item-quantity">
                <button class="decrease" data-id="${item.id}" ${item.quantity <= 1 ? 'disabled' : ''}>-</button>
                <span>${item.quantity}</span>
                <button class="increase" data-id="${item.id}" ${item.quantity >= item.stock ? 'disabled' : ''}>+</button>
            </div>
            <button class="main__cart-item-remove" data-id="${item.id}">×</button>
        `;
        
        const decreaseButton = cartItem.querySelector('.decrease');
        const increaseButton = cartItem.querySelector('.increase');
        const removeButton = cartItem.querySelector('.main__cart-item-remove');
        
        decreaseButton.addEventListener('click', () => decreaseQuantity(item.id));
        increaseButton.addEventListener('click', () => increaseQuantity(item.id));
        removeButton.addEventListener('click', () => removeFromCart(item.id));
        
        DOM.cartItems.appendChild(cartItem);
    });
    
    DOM.cartTotal.textContent = `${totalPrice.toFixed(2)} ₸`;
    
    updateCheckoutButton();
}

function decreaseQuantity(id) {
    const item = state.cart.find(item => item.id === id);
    if (!item) return;
    
    item.quantity -= 1;
    
    if (item.quantity <= 0) {
        removeFromCart(id);
    } else {
        localStorage.setItem('cart', JSON.stringify(state.cart));
        updateCart();
    }
}

function increaseQuantity(id) {
    const item = state.cart.find(item => item.id === id);
    if (!item || item.quantity >= item.stock) return;
    
    item.quantity += 1;
    localStorage.setItem('cart', JSON.stringify(state.cart));
    updateCart();
}

function removeFromCart(id) {
    state.cart = state.cart.filter(item => item.id !== id);
    localStorage.setItem('cart', JSON.stringify(state.cart));
    
    const productButton = document.querySelector(`button[data-id="${id}"]`);
    if (productButton) {
        productButton.disabled = false;
        productButton.textContent = 'Add to cart';
    }
    
    updateCart();
}

function applyFilters() {
    state.filters.name = DOM.nameFilter.value;
    state.filters.minPrice = DOM.minPriceFilter.value;
    state.filters.maxPrice = DOM.maxPriceFilter.value;
    state.currentPage = 1;
    fetchProducts();
}

function resetFilters() {
    DOM.nameFilter.value = '';
    DOM.minPriceFilter.value = '';
    DOM.maxPriceFilter.value = '';
    state.filters = { name: '', minPrice: '', maxPrice: '' };
    state.currentPage = 1;
    fetchProducts();
}

function updateCheckoutButton() {
    const userId = DOM.userIdInput.value.trim();
    state.userId = userId;
    
    if (userId && state.cart.length > 0) {
        DOM.checkoutButton.disabled = false;
    } else {
        DOM.checkoutButton.disabled = true;
    }
    
    if (userId) {
        localStorage.setItem('userId', userId);
    }
}

async function placeOrder(event) {
    event.preventDefault();
    
    const userId = DOM.userIdInput.value.trim();
    if (!userId || state.cart.length === 0) return;
    
    try {
        const orderData = {
            user_id: userId,
            items: state.cart.map(item => ({
                product_id: item.id,
                quantity: item.quantity
            }))
        };
        
        const response = await fetch('/api/orders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(orderData)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Error creating order: ${errorText}`);
        }
        
        const result = await response.json();
        
        state.cart = [];
        localStorage.setItem('cart', JSON.stringify(state.cart));
        updateCart();
        
        fetchProducts();
        
        fetchOrders();
        
        alert(`Order created! Order number: ${result.order_id}`);
        
    } catch (error) {
        console.error('Error creating order:', error);
        alert(error.message);
    }
}

async function fetchOrders() {
    const userId = DOM.userIdFilter.value.trim();
    if (!userId) return;
    
    try {
        const url = new URL('/api/orders', window.location.origin);
        url.searchParams.append('user_id', userId);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Error fetching orders: ${response.statusText}`);
        }
        
        const orders = await response.json();
        state.orders = orders;
        
        renderOrders();
    } catch (error) {
        console.error('Error fetching orders:', error);
        DOM.ordersList.innerHTML = '<div class="main__orders-error">Error fetching orders</div>';
    }
}

function renderOrders() {
    DOM.ordersList.innerHTML = '';
    
    if (state.orders.length === 0) {
        DOM.ordersList.innerHTML = '<div class="main__orders-empty">You do not have any orders</div>';
        return;
    }
    
    state.orders.forEach(order => {
        const orderElement = document.createElement('div');
        orderElement.className = 'main__order-item';
        
        const createdDate = new Date(order.created_at).toLocaleString();
        
        orderElement.innerHTML = `
            <div class="main__order-item-header">
                <div class="main__order-item-date">${createdDate}</div>
                <div class="main__order-item-id">ID: ${order.id}</div>
                <div class="main__order-item-status ${order.status}">${getStatusText(order.status)}</div>
            </div>
            ${order.items ? renderOrderItems(order.items) : '<div class="main__order-item-products-empty">Details are not available</div>'}
            <div class="main__order-item-total">
                <span>Amount:</span>
                <span>${order.total_price.toFixed(2)} ₸</span>
            </div>
            ${renderOrderActions(order)}
        `;
        
        DOM.ordersList.appendChild(orderElement);
        
        const actionButtons = orderElement.querySelectorAll('.main__order-item-button');
        actionButtons.forEach(button => {
            button.addEventListener('click', () => updateOrderStatus(order.id, button.dataset.status));
        });
    });
}

function renderOrderItems(items) {
    if (!items || items.length === 0) return '';
    
    let html = '<div class="main__order-item-products">';
    
    items.forEach(item => {
        const productName = item.product ? item.product.Name : 'Product not found';
        html += `
            <div class="main__order-item-product">
                <div class="main__order-item-product-name">${productName}</div>
                <div class="main__order-item-product-quantity">Quantity: ${item.quantity}</div>
                <div class="main__order-item-product-price">${item.price.toFixed(2)} ₸</div>
            </div>
        `;
    });
    
    html += '</div>';
    return html;
}

function renderOrderActions(order) {
    if (order.status === 'completed' || order.status === 'cancelled') {
        return '';
    }
    
    return `
        <div class="main__order-item-actions">
            ${order.status === 'pending' ? `
                <button class="main__order-item-button complete" data-status="completed">Complete</button>
                <button class="main__order-item-button cancel" data-status="cancelled">Cancel</button>
            ` : ''}
        </div>
    `;
}

async function updateOrderStatus(orderId, status) {
    try {
        const response = await fetch(`/api/orders/${orderId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status })
        });
        
        if (!response.ok) {
            throw new Error(`Error updating order: ${response.statusText}`);
        }
        
        fetchOrders();
        
        if (status === 'cancelled') {
            fetchProducts();
        }
    } catch (error) {
        console.error('Error updating order:', error);
        alert(error.message);
    }
}

function getStatusText(status) {
    switch (status) {
        case 'pending':
            return 'Pending';
        case 'completed':
            return 'Completed';
        case 'cancelled':
            return 'Cancelled';
        default:
            return status;
    }
}