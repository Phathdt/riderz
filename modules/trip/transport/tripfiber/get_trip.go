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

func GetTrip(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int64)
		tripCode := ctx.Params("trip_code")

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := tripRepo.New(conn)
		hdl := handlers.NewGetTripHdl(repo)

		trip, err := hdl.Response(ctx.Context(), userId, tripCode)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(trip))
	}
}
