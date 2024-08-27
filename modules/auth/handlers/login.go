package handlers

import (
	"context"
	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
	"riderz/modules/auth/dto"
	authRepo "riderz/modules/auth/repository/sql"
	"riderz/plugins/tokenprovider"
	"riderz/shared/common"
	"riderz/shared/errorx"
)

type LoginDbStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*authRepo.User, error)
}

type LoginSessionStorage interface {
	SetUserToken(ctx context.Context, userId int64, token, subToken string, expiredTime int) error
}

type loginHandler struct {
	store         LoginDbStorage
	sStore        LoginSessionStorage
	tokenProvider tokenprovider.Provider
}

func NewLoginHandler(store LoginDbStorage, sStore LoginSessionStorage, tokenProvider tokenprovider.Provider) *loginHandler {
	return &loginHandler{store: store, sStore: sStore, tokenProvider: tokenProvider}
}

func (h *loginHandler) Response(ctx context.Context, params *dto.LoginRequest) (tokenprovider.Token, error) {
	user, err := h.store.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	userPass := []byte(params.Password)
	dbPass := []byte(user.Password)

	if err = bcrypt.CompareHashAndPassword(dbPass, userPass); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrPasswordNotMatch.Error()).
			WithDebug(err.Error())
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   user.ID,
		Email:    user.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := h.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	if err = h.sStore.SetUserToken(ctx, user.ID, accessToken.GetToken(), subToken, expiredTime); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
