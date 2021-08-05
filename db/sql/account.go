package db

import "context"

const insertNewAccount = `-- name: CreateNewAccount :one
INSERT INTO accounts (
	owner, balance, currency
) VALUES (
	$1, $2, $3
) RETURNING *
`

type CreateNewAccountArgs struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (query *Query) CreateNewAccount(ctx context.Context, arg CreateNewAccountArgs) (Account, error) {
	row := query.db.QueryRowContext(ctx, insertNewAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

const selectAccountByID = `-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (query *Query) GetAccountByID(ctx context.Context, id int64) (Account, error) {
	row := query.db.QueryRowContext(ctx, selectAccountByID, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const selectAllAccounts = `-- name: GetListAccounts :many
SELECT * FROM accounts WHERE owner = $1 LIMIT $2 OFFSET $3
`

type GetListAccountsArgs struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (query *Query) GetListAccounts(ctx context.Context, arg GetListAccountsArgs) ([]Account, error) {
	rows, err := query.db.QueryContext(ctx, selectAllAccounts, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []Account{}
	for rows.Next() {
		var i Account
		if rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
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

const updateAccountByIDQuery = `-- name: UpdateAccountByID :one
UPDATE accounts SET balance = $2 WHERE id = $1
RETURNING *
`

type UpdateAccountByIDArgs struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (query *Query) UpdateAccountByID(ctx context.Context, arg UpdateAccountByIDArgs) (Account, error) {
	row := query.db.QueryRowContext(ctx, updateAccountByIDQuery, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts SET balance = balance + $2 WHERE id = $1
RETURNING *
`

type UpdateAccountBalanceArgs struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (query *Query) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceArgs) (Account, error) {
	row := query.db.QueryRowContext(ctx, updateAccountBalance, arg.ID, arg.Amount)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccountBuIDQuery = `-- name: DeleteAccountBuID :exec
DELETE FROM accounts WHERE id = $1
`

func (query *Query) DeleteAccountBuID(ctx context.Context, id int64) error {
	_, err := query.db.ExecContext(ctx, deleteAccountBuIDQuery, id)
	return err
}
