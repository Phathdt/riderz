package handlers

import (
	"context"
	"github.com/guregu/null/v5"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"time"
)

type DriverArrivedRepo interface {
	GetTrip(ctx context.Context, tripCode string) (*tripRepo.Trip, error)
	UpdateTripStatusWithEvent(ctx context.Context, arg tripRepo.UpdateTripStatusParams, eventArg tripRepo.CreateTripEventParams) error
}

type driverArrivedHdl struct {
	repo     DriverArrivedRepo
	producer kcomp.KProducer
}

func NewDriverArrivedHdl(repo DriverArrivedRepo, producer kcomp.KProducer) *driverArrivedHdl {
	return &driverArrivedHdl{repo: repo, producer: producer}
}

func (h *driverArrivedHdl) Response(ctx context.Context, tripCode string, data *dto.DriverArrivedData) error {
	trip, err := h.repo.GetTrip(ctx, tripCode)
	if err != nil {
		return err
	}

	if !trip.DriverID.Equal(null.IntFrom(data.DriverID)) {
		return err
	}

	tripData := tripRepo.UpdateTripStatusParams{
		TripCode: tripCode,
		Status:   domain.TripStatusDriverArrived,
	}

	eventData := tripRepo.CreateTripEventParams{
		TripID:    trip.ID,
		TripCode:  trip.TripCode,
		EventType: domain.TripEventTypeDriverArrived,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			DriverArrived: &domain.DriverArrivedData{ArrivalLocation: domain.Location{
				Latitude:  data.Lat,
				Longitude: data.Long,
			}},
		},
	}

	if err = h.repo.UpdateTripStatusWithEvent(ctx, tripData, eventData); err != nil {
		return err
	}

	message := domain.TripDriverArrivedMessage{
		EventType: string(domain.TripEventTypeDriverArrived),
		TripCode:  tripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID         int64           `json:"trip_id"`
			DriverID       int64           `json:"driver_id"`
			ArriveLocation domain.Location `json:"arrive_location"`
		}{
			TripID:   trip.ID,
			DriverID: data.DriverID,
			ArriveLocation: domain.Location{
				Latitude:  data.Lat,
				Longitude: data.Long,
			},
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripTopicDriverArrived), "", message)

	return nil
}
