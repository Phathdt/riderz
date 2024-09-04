package domain

type TripTopic string

const (
	TripTopicRequested     TripTopic = "trips.requested"
	TripTopicAssignDriver  TripTopic = "trips.assign_driver"
	TripTopicDriverArrived TripTopic = "trips.driver_arrived"
	TripCancelTopic        TripTopic = "trips.cancel"
	TripStartTopic         TripTopic = "trips.start"
	TripCompletedTopic     TripTopic = "trips.start"
)
