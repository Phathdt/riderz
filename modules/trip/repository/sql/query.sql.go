// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package tripRepo

import (
	"context"
)

const createTrip = `-- name: CreateTrip :one
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
    ST_SetSRID(ST_MakePoint($6::decimal, $7::decimal), 4326),
    $4,
    ST_SetSRID(ST_MakePoint($8::decimal, $9::decimal), 4326),
    $5
) RETURNING id, trip_code, user_id, driver_id, status, pickup_location, pickup_address, dropoff_location, dropoff_address, request_time, start_time, end_time, price, distance, created_at, updated_at
`

type CreateTripParams struct {
	TripCode       string  `db:"trip_code" json:"trip_code"`
	UserID         int64   `db:"user_id" json:"user_id"`
	Status         string  `db:"status" json:"status"`
	PickupAddress  string  `db:"pickup_address" json:"pickup_address"`
	DropoffAddress string  `db:"dropoff_address" json:"dropoff_address"`
	PickupLong     float64 `db:"pickup_long" json:"pickup_long"`
	PickupLat      float64 `db:"pickup_lat" json:"pickup_lat"`
	DropoffLong    float64 `db:"dropoff_long" json:"dropoff_long"`
	DropoffLat     float64 `db:"dropoff_lat" json:"dropoff_lat"`
}

func (q *Queries) CreateTrip(ctx context.Context, arg CreateTripParams) (*Trip, error) {
	row := q.db.QueryRow(ctx, createTrip,
		arg.TripCode,
		arg.UserID,
		arg.Status,
		arg.PickupAddress,
		arg.DropoffAddress,
		arg.PickupLong,
		arg.PickupLat,
		arg.DropoffLong,
		arg.DropoffLat,
	)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.TripCode,
		&i.UserID,
		&i.DriverID,
		&i.Status,
		&i.PickupLocation,
		&i.PickupAddress,
		&i.DropoffLocation,
		&i.DropoffAddress,
		&i.RequestTime,
		&i.StartTime,
		&i.EndTime,
		&i.Price,
		&i.Distance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getTrip = `-- name: GetTrip :one
SELECT id, trip_code, user_id, driver_id, status, pickup_location, pickup_address, dropoff_location, dropoff_address, request_time, start_time, end_time, price, distance, created_at, updated_at FROM trips
WHERE id = $1
`

func (q *Queries) GetTrip(ctx context.Context, id int64) (*Trip, error) {
	row := q.db.QueryRow(ctx, getTrip, id)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.TripCode,
		&i.UserID,
		&i.DriverID,
		&i.Status,
		&i.PickupLocation,
		&i.PickupAddress,
		&i.DropoffLocation,
		&i.DropoffAddress,
		&i.RequestTime,
		&i.StartTime,
		&i.EndTime,
		&i.Price,
		&i.Distance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listTrips = `-- name: ListTrips :many
SELECT id, trip_code, user_id, driver_id, status, pickup_location, pickup_address, dropoff_location, dropoff_address, request_time, start_time, end_time, price, distance, created_at, updated_at
FROM trips
WHERE user_id = $1
ORDER BY id DESC
`

func (q *Queries) ListTrips(ctx context.Context, userID int64) ([]*Trip, error) {
	rows, err := q.db.Query(ctx, listTrips, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Trip
	for rows.Next() {
		var i Trip
		if err := rows.Scan(
			&i.ID,
			&i.TripCode,
			&i.UserID,
			&i.DriverID,
			&i.Status,
			&i.PickupLocation,
			&i.PickupAddress,
			&i.DropoffLocation,
			&i.DropoffAddress,
			&i.RequestTime,
			&i.StartTime,
			&i.EndTime,
			&i.Price,
			&i.Distance,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTripStatus = `-- name: UpdateTripStatus :exec
UPDATE trips
SET status = $2
WHERE id = $1
`

type UpdateTripStatusParams struct {
	ID     int64  `db:"id" json:"id"`
	Status string `db:"status" json:"status"`
}

func (q *Queries) UpdateTripStatus(ctx context.Context, arg UpdateTripStatusParams) error {
	_, err := q.db.Exec(ctx, updateTripStatus, arg.ID, arg.Status)
	return err
}
