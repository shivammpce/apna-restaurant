-- -- name: CreateMenu :one
-- INSERT INTO menus (category, menu_item_ids)
-- VALUES (
--     $1,
--     (SELECT ARRAY[$2])
-- )
-- RETURNING *;

-- name: CreateMenu :one
-- params: category:string, menu_item_ids:array[string]
WITH inserted AS (
    INSERT INTO menus (category, menu_item_ids)
    VALUES (
        $1,
        (SELECT ARRAY[$2::UUID])
    )
    RETURNING id
)
SELECT id, $1 AS category, $2 AS menu_item_ids
FROM inserted;


-- name: UpdateMenu :one
-- param: category: string
-- param: menu_item_ids: []string
-- param: id: uuid
UPDATE menus
SET
    category = $1,
    menu_item_ids = $2
WHERE
    id = $3
RETURNING *;

-- name: GetAllMenus :many
SELECT * FROM menus
ORDER BY id;

-- name: GetMenuByID :one
-- param: id: uuid
SELECT id, category, menu_item_ids
FROM menus
WHERE id = $1;

-- name: CheckExistingMenu :one
SELECT COUNT(*) AS menu_count
FROM menus
WHERE id = $1;

-- name: DeleteMenuByID :exec
-- param: id: uuid
DELETE FROM menus
WHERE id = $1;


-- name: CreateMenuitem :one
INSERT INTO menuitems (
  name,
  price,
  image_url,
  menu_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetMenuitemsById :one
SELECT * FROM menuitems
WHERE id = $1 LIMIT 1;

-- name: ListMenuitems :many
SELECT * FROM menuitems
ORDER BY id;

-- name: UpdateMenuitem :one
UPDATE menuitems
SET name = $2,
    price = $3,
    image_url = $4
WHERE id = $1
RETURNING *;

-- name: CheckExistingMenuitem :one
SELECT COUNT(*) AS menuitem_count
FROM menuitems
WHERE id = $1;

-- name: DeleteMenuitem :exec
DELETE FROM menuitems
WHERE id = $1;