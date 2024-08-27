package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"riderz/plugins/authcomp"
	"riderz/shared/common"
	"strings"

	"github.com/pkg/errors"
)

func ExtractTokenFromHeaderString(headers []string) (string, error) {
	if len(headers) == 0 {
		return "", errors.New("missing token")
	}
	//"Authorization" : "Bearer {token}"

	parts := strings.Split(headers[0], " ")

	if len(parts) == 0 {
		return "", errors.New("missing token")
	}

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("wrong authen header")
	}

	return parts[1], nil
}

func RequiredAuth(sc sctx.ServiceContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		token, err := ExtractTokenFromHeaderString(headers["Authorization"])
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		comp := sc.MustGet(common.KeyAuthen).(authcomp.AuthenComp)
		if err = comp.ValidateToken(c.Context(), token); err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			panic(fmt.Errorf("invalid token format"))
		}

		rawPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			panic(fmt.Errorf("error decoding payload: %v", err))
		}

		var data struct {
			Payload common.TokenPayload `json:"payload"`
		}

		err = json.Unmarshal(rawPayload, &data)
		if err != nil {
			panic(fmt.Errorf("error unmarshalling payload: %v", err))
		}

		c.Context().SetUserValue("userId", data.Payload.GetUserId())

		return c.Next()
	}
}
