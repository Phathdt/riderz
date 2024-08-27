package pgxc

import (
	"context"
	"github.com/jackc/pgx/v5/tracelog"
	sctx "github.com/phathdt/service-context"
)

type PgxLogAdapter struct {
	logger sctx.Logger
}

func (l *PgxLogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case tracelog.LogLevelTrace:
		l.logger.Debugf("%s %v", msg, data)
	case tracelog.LogLevelDebug:
		l.logger.Debugf("%s %v", msg, data)
	case tracelog.LogLevelInfo:
		l.logger.Infof("%s %v", msg, data)
	case tracelog.LogLevelWarn:
		l.logger.Warnf("%s %v", msg, data)
	case tracelog.LogLevelError:
		l.logger.Errorf("%s %v", msg, data)
	default:
		l.logger.Infof("%s %v", msg, data)
	}
}
