-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    email,
    password,
    role,
    address,
    full_name
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUsersPaginated :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) as total
FROM users;

-- name: GetUsersCursorFirst :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchUsers :many
SELECT *
FROM users
WHERE (
    username ILIKE '%' || $1 || '%'
    OR email ILIKE '%' || $1 || '%'
    OR full_name ILIKE '%' || $1 || '%'
  )
ORDER BY CASE
    WHEN username ILIKE $1 || '%' THEN 1
    WHEN email ILIKE $1 || '%' THEN 2
    WHEN full_name ILIKE $1 || '%' THEN 3
    ELSE 4
  END,
  created_at DESC;

-- name: SearchUsersPaginated :many
SELECT *
FROM users
WHERE (
    username ILIKE '%' || $1 || '%'
    OR email ILIKE '%' || $1 || '%'
    OR full_name ILIKE '%' || $1 || '%'
  )
ORDER BY CASE
    WHEN username ILIKE $1 || '%' THEN 1
    WHEN email ILIKE $1 || '%' THEN 2
    WHEN full_name ILIKE $1 || '%' THEN 3
    ELSE 4
  END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchUsersByUsername :many
SELECT *
FROM users
WHERE username ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN username ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  username ASC;

-- name: SearchUsersByEmail :many
SELECT *
FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN email ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  email ASC;

-- name: SearchUsersByFullName :many
SELECT *
FROM users
WHERE full_name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN full_name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  full_name ASC;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CheckUsernameExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE username = $1
  ) as exists;

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
  ) as exists;

-- name: CheckUserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE id = $1
  ) as exists;

-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username),
  email = COALESCE($3, email),
  password = COALESCE($4, password),
  role = COALESCE($5, role),
  address = COALESCE($6, address),
  full_name = COALESCE($7, full_name)
WHERE id = $1
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE users
SET address = COALESCE($2, address),
  full_name = COALESCE($3, full_name)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
