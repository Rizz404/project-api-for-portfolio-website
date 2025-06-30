-- name: CreateProjectTranslation :one
INSERT INTO project_translations (id, id_project, id_language, name, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateProjectTranslationsBatch :many
INSERT INTO project_translations (id, id_project, id_language, name, description)
SELECT (item->>'id')::uuid,
  (item->>'id_project')::uuid,
  (item->>'id_language')::uuid,
  item->>'name',
  item->>'description'
FROM json_array_elements($1::json) AS item
RETURNING *;

-- name: GetProjectTranslationsPaginated :many
SELECT *
FROM project_translations
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProjectTranslations :one
SELECT COUNT(*) as total
FROM project_translations;

-- name: GetProjectTranslationsCursorFirst :many
SELECT *
FROM project_translations
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchProjectTranslations :many
SELECT *
FROM project_translations
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC;

-- name: SearchProjectTranslationsPaginated :many
SELECT *
FROM project_translations
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProjectTranslation :one
SELECT *
FROM project_translations
WHERE id = $1
LIMIT 1;

-- name: GetProjectTranslationByName :one
SELECT *
FROM project_translations
WHERE name = $1
LIMIT 1;

-- name: CheckProjectTranslationExists :one
SELECT EXISTS (
    SELECT 1
    FROM project_translations
    WHERE id = $1
  ) as exists;

-- name: UpdateProjectTranslation :one
UPDATE project_translations
SET id_language = COALESCE($2, id_language),
  name = COALESCE($3, name),
  description = COALESCE($4, description)
WHERE id = $1
RETURNING *;

-- name: DeleteProjectTranslation :exec
DELETE FROM project_translations
WHERE id = $1;
