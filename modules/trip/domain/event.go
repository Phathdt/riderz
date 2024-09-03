package domain

import (
	"time"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

type TripRequestedData struct {
	UserID          int64    `json:"user_id"`
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`
}

type DriverAssignedData struct {
	DriverID            int64     `json:"driver_id"`
	EstimatedPickupTime time.Time `json:"estimated_pickup_time"`
}

type DriverArrivedData struct {
	ArrivalLocation Location `json:"arrival_location"`
}

type TripStartedData struct {
	StartLocation Location  `json:"start_location"`
	StartTime     time.Time `json:"start_time"`
}

type TripCompletedData struct {
	EndLocation  Location  `json:"end_location"`
	EndTime      time.Time `json:"end_time"`
	TripDistance float64   `json:"trip_distance"`
	TripDuration float64   `json:"trip_duration"`
}

type TripCancelledData struct {
	CancellationReason string    `json:"cancellation_reason"`
	CancelledBy        string    `json:"cancelled_by"`
	CancellationTime   time.Time `json:"cancellation_time"`
}

type TripEventData struct {
	TripRequested  *TripRequestedData  `json:"trip_requested,omitempty"`
	DriverAssigned *DriverAssignedData `json:"driver_assigned,omitempty"`
	DriverArrived  *DriverArrivedData  `json:"driver_arrived,omitempty"`
	TripStarted    *TripStartedData    `json:"trip_started,omitempty"`
	TripCompleted  *TripCompletedData  `json:"trip_completed,omitempty"`
	TripCancelled  *TripCancelledData  `json:"trip_cancelled,omitempty"`
}
