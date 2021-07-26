package db

import "context"

const insertNewTransferQuery = `-- name: CreateNewTransfer :one
INSERT INTO transfers (
	from_account_id, to_account_id, amount
) VALUES (
	$1, $2, $3
) RETURNING *
`

type CreateNewTransferArgs struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (query *Query) CreateNewTransfer(ctx context.Context, arg CreateNewTransferArgs) (Transfer, error) {
	row := query.db.QueryRowContext(ctx, insertNewTransferQuery, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const selectTransferByIDQuery = `-- name: GetTransferByID :one
SELECT * FROM transfers WHERE id = $1 LIMIT 1
`

func (query *Query) GetTransferByID(ctx context.Context, id int64) (Transfer, error) {
	row := query.db.QueryRowContext(ctx, selectTransferByIDQuery, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const selectTransferListQuery = `-- name: GetTransfersList :many
SELECT * FROM transfers LIMIT $1 OFFSET $2
`

type GetTransfersListArgs struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (query *Query) GetTransfersList(ctx context.Context, arg GetTransfersListArgs) ([]Transfer, error) {
	rows, err := query.db.QueryContext(ctx, selectTransferListQuery, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
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

const updateTransferByIDQuery = `-- name: UpdateTransferByID :one
UPDATE transfers SET amount = $2 WHERE id = $1
RETURNING *
`

type UpdateTransferByIDArgs struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (query *Query) UpdateTransferByID(ctx context.Context, arg UpdateTransferByIDArgs) (Transfer, error) {
	row := query.db.QueryRowContext(ctx, updateTransferByIDQuery, arg.ID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransferByIDQuery = `-- name: DeleteTransferByID :exec
DELETE FROM transfers WHERE id = $1
`

func (query *Query) DeleteTransferByID(ctx context.Context, id int64) error {
	_, err := query.db.ExecContext(ctx, deleteTransferByIDQuery, id)
	return err
}
