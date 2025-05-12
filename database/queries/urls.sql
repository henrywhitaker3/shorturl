-- name: CreateUrl :one
WITH alias AS (
    SELECT
        *
    FROM
        alias_buffer
    LIMIT
        1 FOR
    UPDATE
),
inserted AS (
    INSERT INTO
        urls (id, alias, url, domain)
    VALUES
        (
            $1,
            (
                SELECT
                    alias
                FROM
                    alias
                LIMIT
                    1
            ), $2, $3
        ) RETURNING *
), deleted AS (
    DELETE FROM
        alias_buffer
    WHERE
        alias = (
            SELECT
                alias
            FROM
                alias
            LIMIT
                1
        )
)
SELECT
    *
FROM
    inserted;

-- name: GetUrl :one
SELECT
    *
FROM
    urls
WHERE
    id = $1;

-- name: GetUrlByAlias :one
SELECT
    *
FROM
    urls
WHERE
    alias = $1;

-- name: CountUrls :one
SELECT
    count(*)
FROM
    urls;

-- name: GetUrlsByAliases :many
SELECT
    id,
    alias
FROM
    urls
WHERE
    alias = ANY(sqlc.Slice(aliases) :: text []);
