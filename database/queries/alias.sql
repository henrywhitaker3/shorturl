-- name: CreateAlias :one
INSERT INTO
    aliases (alias, used)
VALUES
    ($1, false) RETURNING *;

-- name: GetFreeAlias :one
SELECT
    *
FROM
    aliases
WHERE
    used = false
LIMIT
    1 FOR
UPDATE
;

-- name: MarkAliasUsed :exec
UPDATE
    aliases
SET
    used = TRUE
WHERE
    alias = $1;

-- name: CountFreeAliases :one
SELECT
    count(*)
FROM
    aliases
WHERE
    used = false;

-- name: CountAliases :one
SELECT
    count(*)
FROM
    aliases;

-- name: GetAliases :many
SELECT
    alias
FROM
    aliases
WHERE
    alias = ANY(sqlc.Slice(aliases) :: text []);
