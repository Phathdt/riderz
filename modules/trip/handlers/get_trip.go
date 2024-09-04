package handlers

import (
	"golang.org/x/net/context"
	tripRepo "riderz/modules/trip/repository/sql"
)

type GetTripRepo interface {
	GetTripByUser(ctx context.Context, arg tripRepo.GetTripByUserParams) (*tripRepo.Trip, error)
}

type getTripHdl struct {
	repo GetTripRepo
}

func NewGetTripHdl(repo GetTripRepo) *getTripHdl {
	return &getTripHdl{repo: repo}
}

func (h *getTripHdl) Response(ctx context.Context, userID int64, tripCode string) (*tripRepo.Trip, error) {
	trip, err := h.repo.GetTripByUser(ctx, tripRepo.GetTripByUserParams{
		TripCode: tripCode,
		UserID:   userID,
	})

	if err != nil {
		return nil, err
	}

	return trip, nil
}
