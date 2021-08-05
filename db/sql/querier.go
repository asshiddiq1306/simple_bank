package db

import "context"

type Querier interface {
	CreateNewAccount(ctx context.Context, arg CreateNewAccountArgs) (Account, error)
	GetAccountByID(ctx context.Context, id int64) (Account, error)
	GetListAccounts(ctx context.Context, arg GetListAccountsArgs) ([]Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceArgs) (Account, error)
	DeleteAccountBuID(ctx context.Context, id int64) error
	CreateNewUser(ctx context.Context, arg CreateNewUserArgs) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetListUser(ctx context.Context, arg GetListUserArgs) ([]User, error)
	GetTransferByID(ctx context.Context, id int64) (Transfer, error)
	GetEntryByID(ctx context.Context, id int64) (Entry, error)
}
