-- Customers Table
CREATE TABLE customers (
    cust_id BIGSERIAL PRIMARY KEY,
    cust_name VARCHAR(255) NOT NULL,
    cust_email VARCHAR(255) UNIQUE NOT NULL,
    cust_address TEXT NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE -- 🛡️ The Admin Flag
);


-- Products Table
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    prodc_name VARCHAR(255) NOT NULL,
    prodc_price INTEGER NOT NULL, -- Stored in cents
    stock_on_hand INTEGER NOT NULL DEFAULT 0,
    expiry_date DATE -- Can be NULL if product doesn't expire
);

-- Orders Table
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    cust_id BIGINT NOT NULL REFERENCES customers(cust_id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'shipped', 'cancelled'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Order Items Table (The Join Table)
CREATE TABLE order_items (
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    prodc_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    total_price INTEGER NOT NULL, -- Price at the time of purchase (in cents)
    PRIMARY KEY (order_id, prodc_id) -- Prevents duplicate product entries in the same order
);

-- View to get product inventory
CREATE VIEW product_inventory AS
SELECT 
    p.id AS product_id,
    p.prodc_name,
    p.stock_on_hand,
    COALESCE(SUM(oi.quantity), 0) AS total_ordered,
    (p.stock_on_hand - COALESCE(SUM(oi.quantity), 0)) AS available_inventory
FROM 
    products p
LEFT JOIN 
    order_items oi ON p.id = oi.prodc_id
GROUP BY 
    p.id;
