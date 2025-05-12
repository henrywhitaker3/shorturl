-- name: InsertAliasBuffer :one
WITH inserted AS (
    INSERT INTO
        alias_buffer (alias)
    VALUES
        ($1) ON CONFLICT DO NOTHING RETURNING *
)
SELECT
    count(*)
FROM
    inserted;

-- name: GetLongestAliasBuffer :one
SELECT
    len(alias)
FROM
    alias_buffer
ORDER BY
    len(alias) DESC
LIMIT
    1;

-- name: CountAliasBuffer :one
SELECT
    count(*)
FROM
    alias_buffer;
