// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package locationRepo

import (
	"context"

	"riderz/shared/common"
)

const createLocation = `-- name: CreateLocation :exec
INSERT INTO locations (user_id, geometry)
VALUES ($1, ST_SetSRID(ST_MakePoint($2::decimal, $3::decimal), 4326))
ON CONFLICT (user_id)
DO UPDATE SET
    geometry = EXCLUDED.geometry,
    updated_at = now()
`

type CreateLocationParams struct {
	UserID int64   `db:"user_id" json:"user_id"`
	Lat    float64 `db:"lat" json:"lat"`
	Long   float64 `db:"long" json:"long"`
}

func (q *Queries) CreateLocation(ctx context.Context, arg CreateLocationParams) error {
	_, err := q.db.Exec(ctx, createLocation, arg.UserID, arg.Lat, arg.Long)
	return err
}

const getLocation = `-- name: GetLocation :one
SELECT id, user_id, geometry, created_at, updated_at FROM locations
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetLocation(ctx context.Context, userID int64) (*Location, error) {
	row := q.db.QueryRow(ctx, getLocation, userID)
	var i Location
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Geometry,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listLocations = `-- name: ListLocations :many
SELECT id,
  user_id,
  geometry,
  ST_Distance(geometry::geography, ST_SetSRID(ST_MakePoint($1::decimal, $2::decimal), 4326)::geography)::decimal AS distance
FROM locations
WHERE ST_DWithin(geometry::geography, ST_SetSRID(ST_MakePoint($1::decimal, $2::decimal), 4326)::geography, $3)
ORDER BY geometry <-> ST_SetSRID(ST_MakePoint($1::decimal, $2::decimal), 4326)
LIMIT $4
`

type ListLocationsParams struct {
	Lat      float64     `db:"lat" json:"lat"`
	Long     float64     `db:"long" json:"long"`
	Distance interface{} `db:"distance" json:"distance"`
	Size     int32       `db:"size" json:"size"`
}

type ListLocationsRow struct {
	ID       int64         `db:"id" json:"id"`
	UserID   int64         `db:"user_id" json:"user_id"`
	Geometry common.PointS `db:"geometry" json:"geometry"`
	Distance float64       `db:"distance" json:"distance"`
}

func (q *Queries) ListLocations(ctx context.Context, arg ListLocationsParams) ([]*ListLocationsRow, error) {
	rows, err := q.db.Query(ctx, listLocations,
		arg.Lat,
		arg.Long,
		arg.Distance,
		arg.Size,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListLocationsRow
	for rows.Next() {
		var i ListLocationsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Geometry,
			&i.Distance,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
