package handlers

import (
	"context"
	"github.com/phathdt/service-context/core"
	"riderz/modules/location/domain"
	"riderz/modules/location/dto"
	locationRepo "riderz/modules/location/repository/sql"
)

type ProcessUpdateLocationRepo interface {
	CreateLocation(ctx context.Context, arg locationRepo.CreateLocationParams) error
}

type processUpdateLocationHdl struct {
	repo ProcessUpdateLocationRepo
}

func NewProcessUpdateLocationHdl(repo ProcessUpdateLocationRepo) *processUpdateLocationHdl {
	return &processUpdateLocationHdl{repo: repo}
}

func (h *processUpdateLocationHdl) Response(ctx context.Context, payload *dto.UpdateLocationRequest) error {
	if err := h.repo.CreateLocation(ctx, locationRepo.CreateLocationParams{
		UserID: payload.UserId,
		Lat:    payload.Latitude,
		Long:   payload.Longitude,
	}); err != nil {
		return core.ErrNotFound.
			WithError(domain.ErrCannotUpdateLocation.Error()).
			WithDebug(err.Error())
	}

	return nil
}
