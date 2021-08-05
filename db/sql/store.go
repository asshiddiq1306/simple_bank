package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxArg) (TransferTxResult, error)
}

type SQLStore struct {
	*Query
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:    db,
		Query: NewQuery(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(query *Query) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQuery(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rb error: %s, tx error : %s", rbErr, err)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxArg struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxArg) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(query *Query) error {
		var err error

		result.Transfer, err = query.CreateNewTransfer(ctx, CreateNewTransferArgs{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = query.CreateNewEntry(ctx, CreateNewEntryArgs{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = query.CreateNewEntry(ctx, CreateNewEntryArgs{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, query, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, query, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	query *Query,
	fromID int64,
	amount1 int64,
	toID int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = query.UpdateAccountBalance(ctx, UpdateAccountBalanceArgs{
		ID:     fromID,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = query.UpdateAccountBalance(ctx, UpdateAccountBalanceArgs{
		ID:     toID,
		Amount: amount2,
	})

	return

}
