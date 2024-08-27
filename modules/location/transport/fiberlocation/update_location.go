package fiberlocation

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/location/dto"
	"riderz/modules/location/handlers"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/plugins/pgxc"
	"riderz/plugins/validation"
	"riderz/shared/common"
)

func UpdateLocation(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int64)
		var p dto.UpdateLocationRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		p.UserId = userId

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()

		sqlStorage := locationRepo.New(conn)
		hdl := handlers.NewUpdateLocationHdl(sqlStorage)

		if err := hdl.Response(ctx.Context(), &p); err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse("ok"))
	}
}
