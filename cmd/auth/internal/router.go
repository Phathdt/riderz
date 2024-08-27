package cmd

import (
	"github.com/phathdt/service-context/component/redisc"
	"github.com/phathdt/service-context/core"
	"riderz/modules/auth/repository/sessionRepo"
	"riderz/modules/auth/transport/fiberauth"
	"riderz/plugins/tokenprovider"
	"riderz/shared/common"
	"strings"

	middleware2 "riderz/shared/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/fiberc"
	"github.com/phathdt/service-context/component/fiberc/middleware"
)

func NewRouter(sc sctx.ServiceContext) {
	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})
	app.Use(flogger.New(flogger.Config{
		Format: `{"ip":${ip}, "timestamp":"${time}", "status":${status}, "method":"${method}", "path":"${path}"}` + "\n",
	}))
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.Recover(sc))

	app.Get("/", ping())
	app.Post("/auth/signup", fiberauth.SignUp(sc))
	app.Post("/auth/login", fiberauth.Login(sc))

	app.Use(auth(sc))

	app.Get("/auth/me", fiberauth.GetMe(sc))
	app.Get("/auth/valid", fiberauth.CheckValid(sc))

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

func auth(sc sctx.ServiceContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		token, err := middleware2.ExtractTokenFromHeaderString(headers["Authorization"])

		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}
		rdClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()
		sessionStore := sessionRepo.NewSessionStore(rdClient)

		signature, err := sessionStore.GetUserToken(c.Context(), payload.GetUserId(), payload.GetSubToken())
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		if signature != strings.Split(token, ".")[2] {
			panic(core.ErrUnauthorized.WithError("signature not matched"))
		}

		c.Context().SetUserValue("userId", payload.GetUserId())
		return c.Next()
	}
}
