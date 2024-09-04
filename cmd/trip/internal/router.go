package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/fiberc"
	"github.com/phathdt/service-context/component/fiberc/middleware"
	"riderz/modules/trip/transport/tripfiber"
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

	//user scope
	app.Post("/trips", tripfiber.RequestTrip(sc))
	app.Get("/trips/:trip_code", tripfiber.GetTrip(sc))
	app.Get("/trips", tripfiber.ListTrip(sc))
	app.Post("/user/trips/:trip_code/cancel", tripfiber.CancelTrip(sc))

	//driver scope
	app.Get("/driver/trips/:trip_code", tripfiber.GetTrip(sc))
	app.Get("/driver/trips", tripfiber.ListTrip(sc))
	app.Post("/driver/trips/:trip_code/driver_arrived", tripfiber.DriverArrived(sc))
	app.Post("/driver/trips/:trip_code/start", tripfiber.StartTrip(sc))
	//app.Post("/driver/trips/:trip_code/end", tripfiber.EndTrip(sc))

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
