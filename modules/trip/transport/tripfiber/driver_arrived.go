package tripfiber

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/validation"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/trip/dto"
	"riderz/modules/trip/handlers"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
)

func DriverArrived(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int64)
		tripCode := ctx.Params("trip_code")

		var p dto.DriverArrivedData

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		p.DriverID = userId

		producer := sc.MustGet(common.KeyProducer).(kcomp.KProducer)
		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := tripRepo.New(conn)
		hdl := handlers.NewDriverArrivedHdl(repo, producer)

		err := hdl.Response(ctx.Context(), tripCode, &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse("ok"))
	}
}
