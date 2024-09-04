package handlers

import (
	"context"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"time"
)

type StartTripRepo interface {
	StartTrip(ctx context.Context, arg tripRepo.StartTripParams) error
	GetTrip(ctx context.Context, tripCode string) (*tripRepo.Trip, error)
	CreateTripEvent(ctx context.Context, arg tripRepo.CreateTripEventParams) (int64, error)
}

type startTripHdl struct {
	repo     StartTripRepo
	producer kcomp.KProducer
}

func NewStartTripHdl(repo StartTripRepo, producer kcomp.KProducer) *startTripHdl {
	return &startTripHdl{repo: repo, producer: producer}
}

func (h *startTripHdl) Response(ctx context.Context, tripCode string, data *dto.StartTripData) error {
	trip, err := h.repo.GetTrip(ctx, tripCode)
	if err != nil {
		return err
	}

	if err := h.repo.StartTrip(ctx, tripRepo.StartTripParams{
		TripCode: tripCode,
		Status:   domain.TripStatusStarted,
	}); err != nil {
		return err
	}

	if _, err := h.repo.CreateTripEvent(ctx, tripRepo.CreateTripEventParams{
		TripID:    trip.ID,
		TripCode:  trip.TripCode,
		EventType: domain.TripEventTypeStarted,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			TripStarted: &domain.TripStartedData{
				StartLocation: domain.Location{
					Latitude:  data.Lat,
					Longitude: data.Long,
				},
				StartTime: time.Now(),
			},
		},
	}); err != nil {
		return err
	}

	message := domain.TripStartMessage{
		EventType: string(domain.TripEventTypeStarted),
		TripCode:  tripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID        int64           `json:"trip_id"`
			StartLocation domain.Location `json:"start_location"`
		}{
			TripID: trip.ID,
			StartLocation: domain.Location{
				Latitude:  data.Lat,
				Longitude: data.Long,
			},
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripStartTopic), "", message)

	return nil
}
