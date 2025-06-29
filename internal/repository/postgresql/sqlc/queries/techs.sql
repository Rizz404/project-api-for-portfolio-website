-- name: CreateTech :one
INSERT INTO techs (id, name, description, logo_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTechsPaginated :many
SELECT *
FROM techs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTechs :one
SELECT COUNT(*) as total
FROM techs;

-- name: GetTechsCursorFirst :many
SELECT *
FROM techs
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchTechs :many
SELECT *
FROM techs
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC;

-- name: SearchTechsPaginated :many
SELECT *
FROM techs
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTech :one
SELECT *
FROM techs
WHERE id = $1
LIMIT 1;

-- name: GetTechByName :one
SELECT *
FROM techs
WHERE name = $1
LIMIT 1;

-- name: CheckTechNameExists :one
SELECT EXISTS (
    SELECT 1
    FROM techs
    WHERE name = $1
  ) as exists;

-- name: CheckTechExists :one
SELECT EXISTS (
    SELECT 1
    FROM techs
    WHERE id = $1
  ) as exists;

-- name: UpdateTech :one
UPDATE techs
SET name = COALESCE($2, name),
  description = COALESCE($3, description),
  logo_url = COALESCE($4, logo_url)
WHERE id = $1
RETURNING *;

-- name: DeleteTech :exec
DELETE FROM techs
WHERE id = $1;
