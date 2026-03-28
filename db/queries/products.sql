-- name: CreateProduct :one
INSERT INTO products (
    prodc_name, prodc_price, stock_on_hand, expiry_date
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: CreateProductsBatch :copyfrom
-- (This tells sqlc to generate a function that accepts a slice of products)
INSERT INTO products (
    prodc_name, prodc_price, stock_on_hand, expiry_date
) VALUES (
    $1, $2, $3, $4
);

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
-- Includes Pagination
SELECT * FROM products 
ORDER BY id DESC 
LIMIT $1 OFFSET $2;

-- name: SearchProducts :many
-- Dynamic search using sqlc.narg() + Pagination
SELECT * FROM products 
WHERE (sqlc.narg('name_query')::text IS NULL OR prodc_name ILIKE '%' || sqlc.narg('name_query') || '%')
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: GetProductInventory :one
-- Queries the View we created to get real-time available stock
SELECT * FROM product_inventory WHERE product_id = $1 LIMIT 1;

-- name: AddProductStock :exec
-- Atomic update: Safe from race conditions when two admins add stock at the same time
UPDATE products 
SET stock_on_hand = stock_on_hand + $2 
WHERE id = $1;

-- name: UpdateProductPrice :exec
UPDATE products 
SET prodc_price = $2 
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: UpdateProductStock :exec
UPDATE products 
SET stock_on_hand = $2 
WHERE id = $1;

-- name: ListProductsByIDs :many
SELECT * FROM products WHERE id = ANY($1);

-- name: UpdateProductName :exec
UPDATE products 
SET prodc_name = $2 
WHERE id = $1;