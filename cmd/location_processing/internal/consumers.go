package cmd

import (
	"fmt"
	sctx "github.com/phathdt/service-context"
	"riderz/plugins/kcomp"
	"riderz/shared/common"
)

func SetupConsumer(sc sctx.ServiceContext) {
	c := sc.MustGet(common.KeyConsumer).(kcomp.KConsumer)

	c.Subscribe("demo-group", "driver-locations", func(msg *kcomp.Message) error {
		fmt.Printf("Received message: Key=%s, Payload=%s\n", string(msg.Key), string(msg.Payload))
		return nil
	})
}
