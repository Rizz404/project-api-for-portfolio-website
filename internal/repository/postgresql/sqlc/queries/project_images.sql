-- name: CreateProjectImage :one
INSERT INTO project_images (id, id_project, file_name, url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateProjectImagesBatch :many
INSERT INTO project_images (id, id_project, file_name, url)
SELECT (item->>'id')::uuid,
  (item->>'id_project')::uuid,
  item->>'file_name',
  item->>'url'
FROM json_array_elements($1::json) AS item
RETURNING *;

-- name: GetProjectImagesPaginated :many
SELECT *
FROM project_images
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProjectImages :one
SELECT COUNT(*) as total
FROM project_images;

-- name: GetProjectImagesCursorFirst :many
SELECT *
FROM project_images
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchProjectImages :many
SELECT *
FROM project_images
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC;

-- name: SearchProjectImagesPaginated :many
SELECT *
FROM project_images
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProjectImage :one
SELECT *
FROM project_images
WHERE id = $1
LIMIT 1;

-- name: GetProjectImageByName :one
SELECT *
FROM project_images
WHERE file_name = $1
LIMIT 1;

-- name: CheckProjectImageExists :one
SELECT EXISTS (
    SELECT 1
    FROM project_images
    WHERE id = $1
  ) as exists;

-- name: UpdateProjectImage :one
UPDATE project_images
SET file_name = COALESCE($2, file_name),
  url = COALESCE($3, url)
WHERE id = $1
RETURNING *;

-- name: DeleteProjectImage :exec
DELETE FROM project_images
WHERE id = $1;
