package handlers

import (
	"context"
	"github.com/phathdt/service-context/core"
	"riderz/modules/location/dto"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/shared/errorx"
)

type UpdateLocationRepo interface {
	CreateLocation(ctx context.Context, arg locationRepo.CreateLocationParams) error
}

type updateLocationHdl struct {
	repo UpdateLocationRepo
}

func NewUpdateLocationHdl(repo UpdateLocationRepo) *updateLocationHdl {
	return &updateLocationHdl{repo: repo}
}

func (h *updateLocationHdl) Response(ctx context.Context, data *dto.UpdateLocationRequest) error {
	//TODO: push to queue instead of update directly
	err := h.repo.CreateLocation(ctx, locationRepo.CreateLocationParams{
		UserID: data.UserId,
		Lat:    data.Latitude,
		Long:   data.Longitude,
	})
	if err != nil {
		return core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return nil
}
