package cmd

import (
	sctx "github.com/phathdt/service-context"
	"riderz/modules/trip/domain"
	"riderz/modules/trip/transport/tripconsumer"
	"riderz/plugins/kcomp"
	"riderz/shared/common"
)

func SetupConsumer(sc sctx.ServiceContext) {
	c := sc.MustGet(common.KeyConsumer).(kcomp.KConsumer)

	c.Subscribe("trip-created", string(domain.TripTopicRequested), tripconsumer.AssignDriver(sc))
}
