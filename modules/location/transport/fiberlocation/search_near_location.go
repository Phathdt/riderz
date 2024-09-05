package fiberlocation

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/validation"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/location/dto"
	"riderz/modules/location/handlers"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
)

func SearchNearLocation(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p dto.SearchNearLocation

		if err := ctx.QueryParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := locationRepo.New(conn)
		hdl := handlers.NewSearchNearLocationHdl(repo)

		response, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(response))
	}
}
