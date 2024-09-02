-- name: CreateTrip :one
INSERT INTO trips (
    trip_code,
    user_id,
    status,
    pickup_location,
    pickup_address,
    dropoff_location,
    dropoff_address
) VALUES (
    $1,
    $2,
    $3,
    ST_SetSRID(ST_MakePoint(@pickup_long::decimal, @pickup_lat::decimal), 4326),
    $4,
    ST_SetSRID(ST_MakePoint(@dropoff_long::decimal, @dropoff_lat::decimal), 4326),
    $5
) RETURNING *;

-- name: GetTrip :one
SELECT * FROM trips
WHERE id = $1;

-- name: ListTrips :many
SELECT *
FROM trips
WHERE user_id = $1
ORDER BY id DESC;

-- name: UpdateTripStatus :exec
UPDATE trips
SET status = $2
WHERE id = $1;

-- name: CreateTripEvent :one
INSERT INTO trip_events (
    trip_id,
    event_type,
    status,
    event_data
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;
