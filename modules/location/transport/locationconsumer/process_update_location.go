package locationconsumer

import (
	"context"
	"encoding/json"
	sctx "github.com/phathdt/service-context"
	"riderz/modules/location/dto"
	"riderz/modules/location/handlers"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/plugins/kcomp"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
)

func ProcessUpdateLocation(sc sctx.ServiceContext) kcomp.HandlerFunc {
	return func(msg *kcomp.Message) error {
		ctx := context.Background()
		defer ctx.Done()

		var payload dto.UpdateLocationRequest

		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			panic(err)
		}

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := locationRepo.New(conn)
		hdl := handlers.NewProcessUpdateLocationHdl(repo)

		return hdl.Response(ctx, &payload)
	}
}

func ProcessUpdateMultiLocation(sc sctx.ServiceContext) kcomp.HandlerMultiFunc {
	return func(msgs []*kcomp.Message) error {
		mapUserRequest := make(map[int64]*dto.UpdateLocationRequest)

		for _, msg := range msgs {
			var payload dto.UpdateLocationRequest

			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				continue
			}

			mapUserRequest[payload.UserId] = &payload
		}

		logger := sc.Logger("process-update-multi-location")

		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()
		repo := locationRepo.New(conn)
		hdl := handlers.NewProcessUpdateLocationHdl(repo)

		for _, locationRequest := range mapUserRequest {
			func() {
				ctx := context.Background()
				defer ctx.Done()
				if err := hdl.Response(ctx, locationRequest); err != nil {
					logger.Fatalf("failed to update userId %d: %v", locationRequest.UserId, err)
				}
			}()
		}

		return nil
	}
}
