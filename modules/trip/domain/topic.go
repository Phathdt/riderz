package domain

type TripTopic string

const (
	TripTopicRequested    TripTopic = "trips.requested"
	TripTopicAssignDriver TripTopic = "trips.assign_driver"
)
