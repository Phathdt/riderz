package tripconsumer

import (
	"context"
	"encoding/json"
	sctx "github.com/phathdt/service-context"
	locationRepo "riderz/modules/location/repository/sql"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/handlers"
	tripRepo "riderz/modules/trip/repository/sql"
	"riderz/plugins/kcomp"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
)

func AssignDriver(sc sctx.ServiceContext) kcomp.HandlerFunc {
	return func(msg *kcomp.Message) error {
		ctx := context.Background()
		defer ctx.Done()

		var payload domain.TripRequestedMessage

		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		producer := sc.MustGet(common.KeyProducer).(kcomp.KProducer)
		conn := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()

		lRepo := locationRepo.New(conn)
		repo := tripRepo.New(conn)
		hdl := handlers.NewAssignDriverHdl(repo, lRepo, producer)

		err = hdl.Response(ctx, &payload)
		if err != nil {
			return err
		}

		return nil
	}
}
