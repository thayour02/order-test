-- name: CreateOrder :one
INSERT INTO orders (customer_id, total)
VALUES ($1, $2)
RETURNING id, customer_id, total, created_at;


-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, subtotal)
VALUES ($1, $2, $3, $4)
RETURNING  order_id, product_id, quantity, subtotal;


-- name: UpdateOrderTotal :exec
UPDATE orders
SET total = $2
WHERE id = $1;


-- name: GetOrderByID :one
SELECT id, customer_id, total, created_at
FROM orders
WHERE id = $1;

-- name: GetOrdersByCustomerID :many
SELECT id, customer_id, total, created_at
FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC;

