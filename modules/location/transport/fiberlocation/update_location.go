package fiberlocation

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"net/http"
	"riderz/modules/location/dto"
	"riderz/modules/location/handlers"
	"riderz/plugins/kcomp"
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

		producer := sc.MustGet(common.KeyProducer).(kcomp.KProducer)
		hdl := handlers.NewUpdateLocationHdl(producer)

		if err := hdl.Response(ctx.Context(), &p); err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse("ok"))
	}
}
