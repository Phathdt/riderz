package handlers

import (
	"context"
	"fmt"
	"github.com/jaevor/go-nanoid"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/dto"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"riderz/shared/common"
	"time"
)

type RequestTripRepo interface {
	CreateTrip(ctx context.Context, arg tripRepo.CreateTripParams) (int64, error)
	CreateTripEvent(ctx context.Context, arg tripRepo.CreateTripEventParams) (int64, error)
}

type requestTripHdl struct {
	producer kcomp.KProducer
	repo     RequestTripRepo
}

func NewRequestTripHdl(producer kcomp.KProducer, repo RequestTripRepo) *requestTripHdl {
	return &requestTripHdl{producer: producer, repo: repo}
}

func (h *requestTripHdl) Response(ctx context.Context, data *dto.CreateTripData) (string, error) {
	tripCode := h.generateIdempotent()

	tripID, err := h.repo.CreateTrip(ctx, tripRepo.CreateTripParams{
		TripCode:       tripCode,
		UserID:         data.UserID,
		Status:         domain.TripStatusRequested,
		PickupAddress:  data.PickupAddress,
		DropoffAddress: data.DropoffAddress,
		PickupLong:     data.PickupLon,
		PickupLat:      data.PickupLat,
		DropoffLong:    data.DropoffLon,
		DropoffLat:     data.DropoffLat,
	})
	if err != nil {
		return "", err
	}

	_, err = h.repo.CreateTripEvent(ctx, tripRepo.CreateTripEventParams{
		TripID:    tripID,
		TripCode:  tripCode,
		EventType: domain.TripEventTypeRequested,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			TripRequested: &domain.TripRequestedData{
				UserID: data.UserID,
				PickupLocation: domain.Location{
					Latitude:  data.PickupLat,
					Longitude: data.PickupLon,
					Address:   data.PickupAddress,
				},
				DropoffLocation: domain.Location{
					Latitude:  data.DropoffLat,
					Longitude: data.DropoffLon,
					Address:   data.DropoffAddress,
				},
			},
		},
	})

	message := domain.TripRequestedMessage{
		EventType: string(domain.TripEventTypeRequested),
		TripCode:  tripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID          int64           `json:"trip_id"`
			UserID          int64           `json:"user_id"`
			PickupLocation  domain.Location `json:"pickup_location"`
			DropoffLocation domain.Location `json:"dropoff_location"`
		}{
			TripID: tripID,
			UserID: data.UserID,
			PickupLocation: domain.Location{
				Latitude:  data.PickupLat,
				Longitude: data.PickupLon,
				Address:   data.PickupAddress,
			},
			DropoffLocation: domain.Location{
				Latitude:  data.DropoffLat,
				Longitude: data.DropoffLon,
				Address:   data.DropoffAddress,
			},
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripTopicRequested), "", message)

	return tripCode, nil
}

func (h *requestTripHdl) generateIdempotent() string {
	canonicID, _ := nanoid.CustomASCII(common.ALPHABET, 10)
	idempotent := canonicID()

	now := time.Now()

	return fmt.Sprintf("%s%s", now.Format("20060102"), idempotent)
}
