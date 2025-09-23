-- name: CreateCustomer :one
INSERT INTO customers (oidc_sub, name, email, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;



-- name: GetCustomerByOIDCSub :one
SELECT id, oidc_sub, name, email, phone, created_at
FROM customers
WHERE oidc_sub = $1;
