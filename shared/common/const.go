package common

import (
	"github.com/namsral/flag"
)

var (
	ShowLog = false
)

func init() {
	flag.BoolVar(&ShowLog, "show-log", false, "show log")
	flag.Parse()
}

const (
	KeyCompFiber = "fiberlocation"
	KeyCompRedis = "redis"
	KeyJwt       = "jwt"
	KeyPgx       = "pgx"
	KeyAuthen    = "authen"
)
