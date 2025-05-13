-- name: StoreClick :exec
INSERT INTO
    clicks (id, url_id, ip, clicked_at)
VALUES
    ($1, $2, $3, $4);

-- name: CountClicks :one
SELECT
    count(*)
FROM
    clicks
WHERE
    url_id = $1;

-- name: DeleteClicks :one
WITH deleted AS (
    DELETE FROM
        clicks
    WHERE
        clicked_at <= $1 RETURNING *
)
SELECT
    count(*)
FROM
    deleted;
