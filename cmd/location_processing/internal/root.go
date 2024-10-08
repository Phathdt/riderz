package cmd

import (
	"fmt"
	"os"
	"os/signal"
	kconsumercomp "riderz/plugins/kcomp/consumercomp"
	"riderz/plugins/pgxc"
	"riderz/shared/common"
	"syscall"
	"time"

	sctx "github.com/phathdt/service-context"
	"github.com/spf13/cobra"
)

const (
	serviceName = "location_processing"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(pgxc.New(common.KeyPgx, "")),
		sctx.WithComponent(kconsumercomp.New(common.KeyConsumer)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 1)

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

		SetupConsumer(sc)
		// gracefully shutdown
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		_ = sc.Stop()
		logger.Info("Server exited")
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
