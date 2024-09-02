package dto

import "riderz/shared/common"

type TripEventData struct {
	UserId         *int64         `json:"user_id"`
	PickupLocation *common.PointS `json:"pickup_location"`
	DriverId       *int64         `json:"driver_id"`
	StartLocation  *common.PointS `json:"start_location"`
	EndLocation    *common.PointS `json:"end_location"`
	Amount         *float32       `json:"amount"`
}
