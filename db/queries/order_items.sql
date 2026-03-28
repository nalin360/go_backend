-- name: CreateOrderItem :exec
INSERT INTO order_items (
    order_id, prodc_id, quantity, total_price
) VALUES (
    $1, $2, $3, $4
);

-- name: ListOrderItemsWithProductDetails :many
-- Joins the items with the products table so the frontend gets the name and historical price
SELECT 
    oi.order_id, 
    oi.prodc_id, 
    oi.quantity, 
    oi.total_price AS purchased_price_cents,
    p.prodc_name
FROM order_items oi
JOIN products p ON oi.prodc_id = p.id
WHERE oi.order_id = $1;