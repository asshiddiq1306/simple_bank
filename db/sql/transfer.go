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
