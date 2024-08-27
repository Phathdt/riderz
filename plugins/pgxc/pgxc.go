package pgxc

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	sctx "github.com/phathdt/service-context"
)

type PgxComp interface {
	GetConn() *pgx.Conn
}

type pgxComp struct {
	id     string
	prefix string
	dsn    string
	logger sctx.Logger
	conn   *pgx.Conn
}

func New(id string, prefix string) *pgxComp {
	return &pgxComp{id: id, prefix: prefix}
}

func (p *pgxComp) ID() string {
	return p.id
}

func (p *pgxComp) InitFlags() {
	prefix := p.prefix
	if p.prefix != "" {
		prefix += "-"
	}

	flag.StringVar(
		&p.dsn,
		fmt.Sprintf("%sdb-dsn", prefix),
		"",
		"Database dsn",
	)
}

func (p *pgxComp) Activate(_ sctx.ServiceContext) error {
	p.logger = sctx.GlobalLogger().GetLogger(p.id)

	p.logger.Info("Connecting to database...")

	var err error

	conn, err := pgx.Connect(context.Background(), p.dsn)
	if err != nil {
		p.logger.Error("Cannot connect to database", err.Error())
		return err
	}

	p.conn = conn

	return nil
}

func (p *pgxComp) Stop() error {
	return p.conn.Close(context.Background())
}

func (p *pgxComp) GetConn() *pgx.Conn {
	return p.conn
}
