-- name: CreateLocation :exec
INSERT INTO locations (user_id, geometry)
VALUES (@user_id, ST_SetSRID(ST_MakePoint(@lat::decimal, @long::decimal), 4326))
ON CONFLICT (user_id)
DO UPDATE SET
    geometry = EXCLUDED.geometry,
    updated_at = now();

-- name: GetLocation :one
SELECT * FROM locations
WHERE user_id = $1 LIMIT 1;

-- name: ListLocations :many
SELECT id, user_id, geometry, ST_Distance(geometry, ST_SetSRID(ST_MakePoint(@lat::decimal, @long::decimal), 4326))::decimal AS distance
FROM locations
ORDER BY geometry <-> ST_SetSRID(ST_MakePoint(@lat::decimal, @long::decimal), 4326)
LIMIT @size;
