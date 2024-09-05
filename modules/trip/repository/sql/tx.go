package tripRepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) execTx(ctx context.Context, fn func(*Queries) error) error {
	db := q.db.(*pgxpool.Pool)

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)

	if err = fn(qtx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (q *Queries) execTxWithEvent(ctx context.Context, fn func(*Queries) error, eventArg CreateTripEventParams) error {
	return q.execTx(ctx, func(qtx *Queries) error {
		if err := fn(qtx); err != nil {
			return err
		}

		if _, err := qtx.CreateTripEvent(ctx, eventArg); err != nil {
			return fmt.Errorf("error creating trip event: %w", err)
		}

		return nil
	})
}

func (q *Queries) CreateTripAndEvent(ctx context.Context, tripArg CreateTripParams, eventArg CreateTripEventParams) (tripID int64, err error) {
	err = q.execTxWithEvent(ctx, func(qtx *Queries) error {
		var err error
		tripID, err = qtx.CreateTrip(ctx, tripArg)
		if err != nil {
			return fmt.Errorf("error creating trip: %w", err)
		}

		eventArg.TripID = tripID

		return nil
	}, eventArg)

	return tripID, err
}

func (q *Queries) StartTripWithEvent(ctx context.Context, arg StartTripParams, eventArg CreateTripEventParams) error {
	return q.execTxWithEvent(ctx, func(qtx *Queries) error {
		return qtx.StartTrip(ctx, arg)
	}, eventArg)
}

func (q *Queries) UpdateTripStatusWithEvent(ctx context.Context, arg UpdateTripStatusParams, eventArg CreateTripEventParams) error {
	return q.execTxWithEvent(ctx, func(qtx *Queries) error {
		return qtx.UpdateTripStatus(ctx, arg)
	}, eventArg)
}

func (q *Queries) AssignDriverWithEvent(ctx context.Context, arg AssignDriverParams, eventArg CreateTripEventParams) error {
	return q.execTxWithEvent(ctx, func(qtx *Queries) error {
		return qtx.AssignDriver(ctx, arg)
	}, eventArg)
}

func (q *Queries) CompleteTripWithEvent(ctx context.Context, arg CompleteTripParams, eventArg CreateTripEventParams) error {
	return q.execTxWithEvent(ctx, func(qtx *Queries) error {
		return qtx.CompleteTrip(ctx, arg)
	}, eventArg)
}
