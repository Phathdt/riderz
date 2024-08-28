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
