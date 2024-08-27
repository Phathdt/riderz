package handlers

import (
	"context"
	"errors"
	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
	"riderz/modules/auth/dto"
	authRepo "riderz/modules/auth/repository/sql"
	"riderz/plugins/tokenprovider"
	"riderz/shared/common"
	"riderz/shared/errorx"
)

type SignupUser interface {
	GetUserByEmail(ctx context.Context, email string) (*authRepo.User, error)
	CreateUser(ctx context.Context, arg authRepo.CreateUserParams) (*authRepo.User, error)
}

type signUpSessionStorage interface {
	SetUserToken(ctx context.Context, userId int64, token, subToken string, expiredTime int) error
}

type signupHdl struct {
	store         SignupUser
	sessionStore  signUpSessionStorage
	tokenProvider tokenprovider.Provider
}

func NewSignupHdl(store SignupUser, sStore signUpSessionStorage, tokenProvider tokenprovider.Provider) *signupHdl {
	return &signupHdl{store: store, sessionStore: sStore, tokenProvider: tokenProvider}
}

func (h *signupHdl) Response(ctx context.Context, params *dto.SignupRequest) (tokenprovider.Token, error) {
	user, err := h.store.GetUserByEmail(ctx, params.Email)
	if err == nil && user != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrUserAlreadyExists.Error()).
			WithDebug(errors.New("user already exist").Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	data := authRepo.CreateUserParams{Email: params.Email, Password: string(hashedPassword)}

	newUser, err := h.store.CreateUser(ctx, data)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   newUser.ID,
		Email:    data.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := h.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	if err = h.sessionStore.SetUserToken(ctx, newUser.ID, accessToken.GetToken(), subToken, expiredTime); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
