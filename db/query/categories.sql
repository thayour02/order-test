-- db/queries.sql

-- name: CreateCategory :one
INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING id, name, parent_id;



-- name: AvgPriceForCategory :one
WITH RECURSIVE cat_tree AS (
  SELECT categories.id AS cat_id
  FROM categories
  WHERE categories.id = $1
  UNION ALL
  SELECT c.id
  FROM categories c
  JOIN cat_tree ct ON c.parent_id = ct.cat_id
)
SELECT CAST(COALESCE(AVG(p.price), 0) AS float8) AS avg_price
FROM products p
JOIN product_categories pc ON pc.product_id = p.id
WHERE pc.category_id IN (SELECT cat_id FROM cat_tree);

