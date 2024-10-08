package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/fiberc"
	"github.com/phathdt/service-context/component/fiberc/middleware"
	"riderz/modules/location/transport/fiberlocation"
	"riderz/shared/common"
	middleware2 "riderz/shared/middleware"
)

func NewRouter(sc sctx.ServiceContext) {
	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})
	app.Use(flogger.New(flogger.Config{
		Format: `{"ip":${ip}, "timestamp":"${time}", "status":${status},  "method":"${method}", "path":"${path}"}` + "\n",
	}))
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.Recover(sc))

	app.Get("/", ping())

	app.Use(middleware2.RequiredAuth(sc))

	app.Get("/locations/search", fiberlocation.SearchNearLocation(sc))

	fiberComp := sc.MustGet(common.KeyCompFiber).(fiberc.FiberComponent)
	fiberComp.SetApp(app)
}

func ping() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(&fiber.Map{
			"msg": "pong",
		})
	}
}
