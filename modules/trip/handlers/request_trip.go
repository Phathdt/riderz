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
	CreateTripAndEvent(ctx context.Context, tripArg tripRepo.CreateTripParams, eventArg tripRepo.CreateTripEventParams) (tripID int64, err error)
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

	now := time.Now()
	createTripData := tripRepo.CreateTripParams{
		TripCode:       tripCode,
		UserID:         data.UserID,
		Status:         domain.TripStatusRequested,
		PickupAddress:  data.PickupAddress,
		DropoffAddress: data.DropoffAddress,
		PickupLong:     data.PickupLon,
		PickupLat:      data.PickupLat,
		DropoffLong:    data.DropoffLon,
		DropoffLat:     data.DropoffLat,
	}

	createTripEventData := tripRepo.CreateTripEventParams{
		TripCode:  tripCode,
		EventType: domain.TripEventTypeRequested,
		EventTime: now,
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
		}}

	tripID, err := h.repo.CreateTripAndEvent(ctx, createTripData, createTripEventData)
	if err != nil {
		return "", err
	}

	message := domain.TripRequestedMessage{
		EventType: string(domain.TripEventTypeRequested),
		TripCode:  tripCode,
		Timestamp: now,
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
