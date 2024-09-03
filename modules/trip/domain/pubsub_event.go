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
