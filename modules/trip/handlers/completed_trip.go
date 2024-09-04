package handlers

import (
	"context"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"time"
)

type CompletedTripRepo interface {
	CompleteTrip(ctx context.Context, arg tripRepo.CompleteTripParams) error
	GetTrip(ctx context.Context, tripCode string) (*tripRepo.Trip, error)
	CreateTripEvent(ctx context.Context, arg tripRepo.CreateTripEventParams) (int64, error)
}

type completeTripHdl struct {
	repo     CompletedTripRepo
	producer kcomp.KProducer
}

func NewCompleteTripHdl(repo CompletedTripRepo, producer kcomp.KProducer) *completeTripHdl {
	return &completeTripHdl{repo: repo, producer: producer}
}

func (h *completeTripHdl) Response(ctx context.Context, tripCode string, data *dto.CompleteTripData) error {
	trip, err := h.repo.GetTrip(ctx, tripCode)
	if err != nil {
		return err
	}

	if err = h.repo.CompleteTrip(ctx, tripRepo.CompleteTripParams{
		TripCode: tripCode,
		Status:   domain.TripStatusCompleted,
	}); err != nil {
		return err
	}

	if _, err = h.repo.CreateTripEvent(ctx, tripRepo.CreateTripEventParams{
		TripID:    trip.ID,
		TripCode:  trip.TripCode,
		EventType: domain.TripEventTypeCompleted,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			TripCompleted: &domain.TripCompletedData{
				EndLocation: domain.Location{
					Latitude:  data.Lat,
					Longitude: data.Long,
				},
				EndTime:      time.Now(),
				TripDistance: 0,
				TripDuration: 0,
			},
		},
	}); err != nil {
		return err
	}

	message := domain.TripEndMessage{
		EventType: string(domain.TripEventTypeStarted),
		TripCode:  tripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID      int64           `json:"trip_id"`
			EndLocation domain.Location `json:"end_location"`
		}{
			TripID: trip.ID,
			EndLocation: domain.Location{
				Latitude:  data.Lat,
				Longitude: data.Long,
			},
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripCompletedTopic), "", message)

	return nil
}
