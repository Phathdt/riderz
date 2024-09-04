// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package tripRepo

import (
	"time"

	null "github.com/guregu/null/v5"
	"riderz/modules/trip/domain"
	"riderz/shared/common"
)

type Trip struct {
	ID              int64             `db:"id" json:"id"`
	TripCode        string            `db:"trip_code" json:"trip_code"`
	UserID          int64             `db:"user_id" json:"user_id"`
	DriverID        null.Int64        `db:"driver_id" json:"driver_id"`
	Status          domain.TripStatus `db:"status" json:"status"`
	PickupLocation  common.PointS     `db:"pickup_location" json:"pickup_location"`
	PickupAddress   string            `db:"pickup_address" json:"pickup_address"`
	DropoffLocation common.PointS     `db:"dropoff_location" json:"dropoff_location"`
	DropoffAddress  string            `db:"dropoff_address" json:"dropoff_address"`
	RequestTime     time.Time         `db:"request_time" json:"request_time"`
	StartTime       null.Time         `db:"start_time" json:"start_time"`
	EndTime         null.Time         `db:"end_time" json:"end_time"`
	Price           null.Float        `db:"price" json:"price"`
	Distance        null.Float        `db:"distance" json:"distance"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
}

type TripEvent struct {
	ID        int64                `db:"id" json:"id"`
	TripID    int64                `db:"trip_id" json:"trip_id"`
	TripCode  string               `db:"trip_code" json:"trip_code"`
	EventType domain.TripEventType `db:"event_type" json:"event_type"`
	EventTime time.Time            `db:"event_time" json:"event_time"`
	EventData domain.TripEventData `db:"event_data" json:"event_data"`
	CreatedAt time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt time.Time            `db:"updated_at" json:"updated_at"`
}
