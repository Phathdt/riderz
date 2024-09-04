package handlers

import (
	"context"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"time"
)

type CancelTripRepo interface {
	GetTrip(ctx context.Context, tripCode string) (*tripRepo.Trip, error)
	UpdateTripStatus(ctx context.Context, arg tripRepo.UpdateTripStatusParams) error
	CreateTripEvent(ctx context.Context, arg tripRepo.CreateTripEventParams) (int64, error)
}

type cancelTripHdl struct {
	repo     CancelTripRepo
	producer kcomp.KProducer
}

func NewCancelTripHdl(repo CancelTripRepo, producer kcomp.KProducer) *cancelTripHdl {
	return &cancelTripHdl{repo: repo, producer: producer}
}

func (h *cancelTripHdl) Response(ctx context.Context, tripCode string, data *dto.CancelTripData) error {
	trip, err := h.repo.GetTrip(ctx, tripCode)
	if err != nil {
		return err
	}

	if err = h.repo.UpdateTripStatus(ctx, tripRepo.UpdateTripStatusParams{
		TripCode: tripCode,
		Status:   domain.TripStatusCancelled,
	}); err != nil {
		return err
	}

	if _, err = h.repo.CreateTripEvent(ctx, tripRepo.CreateTripEventParams{
		TripID:    trip.ID,
		TripCode:  trip.TripCode,
		EventType: domain.TripEventTypeCancelled,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			TripCancelled: &domain.TripCancelledData{
				CancellationReason: data.Reason,
				CancelledBy:        "user",
				CancellationTime:   time.Now(),
			},
		},
	}); err != nil {
		return err
	}

	message := domain.TripCancelMessage{
		EventType: string(domain.TripEventTypeCancelled),
		TripCode:  tripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID             int64     `json:"trip_id"`
			CancellationReason string    `json:"cancellation_reason"`
			CancelledBy        string    `json:"cancelled_by"`
			CancellationTime   time.Time `json:"cancellation_time"`
		}{
			TripID:             trip.ID,
			CancellationReason: data.Reason,
			CancelledBy:        "user",
			CancellationTime:   time.Now(),
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripCancelTopic), "", message)

	return nil
}
