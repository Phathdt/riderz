package dto

type CreateTripData struct {
	UserID         int64   `json:"-"`
	PickupLat      float64 `json:"pickup_lat" validate:"required"`
	PickupLon      float64 `json:"pickup_lon" validate:"required"`
	PickupAddress  string  `json:"pickup_address" validate:"required"`
	DropoffLat     float64 `json:"dropoff_lat" validate:"required"`
	DropoffLon     float64 `json:"dropoff_lon" validate:"required"`
	DropoffAddress string  `json:"dropoff_address" validate:"required"`
}
