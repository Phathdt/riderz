package handlers

import (
	"golang.org/x/net/context"
	tripRepo "riderz/modules/trip/repository/sql"
)

type GetTripRepo interface {
	GetTripByUser(ctx context.Context, arg tripRepo.GetTripByUserParams) (*tripRepo.Trip, error)
	ListTripEvents(ctx context.Context, tripCode string) ([]*tripRepo.TripEvent, error)
}

type getTripHdl struct {
	repo GetTripRepo
}

func NewGetTripHdl(repo GetTripRepo) *getTripHdl {
	return &getTripHdl{repo: repo}
}

func (h *getTripHdl) Response(ctx context.Context, userID int64, tripCode string) (*tripRepo.TripWithEvents, error) {
	trip, err := h.repo.GetTripByUser(ctx, tripRepo.GetTripByUserParams{
		TripCode: tripCode,
		UserID:   userID,
	})

	if err != nil {
		return nil, err
	}

	events, err := h.repo.ListTripEvents(ctx, tripCode)
	if err != nil {
		return nil, err
	}

	return &tripRepo.TripWithEvents{
		Trip:   trip,
		Events: events,
	}, nil
}
