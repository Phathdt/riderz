package fiberauth

import (
	"net/http"
	"riderz/modules/auth/handlers"
	authRepo "riderz/modules/auth/repository/sql"
	"riderz/plugins/pgxc"
	"riderz/shared/common"

	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
)

func GetMe(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int64)
		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		sqlStorage := authRepo.New(conn)
		hdl := handlers.NewGetMeHdl(sqlStorage)

		user, err := hdl.Response(ctx.Context(), userId)
		if err != nil {
			panic(err)
		}

		user.Mask()

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(user))
	}
}
