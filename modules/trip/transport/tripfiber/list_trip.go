package tripfiber

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/trip/handlers"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
)

func ListTrip(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int64)

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := tripRepo.New(conn)
		hdl := handlers.NewListTripHdl(repo)

		trip, err := hdl.Response(ctx.Context(), userId)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(trip))
	}
}
