-- name: CreateProject :one
INSERT INTO projects (
    id,
    id_user,
    id_category,
    is_deployed,
    is_maintained,
    live_demo,
    source_code
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CountProjects :one
SELECT COUNT(*)
FROM projects;

-- name: CheckProjectExists :one
SELECT EXISTS (
    SELECT 1
    FROM projects
    WHERE id = $1
  );

-- name: GetProject :one
SELECT sqlc.embed(p),
  json_build_object(
    'id',
    u.id,
    'username',
    u.username,
    'email',
    u.email
  ) as user,
  json_build_object(
    'id',
    c.id,
    'name',
    c.name,
    'description',
    c.description
  ) as category,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pt.id,
            'id_language',
            pt.id_language,
            'name',
            pt.name,
            'description',
            pt.description,
            'language',
            json_build_object('id', l.id, 'code', l.code, 'name', l.name)
          )
          ORDER BY l.name
        ),
        '[]'::json
      )
    FROM project_translations pt
      JOIN languages l ON pt.id_language = l.id
    WHERE pt.id_project = p.id
  ) as translations,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            t.id,
            'name',
            t.name,
            'description',
            t.description,
            'logo_url',
            t.logo_url
          )
          ORDER BY t.name
        ),
        '[]'::json
      )
    FROM tech_stacks ts
      JOIN techs t ON ts.id_tech = t.id
    WHERE ts.id_project = p.id
  ) as techs,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pi.id,
            'file_name',
            pi.file_name,
            'url',
            pi.url
          )
          ORDER BY pi.created_at
        ),
        '[]'::json
      )
    FROM project_images pi
    WHERE pi.id_project = p.id
  ) as images
FROM projects p
  JOIN users u ON p.id_user = u.id
  JOIN categories c ON p.id_category = c.id
WHERE p.id = $1
LIMIT 1;

-- name: GetProjectByTranslatedName :one
SELECT sqlc.embed(p),
  json_build_object(
    'id',
    u.id,
    'username',
    u.username,
    'email',
    u.email
  ) as user,
  json_build_object(
    'id',
    c.id,
    'name',
    c.name,
    'description',
    c.description
  ) as category,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pt_agg.id,
            'id_language',
            pt_agg.id_language,
            'name',
            pt_agg.name,
            'description',
            pt_agg.description,
            'language',
            json_build_object('id', l.id, 'code', l.code, 'name', l.name)
          )
          ORDER BY l.name
        ),
        '[]'::json
      )
    FROM project_translations pt_agg
      JOIN languages l ON pt_agg.id_language = l.id
    WHERE pt_agg.id_project = p.id
  ) as translations,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            t.id,
            'name',
            t.name,
            'description',
            t.description,
            'logo_url',
            t.logo_url
          )
          ORDER BY t.name
        ),
        '[]'::json
      )
    FROM tech_stacks ts
      JOIN techs t ON ts.id_tech = t.id
    WHERE ts.id_project = p.id
  ) as techs,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pi.id,
            'file_name',
            pi.file_name,
            'url',
            pi.url
          )
          ORDER BY pi.created_at
        ),
        '[]'::json
      )
    FROM project_images pi
    WHERE pi.id_project = p.id
  ) as images
FROM projects p
  JOIN users u ON p.id_user = u.id
  JOIN categories c ON p.id_category = c.id
  JOIN project_translations pt ON p.id = pt.id_project
WHERE pt.name = $1
LIMIT 1;

-- name: GetProjectsPaginated :many
SELECT sqlc.embed(p),
  json_build_object(
    'id',
    u.id,
    'username',
    u.username,
    'email',
    u.email
  ) as user,
  json_build_object(
    'id',
    c.id,
    'name',
    c.name,
    'description',
    c.description
  ) as category,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pt.id,
            'id_language',
            pt.id_language,
            'name',
            pt.name,
            'description',
            pt.description,
            'language',
            json_build_object('id', l.id, 'code', l.code, 'name', l.name)
          )
          ORDER BY l.name
        ),
        '[]'::json
      )
    FROM project_translations pt
      JOIN languages l ON pt.id_language = l.id
    WHERE pt.id_project = p.id
  ) as translations,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            t.id,
            'name',
            t.name,
            'description',
            t.description,
            'logo_url',
            t.logo_url
          )
          ORDER BY t.name
        ),
        '[]'::json
      )
    FROM tech_stacks ts
      JOIN techs t ON ts.id_tech = t.id
    WHERE ts.id_project = p.id
  ) as techs,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pi.id,
            'file_name',
            pi.file_name,
            'url',
            pi.url
          )
          ORDER BY pi.created_at
        ),
        '[]'::json
      )
    FROM project_images pi
    WHERE pi.id_project = p.id
  ) as images
FROM projects p
  JOIN users u ON p.id_user = u.id
  JOIN categories c ON p.id_category = c.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchProjectsPaginated :many
WITH relevant_projects AS (
  SELECT DISTINCT ON (p.id) p.id,
    pt.name as matched_name,
    p.created_at
  FROM projects p
    JOIN project_translations pt ON p.id = pt.id_project
  WHERE pt.name ILIKE '%' || $1::text || '%'
)
SELECT sqlc.embed(p),
  json_build_object(
    'id',
    u.id,
    'username',
    u.username,
    'email',
    u.email
  ) as user,
  json_build_object(
    'id',
    c.id,
    'name',
    c.name,
    'description',
    c.description
  ) as category,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pt.id,
            'id_language',
            pt.id_language,
            'name',
            pt.name,
            'description',
            pt.description,
            'language',
            json_build_object('id', l.id, 'code', l.code, 'name', l.name)
          )
          ORDER BY l.name
        ),
        '[]'::json
      )
    FROM project_translations pt
      JOIN languages l ON pt.id_language = l.id
    WHERE pt.id_project = p.id
  ) as translations,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            t.id,
            'name',
            t.name,
            'description',
            t.description,
            'logo_url',
            t.logo_url
          )
          ORDER BY t.name
        ),
        '[]'::json
      )
    FROM tech_stacks ts
      JOIN techs t ON ts.id_tech = t.id
    WHERE ts.id_project = p.id
  ) as techs,
  (
    SELECT COALESCE(
        json_agg(
          json_build_object(
            'id',
            pi.id,
            'file_name',
            pi.file_name,
            'url',
            pi.url
          )
          ORDER BY pi.created_at
        ),
        '[]'::json
      )
    FROM project_images pi
    WHERE pi.id_project = p.id
  ) as images
FROM projects p
  JOIN users u ON p.id_user = u.id
  JOIN categories c ON p.id_category = c.id
  JOIN relevant_projects rp ON p.id = rp.id
ORDER BY CASE
    WHEN rp.matched_name ILIKE $1::text || '%' THEN 1
    ELSE 2
  END,
  rp.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProject :one
UPDATE projects
SET id_category = COALESCE(sqlc.narg('id_category'), id_category),
  is_deployed = COALESCE(sqlc.narg('is_deployed'), is_deployed),
  is_maintained = COALESCE(sqlc.narg('is_maintained'), is_maintained),
  live_demo = COALESCE(sqlc.narg('live_demo'), live_demo),
  source_code = COALESCE(sqlc.narg('source_code'), source_code)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;
