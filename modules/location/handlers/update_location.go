package handlers

import (
	"context"
	"github.com/phathdt/service-context/core"
	"riderz/modules/location/dto"
	"riderz/plugins/kcomp"
	"riderz/shared/errorx"
)

type updateLocationHdl struct {
	producer kcomp.KProducer
}

func NewUpdateLocationHdl(producer kcomp.KProducer) *updateLocationHdl {
	return &updateLocationHdl{producer: producer}
}

func (h *updateLocationHdl) Response(ctx context.Context, data *dto.UpdateLocationRequest) error {
	if err := h.producer.Publish(ctx, "driver-locations", "", data); err != nil {
		return core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return nil
}
