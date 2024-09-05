package handlers

import (
	"github.com/guregu/null/v5"
	"golang.org/x/net/context"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/modules/trip/domain"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"time"
)

type AssignDriverLocationRepo interface {
	ListLocations(ctx context.Context, arg locationRepo.ListLocationsParams) ([]*locationRepo.ListLocationsRow, error)
}

type AssignDriverRepo interface {
	AssignDriverWithEvent(ctx context.Context, arg tripRepo.AssignDriverParams, argEvent tripRepo.CreateTripEventParams) error
}

type assignDriverHdl struct {
	repo     AssignDriverRepo
	lRepo    AssignDriverLocationRepo
	producer kcomp.KProducer
}

func NewAssignDriverHdl(repo AssignDriverRepo, lRepo AssignDriverLocationRepo, producer kcomp.KProducer) *assignDriverHdl {
	return &assignDriverHdl{repo: repo, lRepo: lRepo, producer: producer}
}

func (h *assignDriverHdl) Response(ctx context.Context, payload *domain.TripRequestedMessage) error {
	locations, err := h.lRepo.ListLocations(ctx, locationRepo.ListLocationsParams{
		Lat:      payload.Data.PickupLocation.Latitude,
		Long:     payload.Data.PickupLocation.Longitude,
		Distance: 10000,
		Size:     1,
	})
	if err != nil {
		return err
	}

	driverID := locations[0].UserID

	tripData := tripRepo.AssignDriverParams{
		TripCode: payload.TripCode,
		DriverID: null.IntFrom(driverID),
		Status:   domain.TripStatusDriverAssigned,
	}

	estimatedTime := time.Now().Add(15 * time.Minute)

	eventData := tripRepo.CreateTripEventParams{
		TripID:    payload.Data.TripID,
		TripCode:  payload.TripCode,
		EventType: domain.TripEventTypeDriverAssigned,
		EventTime: time.Now(),
		EventData: domain.TripEventData{
			DriverAssigned: &domain.DriverAssignedData{
				DriverID:            driverID,
				EstimatedPickupTime: estimatedTime,
			},
		},
	}

	if err = h.repo.AssignDriverWithEvent(ctx, tripData, eventData); err != nil {
		return err
	}

	message := domain.TripDriverAssignedMessage{
		EventType: string(domain.TripEventTypeDriverAssigned),
		TripCode:  payload.TripCode,
		Timestamp: time.Now(),
		Data: struct {
			TripID              int64     `json:"trip_id"`
			DriverID            int64     `json:"driver_id"`
			EstimatedPickupTime time.Time `json:"estimated_pickup_time"`
		}{
			TripID:              payload.Data.TripID,
			DriverID:            driverID,
			EstimatedPickupTime: estimatedTime,
		},
	}

	_ = h.producer.Publish(ctx, string(domain.TripTopicAssignDriver), "", message)

	return nil
}
