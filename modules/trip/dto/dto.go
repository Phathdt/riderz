package dto

import "time"

type CreateTripData struct {
	UserID         int64   `json:"-"`
	PickupLat      float64 `json:"pickup_lat" validate:"required"`
	PickupLon      float64 `json:"pickup_lon" validate:"required"`
	PickupAddress  string  `json:"pickup_address" validate:"required"`
	DropoffLat     float64 `json:"dropoff_lat" validate:"required"`
	DropoffLon     float64 `json:"dropoff_lon" validate:"required"`
	DropoffAddress string  `json:"dropoff_address" validate:"required"`
}

type DriverArrivedData struct {
	DriverID int64   `json:"-"`
	Lat      float64 `json:"lat" validate:"required"`
	Long     float64 `json:"long" validate:"required"`
}

type CancelTripData struct {
	UserID int64  `json:"-"`
	Reason string `json:"reason" validate:"required"`
}

type StartTripData struct {
	DriverID  int64     `json:"-"`
	Lat       float64   `json:"lat" validate:"required"`
	Long      float64   `json:"long" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
}

type EndTripData struct {
	DriverID int64     `json:"-"`
	Lat      float64   `json:"lat" validate:"required"`
	Long     float64   `json:"long" validate:"required"`
	EndTime  time.Time `json:"end_time" validate:"required"`
}
