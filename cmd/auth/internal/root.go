package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"riderz/plugins/pgxc"
	"riderz/plugins/tokenprovider/jwt"
	"riderz/shared/common"
	"syscall"
	"time"

	"github.com/phathdt/service-context/component/fiberc"
	"github.com/phathdt/service-context/component/redisc"

	sctx "github.com/phathdt/service-context"

	"github.com/spf13/cobra"
)

const (
	serviceName = "auth"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(fiberc.New(common.KeyCompFiber)),
		sctx.WithComponent(pgxc.New(common.KeyPgx, "")),
		sctx.WithComponent(jwt.New(common.KeyJwt)),
		sctx.WithComponent(redisc.New(common.KeyCompRedis)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 1)

		NewRouter(sc)

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

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
