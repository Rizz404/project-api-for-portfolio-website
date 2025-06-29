-- name: CreateLanguage :one
INSERT INTO languages(id, name, lang_code)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetLanguagesPaginated :many
SELECT *
FROM languages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountLanguages :one
SELECT COUNT(*) as total
FROM languages;

-- name: GetLanguagesCursorFirst :many
SELECT *
FROM languages
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchLanguages :many
SELECT *
FROM languages
WHERE (
    name ILIKE '%' || $1 || '%'
    OR lang_code ILIKE '%' || $1 || '%'
  )
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    WHEN lang_code ILIKE $1 || '%' THEN 2
    ELSE 3
  END,
  created_at DESC;

-- name: SearchLanguagesPaginated :many
SELECT *
FROM languages
WHERE (
    name ILIKE '%' || $1 || '%'
    OR lang_code ILIKE '%' || $1 || '%'
  )
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    WHEN lang_code ILIKE $1 || '%' THEN 2
    ELSE 3
  END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchLanguagesByName :many
SELECT *
FROM languages
WHERE name ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN name ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  name ASC;

-- name: SearchLanguagesByLangCode :many
SELECT *
FROM languages
WHERE lang_code ILIKE '%' || $1 || '%'
ORDER BY CASE
    WHEN lang_code ILIKE $1 || '%' THEN 1
    ELSE 2
  END,
  lang_code ASC;

-- name: GetLanguage :one
SELECT *
FROM languages
WHERE id = $1
LIMIT 1;

-- name: GetLanguageByName :one
SELECT *
FROM languages
WHERE name = $1
LIMIT 1;

-- name: GetLanguageByLangCode :one
SELECT *
FROM languages
WHERE lang_code = $1
LIMIT 1;

-- name: CheckNameExists :one
SELECT EXISTS (
    SELECT 1
    FROM languages
    WHERE name = $1
  ) as exists;

-- name: CheckLangCodeExists :one
SELECT EXISTS (
    SELECT 1
    FROM languages
    WHERE lang_code = $1
  ) as exists;

-- name: CheckLanguageExists :one
SELECT EXISTS (
    SELECT 1
    FROM languages
    WHERE id = $1
  ) as exists;

-- name: UpdateLanguage :one
UPDATE languages
SET name = COALESCE($2, name),
  lang_code = COALESCE($3, lang_code)
WHERE id = $1
RETURNING *;

-- name: DeleteLanguage :exec
DELETE FROM languages
WHERE id = $1;
