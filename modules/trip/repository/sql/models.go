// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package tripRepo

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"riderz/modules/trip/dto"
	"riderz/shared/common"
)

type Trip struct {
	ID              int64         `db:"id" json:"id"`
	TripCode        string        `db:"trip_code" json:"trip_code"`
	UserID          int64         `db:"user_id" json:"user_id"`
	DriverID        int64         `db:"driver_id" json:"driver_id"`
	Status          string        `db:"status" json:"status"`
	PickupLocation  common.PointS `db:"pickup_location" json:"pickup_location"`
	PickupAddress   string        `db:"pickup_address" json:"pickup_address"`
	DropoffLocation common.PointS `db:"dropoff_location" json:"dropoff_location"`
	DropoffAddress  string        `db:"dropoff_address" json:"dropoff_address"`
	RequestTime     time.Time     `db:"request_time" json:"request_time"`
	StartTime       time.Time     `db:"start_time" json:"start_time"`
	EndTime         time.Time     `db:"end_time" json:"end_time"`
	Price           float64       `db:"price" json:"price"`
	Distance        float64       `db:"distance" json:"distance"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
}

type TripEvent struct {
	ID        int64             `db:"id" json:"id"`
	TripID    int64             `db:"trip_id" json:"trip_id"`
	EventType string            `db:"event_type" json:"event_type"`
	Status    pgtype.Text       `db:"status" json:"status"`
	EventData dto.TripEventData `db:"event_data" json:"event_data"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
}
