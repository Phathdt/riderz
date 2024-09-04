// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package tripRepo

import (
	"context"
	"time"

	"riderz/modules/trip/domain"
)

const assignDriver = `-- name: AssignDriver :exec
UPDATE trips
SET driver_id = $2, status = $3
WHERE trip_code = $1
`

type AssignDriverParams struct {
	TripCode string            `db:"trip_code" json:"trip_code"`
	DriverID *int64            `db:"driver_id" json:"driver_id"`
	Status   domain.TripStatus `db:"status" json:"status"`
}

func (q *Queries) AssignDriver(ctx context.Context, arg AssignDriverParams) error {
	_, err := q.db.Exec(ctx, assignDriver, arg.TripCode, arg.DriverID, arg.Status)
	return err
}

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
) RETURNING id
`

type CreateTripParams struct {
	TripCode       string            `db:"trip_code" json:"trip_code"`
	UserID         int64             `db:"user_id" json:"user_id"`
	Status         domain.TripStatus `db:"status" json:"status"`
	PickupAddress  string            `db:"pickup_address" json:"pickup_address"`
	DropoffAddress string            `db:"dropoff_address" json:"dropoff_address"`
	PickupLong     float64           `db:"pickup_long" json:"pickup_long"`
	PickupLat      float64           `db:"pickup_lat" json:"pickup_lat"`
	DropoffLong    float64           `db:"dropoff_long" json:"dropoff_long"`
	DropoffLat     float64           `db:"dropoff_lat" json:"dropoff_lat"`
}

func (q *Queries) CreateTrip(ctx context.Context, arg CreateTripParams) (int64, error) {
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
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createTripEvent = `-- name: CreateTripEvent :one
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
) RETURNING id
`

type CreateTripEventParams struct {
	TripID    int64                `db:"trip_id" json:"trip_id"`
	TripCode  string               `db:"trip_code" json:"trip_code"`
	EventType domain.TripEventType `db:"event_type" json:"event_type"`
	EventTime time.Time            `db:"event_time" json:"event_time"`
	EventData domain.TripEventData `db:"event_data" json:"event_data"`
}

func (q *Queries) CreateTripEvent(ctx context.Context, arg CreateTripEventParams) (int64, error) {
	row := q.db.QueryRow(ctx, createTripEvent,
		arg.TripID,
		arg.TripCode,
		arg.EventType,
		arg.EventTime,
		arg.EventData,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getTrip = `-- name: GetTrip :one
SELECT id, trip_code, user_id, driver_id, status, pickup_location, pickup_address, dropoff_location, dropoff_address, request_time, start_time, end_time, price, distance, created_at, updated_at FROM trips
WHERE trip_code = $1
`

func (q *Queries) GetTrip(ctx context.Context, tripCode string) (*Trip, error) {
	row := q.db.QueryRow(ctx, getTrip, tripCode)
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

const getTripByUser = `-- name: GetTripByUser :one
SELECT id, trip_code, user_id, driver_id, status, pickup_location, pickup_address, dropoff_location, dropoff_address, request_time, start_time, end_time, price, distance, created_at, updated_at FROM trips
WHERE trip_code = $1 and user_id = $2
`

type GetTripByUserParams struct {
	TripCode string `db:"trip_code" json:"trip_code"`
	UserID   int64  `db:"user_id" json:"user_id"`
}

func (q *Queries) GetTripByUser(ctx context.Context, arg GetTripByUserParams) (*Trip, error) {
	row := q.db.QueryRow(ctx, getTripByUser, arg.TripCode, arg.UserID)
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

const listTripEvents = `-- name: ListTripEvents :many
SELECT id, trip_id, trip_code, event_type, event_time, event_data, created_at, updated_at
FROM trip_events
WHERE trip_code = $1
ORDER BY event_time DESC
`

func (q *Queries) ListTripEvents(ctx context.Context, tripCode string) ([]*TripEvent, error) {
	rows, err := q.db.Query(ctx, listTripEvents, tripCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*TripEvent
	for rows.Next() {
		var i TripEvent
		if err := rows.Scan(
			&i.ID,
			&i.TripID,
			&i.TripCode,
			&i.EventType,
			&i.EventTime,
			&i.EventData,
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
WHERE trip_code = $1
`

type UpdateTripStatusParams struct {
	TripCode string            `db:"trip_code" json:"trip_code"`
	Status   domain.TripStatus `db:"status" json:"status"`
}

func (q *Queries) UpdateTripStatus(ctx context.Context, arg UpdateTripStatusParams) error {
	_, err := q.db.Exec(ctx, updateTripStatus, arg.TripCode, arg.Status)
	return err
}
