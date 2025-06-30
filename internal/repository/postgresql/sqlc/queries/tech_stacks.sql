-- name: CreateTechStack :one
INSERT INTO tech_stacks (id_project, id_tech)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTechStack :one
UPDATE tech_stacks
SET id_project = COALESCE($3, id_project),
  id_tech = COALESCE($4, id_tech)
WHERE id_project = $1
  AND id_tech = $2
RETURNING *;

-- name: DeleteTechStack :exec
DELETE FROM tech_stacks
WHERE id_project = $1
  AND id_tech = $2;
