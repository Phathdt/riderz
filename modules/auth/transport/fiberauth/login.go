package fiberauth

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/redisc"
	"github.com/phathdt/service-context/component/validation"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/auth/dto"
	"riderz/modules/auth/handlers"
	"riderz/modules/auth/repository/sessionRepo"
	authRepo "riderz/modules/auth/repository/sql"
	"riderz/plugins/pgxc"
	"riderz/plugins/tokenprovider"
	"riderz/shared/common"
)

func Login(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.LoginRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)
		rdClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()

		sqlStorage := authRepo.New(conn)
		sessionStore := sessionRepo.NewSessionStore(rdClient)
		hdl := handlers.NewLoginHandler(sqlStorage, sessionStore, tokenProvider)

		token, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(token))
	}
}
