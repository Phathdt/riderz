package domain

import "time"

type TripRequestedMessage struct {
	EventType string    `json:"event_type"`
	TripCode  string    `json:"trip_code"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		TripID          int64    `json:"trip_id"`
		UserID          int64    `json:"user_id"`
		PickupLocation  Location `json:"pickup_location"`
		DropoffLocation Location `json:"dropoff_location"`
	} `json:"data"`
}

type TripDriverAssignedMessage struct {
	EventType string    `json:"event_type"`
	TripCode  string    `json:"trip_code"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		TripID              int64     `json:"trip_id"`
		DriverID            int64     `json:"driver_id"`
		EstimatedPickupTime time.Time `json:"estimated_pickup_time"`
	} `json:"data"`
}

type TripDriverArrivedMessage struct {
	EventType string    `json:"event_type"`
	TripCode  string    `json:"trip_code"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		TripID         int64    `json:"trip_id"`
		DriverID       int64    `json:"driver_id"`
		ArriveLocation Location `json:"arrive_location"`
	} `json:"data"`
}
type TripCancelMessage struct {
	EventType string    `json:"event_type"`
	TripCode  string    `json:"trip_code"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		TripID             int64     `json:"trip_id"`
		CancellationReason string    `json:"cancellation_reason"`
		CancelledBy        string    `json:"cancelled_by"`
		CancellationTime   time.Time `json:"cancellation_time"`
	} `json:"data"`
}
