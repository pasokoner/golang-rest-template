-- name: GetManyUser :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $1, password = $2, name = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;



