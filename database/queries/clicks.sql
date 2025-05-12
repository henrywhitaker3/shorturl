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
