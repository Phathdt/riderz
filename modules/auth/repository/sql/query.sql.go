// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package authRepo

import (
	"context"
)

const getUser = `-- name: GetUser :one
SELECT id, email, password, active, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (*User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
