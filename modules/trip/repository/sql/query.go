package tripRepo

type TripWithEvents struct {
	*Trip  `json:",inline"`
	Events []*TripEvent `json:"events,omitempty"`
}
