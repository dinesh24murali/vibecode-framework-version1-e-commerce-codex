-- name: GetInventoryByProduct :one
SELECT * FROM inventory_items
WHERE product_id = $1;

-- name: GetInventoryForUpdate :one
SELECT * FROM inventory_items
WHERE product_id = $1
FOR UPDATE;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items (product_id, on_hand, reserved, reorder_threshold)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateInventory :one
UPDATE inventory_items
SET on_hand           = $2,
    reorder_threshold = $3,
    updated_at        = now()
WHERE product_id = $1
RETURNING *;

-- name: ReserveStock :one
UPDATE inventory_items
SET reserved   = reserved + $2,
    updated_at = now()
WHERE product_id = $1
RETURNING *;

-- name: ReleaseStock :one
UPDATE inventory_items
SET reserved   = reserved - $2,
    updated_at = now()
WHERE product_id = $1
RETURNING *;

-- name: DeductStock :one
UPDATE inventory_items
SET on_hand    = on_hand - $2,
    reserved   = reserved - $2,
    updated_at = now()
WHERE product_id = $1
RETURNING *;
