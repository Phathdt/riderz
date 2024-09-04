package domain

type TripStatus string

const (
	TripStatusRequested      TripStatus = "REQUESTED"
	TripStatusDriverAssigned TripStatus = "DRIVER_ASSIGNED"
	TripStatusDriverArrived  TripStatus = "DRIVER_ARRIVED"
	TripStatusStarted        TripStatus = "STARTED"
	TripStatusCompleted      TripStatus = "COMPLETED"
	TripStatusCancelled      TripStatus = "CANCELLED"
)

type TripEventType string

const (
	TripEventTypeRequested      TripEventType = "TRIP_REQUESTED"
	TripEventTypeDriverAssigned TripEventType = "DRIVER_ASSIGNED"
	TripEventTypeDriverArrived  TripEventType = "DRIVER_ARRIVED"
	TripEventTypeStarted        TripEventType = "TRIP_STARTED"
	TripEventTypeCompleted      TripEventType = "TRIP_COMPLETED"
	TripEventTypeCancelled      TripEventType = "TRIP_CANCELLED"
)
