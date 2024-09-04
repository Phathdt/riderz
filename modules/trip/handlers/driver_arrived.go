package handlers

import (
	"context"
	"fmt"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
)

type DriverArrivedRepo interface {
	GetTrip(ctx context.Context, tripCode string) (*tripRepo.Trip, error)
	UpdateTripStatus(ctx context.Context, arg tripRepo.UpdateTripStatusParams) error
	CreateTripEvent(ctx context.Context, arg tripRepo.CreateTripEventParams) (int64, error)
}

type driverArrivedHdl struct {
	repo DriverArrivedRepo
}

func (h *driverArrivedHdl) Response(ctx context.Context, tripCode string, data *dto.DriverArrivedData) error {
	trip, err := h.repo.GetTrip(ctx, tripCode)
	if err != nil {
		return err
	}

	fmt.Println(trip)

	//if trip.DriverID != data.DriverID {
	//
	//}

	return nil
}
