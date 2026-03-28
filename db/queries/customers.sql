-- name: CreateCustomer :one
INSERT INTO customers (
    cust_name, cust_email, cust_address, password_hash
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers WHERE cust_id = $1 LIMIT 1;

-- name: GetCustomerByEmail :one
-- Essential for the JWT Login process
SELECT * FROM customers WHERE cust_email = $1 LIMIT 1;

-- name: UpdateCustomerProfile :exec
UPDATE customers 
SET cust_name = $2, cust_address = $3 
WHERE cust_id = $1;

-- name: UpdateCustomerPassword :exec
UPDATE customers 
SET password_hash = $2 
WHERE cust_id = $1;

-- name: DeleteCustomer :exec
-- Because we set ON DELETE CASCADE, this automatically wipes their orders too!
DELETE FROM customers WHERE cust_id = $1;

-- name: ListCustomers :many
-- Used for the Admin Dashboard to view all registered users with pagination
SELECT * FROM customers 
ORDER BY cust_id DESC 
LIMIT $1 OFFSET $2;

-- name: SearchCustomers :many
-- Allows an admin to quickly find a user by typing part of their name or email
SELECT * FROM customers 
WHERE (sqlc.narg('search_term')::text IS NULL OR 
       cust_name ILIKE '%' || sqlc.narg('search_term') || '%' OR 
       cust_email ILIKE '%' || sqlc.narg('search_term') || '%')
ORDER BY cust_id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCustomerEmail :exec
UPDATE customers 
SET cust_email = $2 
WHERE cust_id = $1;