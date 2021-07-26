package db

import (
	"context"
)

const insertNewEntryQuery = `-- name: CreateNewEntry :one
INSERT INTO entries (
	account_id, amount
) VALUES (
	$1, $2
) RETURNING *
`

type CreateNewEntryArgs struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (query *Query) CreateNewEntry(ctx context.Context, arg CreateNewEntryArgs) (Entry, error) {
	row := query.db.QueryRowContext(ctx, insertNewEntryQuery, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const selectEntryByIDQuery = `-- name: GetEntryByID :one
SELECT * FROM entries WHERE id = $1 LIMIT 1
`

func (query *Query) GetEntryByID(ctx context.Context, id int64) (Entry, error) {
	row := query.db.QueryRowContext(ctx, selectEntryByIDQuery, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const selectEntryListQuery = `-- name: GetEntryList :many
SELECT * FROM entries LIMIT $1 OFFSET $2
`

type GetEntryListArgs struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (query *Query) GetEntryList(ctx context.Context, arg GetEntryListArgs) ([]Entry, error) {
	rows, err := query.db.QueryContext(ctx, selectEntryListQuery, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEntryByIDQuery = `-- name: UpdateEntryByID :one
UPDATE entries SET amount = $2 WHERE id = $1
RETURNING *
`

type UpdateEntryByIDArgs struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (query *Query) UpdateEntryByID(ctx context.Context, arg UpdateEntryByIDArgs) (Entry, error) {
	row := query.db.QueryRowContext(ctx, updateEntryByIDQuery, arg.ID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntryByIDQuery = `-- name: DeleteEntryByID :exec
DELETE FROM entries WHERE id = $1
`

func (query *Query) DeleteEntryByID(ctx context.Context, id int64) error {
	_, err := query.db.ExecContext(ctx, deleteEntryByIDQuery, id)
	return err
}
