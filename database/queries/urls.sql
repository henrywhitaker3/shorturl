-- name: CreateUrl :one
INSERT INTO
    urls (id, alias, url, domain)
VALUES
    ($1, $2, $3, $4) RETURNING *;

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
