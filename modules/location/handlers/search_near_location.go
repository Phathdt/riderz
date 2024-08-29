package handlers

import (
	"context"
	"github.com/phathdt/service-context/core"
	"riderz/modules/location/domain"
	"riderz/modules/location/dto"
	locationRepo "riderz/modules/location/repository/sql"
)

type SearchNearLocationRepo interface {
	ListLocations(ctx context.Context, arg locationRepo.ListLocationsParams) ([]*locationRepo.ListLocationsRow, error)
}

type searchNearLocationHdl struct {
	repo SearchNearLocationRepo
}

func NewSearchNearLocationHdl(repo SearchNearLocationRepo) *searchNearLocationHdl {
	return &searchNearLocationHdl{repo: repo}
}

func (h *searchNearLocationHdl) Response(ctx context.Context, data *dto.SearchNearLocation) ([]*locationRepo.ListLocationsRow, error) {
	params := locationRepo.ListLocationsParams{
		Lat:      data.Latitude,
		Long:     data.Longitude,
		Size:     10,
		Distance: 10000,
	}

	locations, err := h.repo.ListLocations(ctx, params)
	if err != nil {
		return nil, core.ErrNotFound.
			WithError(domain.ErrCannotSearchNearLocation.Error()).
			WithDebug(err.Error())
	}

	return locations, nil
}
