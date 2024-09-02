// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package tripRepo

import (
	"time"

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
