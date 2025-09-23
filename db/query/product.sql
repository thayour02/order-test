-- name: CreateProduct :one
INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id,  name, description, price, created_at;

-- name: AddProductCategory :exec
INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;



-- name: GetProductByID :one
SELECT id, name, description, price, created_at FROM products WHERE id = $1;



-- name: ProductsInCategoryRecursive :many
WITH RECURSIVE cat_tree AS (
  SELECT categories.id AS cat_id
  FROM categories
  WHERE categories.id = $1
  UNION ALL
  SELECT c.id
  FROM categories c
  JOIN cat_tree ct ON c.parent_id = ct.cat_id
)
SELECT p.id,  p.name, p.description, p.price, p.created_at
FROM products p
JOIN product_categories pc ON pc.product_id = p.id
WHERE pc.category_id IN (SELECT cat_id FROM cat_tree);


-- name: GetAllProducts :many
SELECT id, name, description, price, created_at
FROM products
ORDER BY created_at DESC;
