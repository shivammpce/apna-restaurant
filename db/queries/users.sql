-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  phone_number,
  password,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    name = $3,
    phone_number = $4,
    password = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;

-- name: CheckExistingUser :one
SELECT COUNT(*) AS email_count
FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;