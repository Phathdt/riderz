package fiberauth

import (
	"net/http"
	"riderz/modules/auth/handlers"
	"riderz/modules/auth/models"
	"riderz/modules/auth/repository/sessionRepo"
	"riderz/modules/auth/storage"
	"riderz/shared/common"
	"riderz/shared/plugins/tokenprovider"
	"riderz/shared/plugins/validation"

	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/component/redisc"
	"github.com/phathdt/service-context/core"
)

func SignUp(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p models.SignupRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)
		rdClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()

		sqlStorage := storage.NewSqlStorage(db)
		sessionStore := sessionRepo.NewSessionStore(rdClient)
		hdl := handlers.NewSignupHdl(sqlStorage, sessionStore, tokenProvider)

		token, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(token))
	}
}
