package dto

type UpdateLocationRequest struct {
	UserId    int64   `json:"user_id" avro:"user_id"`
	Latitude  float64 `json:"latitude" avro:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" avro:"longitude" validate:"required"`
}

type SearchNearLocation struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude"  validate:"required"`
}
