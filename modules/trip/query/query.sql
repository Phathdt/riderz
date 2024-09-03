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
) RETURNING id;

-- name: GetTrip :one
SELECT * FROM trips
WHERE trip_code = $1 and user_id = $2;

-- name: ListTrips :many
SELECT *
FROM trips
WHERE user_id = $1
ORDER BY id DESC;

-- name: UpdateTripStatus :exec
UPDATE trips
SET status = $2
WHERE trip_code = $1;

-- name: CreateTripEvent :one
INSERT INTO trip_events (
    trip_id,
    trip_code,
    event_type,
    event_time,
    event_data
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING id;

-- name: UpdateDriverId :exec
UPDATE trips
SET driver_id = $2
WHERE trip_code = $1;

-- name: ListTripEvents :many
SELECT *
FROM trip_events
WHERE trip_code = $1
ORDER BY event_time DESC;
