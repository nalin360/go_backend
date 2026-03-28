-- name: CreateOrder :one
INSERT INTO orders (
    cust_id, status
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: ListCustomerOrders :many
-- Get a user's order history, newest first, with pagination
SELECT * FROM orders 
WHERE cust_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :exec
-- Used by admins to mark orders as 'shipped' or 'cancelled'
UPDATE orders 
SET status = $2 
WHERE id = $1;