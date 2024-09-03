package handlers

import (
	"context"
	tripRepo "riderz/modules/trip/repository/sql"
)

type ListTripRepo interface {
	ListTrips(ctx context.Context, userID int64) ([]*tripRepo.Trip, error)
}

type listTripHdl struct {
	repo ListTripRepo
}

func NewListTripHdl(repo ListTripRepo) *listTripHdl {
	return &listTripHdl{repo: repo}
}

func (h *listTripHdl) Response(ctx context.Context, userID int64) ([]*tripRepo.Trip, error) {
	trips, err := h.repo.ListTrips(ctx, userID)
	if err != nil {
		return nil, err
	}

	return trips, nil
}
