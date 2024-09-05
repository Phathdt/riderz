package tripRepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) CreateTripAndEvent(ctx context.Context, tripArg CreateTripParams, eventArg CreateTripEventParams) (tripID int64, err error) {
	db := q.db.(*pgxpool.Pool)

	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)

	tripID, err = qtx.CreateTrip(ctx, tripArg)
	if err != nil {
		return 0, fmt.Errorf("error creating trip: %w", err)
	}

	eventArg.TripID = tripID

	_, err = qtx.CreateTripEvent(ctx, eventArg)
	if err != nil {
		return 0, fmt.Errorf("error creating trip event: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return tripID, nil
}
