-- name: CreateTable :one
WITH inserted AS (
    INSERT INTO tables (table_number, order_ids)
    VALUES (
        $1,
        (SELECT ARRAY[$2::UUID])
    )
    RETURNING id
)
SELECT id, $1 AS table_number, $2 AS order_ids
FROM inserted;

-- name: UpdateTable :one
UPDATE tables
SET
    table_number = $1,
    order_ids = $2
WHERE
    id = $3
RETURNING *;

-- name: GetTableByID :one
SELECT * FROM tables WHERE id = $1;


-- name: GetAllTables :many
SELECT * FROM tables ORDER BY id;

-- name: CheckExistingTable :one
SELECT COUNT(*) AS table_count
FROM tables
WHERE id = $1;

-- name: DeleteTableByID :exec
DELETE FROM tables
WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO orders (
    table_id, amount, order_items
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET
    table_id = $1,
    amount = $2,
    order_items = $3,
    created_at = $4, 
    delivered_at = $5
WHERE
    id = $6
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;


-- name: GetAllOrders :many
SELECT * FROM orders
ORDER BY created_at;

-- name: CheckExistingOrder :one
SELECT COUNT(*) AS order_count
FROM orders
WHERE id = $1;

-- name: DeleteOrderByID :exec
DELETE FROM orders
WHERE id = $1;

-- name: RemoveOrderIDFromTables :exec
UPDATE tables
SET order_ids = array_remove(order_ids, $1)
WHERE $1 = ANY(order_ids);
