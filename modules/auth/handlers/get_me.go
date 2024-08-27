package handlers

import (
	"context"
	authRepo "riderz/modules/auth/repository/sql"
	"riderz/shared/errorx"

	"github.com/phathdt/service-context/core"
)

type GetMeStorage interface {
	GetUserById(ctx context.Context, id int64) (*authRepo.User, error)
}

type getMeHdl struct {
	store GetMeStorage
}

func NewGetMeHdl(store GetMeStorage) *getMeHdl {
	return &getMeHdl{store}
}

func (h *getMeHdl) Response(ctx context.Context, userId int64) (*authRepo.User, error) {
	user, err := h.store.GetUserById(ctx, userId)
	if err != nil {
		return nil, core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return user, nil
}
