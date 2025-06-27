-- name: CreateUserTranslation :one
INSERT INTO
  user_translations (
    id,
    id_user,
    id_language,
    bio,
    about_me,
    additional_skills,
    languages,
    quote
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  *;

-- name: GetUserTranslationsPaginated :many
SELECT
  *
FROM
  user_translations
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: GetUserTranslationsByUserIDPaginated :many
SELECT
  *
FROM
  user_translations
WHERE
  id_user = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: CountUserTranslations :one
SELECT
  COUNT(*) as total
FROM
  user_translations;

-- name: GetUserTranslationsCursorFirst :many
SELECT
  *
FROM
  user_translations
ORDER BY
  created_at DESC
LIMIT
  $1;

-- name: GetUserTranslation :one
SELECT
  *
FROM
  user_translations
WHERE
  id = $1
LIMIT
  1;

-- name: GetUserTranslationByUserIDAndLangID :one
SELECT
  *
FROM
  user_translations
WHERE
  id_user = $1
  AND id_language = $2
LIMIT
  1;

-- name: GetUserTranslationByUserIDAndLangName :one
SELECT
  ut.*
FROM
  user_translations ut
  JOIN languages l ON ut.id_language = l.id
WHERE
  ut.id_user = $1
  AND l.name = $2
LIMIT
  1;

-- name: GetUserTranslationByUserIDAndLangCode :one
SELECT
  ut.*
FROM
  user_translations ut
  JOIN languages l ON ut.id_language = l.id
WHERE
  ut.id_user = $1
  AND l.lang_code = $2
LIMIT
  1;

-- name: CheckUserTranslationExists :one
SELECT
  EXISTS (
    SELECT
      1
    FROM
      user_translations
    WHERE
      id = $1
  ) as exists;

-- name: UpdateUserTranslation :one
UPDATE user_translations
SET
  id_user = COALESCE($2, id_user),
  id_language = COALESCE($3, id_language),
  bio = COALESCE($4, bio),
  about_me = COALESCE($5, about_me),
  additional_skills = COALESCE($6, additional_skills),
  languages = COALESCE($7, languages),
  quote = COALESCE($8, quote)
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteUserTranslation :exec
DELETE FROM user_translations
WHERE
  id = $1;
