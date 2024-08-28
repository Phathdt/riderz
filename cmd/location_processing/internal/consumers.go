package cmd

import (
	sctx "github.com/phathdt/service-context"
	"riderz/modules/location/transport/locationconsumer"
	"riderz/plugins/kcomp"
	"riderz/shared/common"
)

func SetupConsumer(sc sctx.ServiceContext) {
	c := sc.MustGet(common.KeyConsumer).(kcomp.KConsumer)

	c.Subscribe("demo-group123", "driver-locations", locationconsumer.ProcessUpdateLocation(sc))
}
