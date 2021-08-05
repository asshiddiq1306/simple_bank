package db

import "context"

const insertNewUserQuery = `-- name: CreateNewUser :one
INSERT INTO users(
	username, hashed_password, full_name, email
) VALUES (
	$1, $2, $3, $4
) RETURNING *
`

type CreateNewUserArgs struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (query *Query) CreateNewUser(ctx context.Context, arg CreateNewUserArgs) (User, error) {
	row := query.db.QueryRowContext(ctx, insertNewUserQuery, arg.Username, arg.HashedPassword, arg.FullName, arg.Email)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const selectUserByUsernameQuery = `-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1
`

func (query *Query) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := query.db.QueryRowContext(ctx, selectUserByUsernameQuery, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const selectListUserQuery = `-- name: GetListUser :many
SELECT * FROM users LIMIT $1 OFFSET $2
`

type GetListUserArgs struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (query *Query) GetListUser(ctx context.Context, arg GetListUserArgs) ([]User, error) {
	rows, err := query.db.QueryContext(ctx, selectListUserQuery, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	items := []User{}

	for rows.Next() {
		var i User

		if err := rows.Scan(
			&i.Username,
			&i.HashedPassword,
			&i.FullName,
			&i.Email,
			&i.PasswordChangedAt,
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

const deleteUserByIDQuery = `-- name: DeleteUserByID :exec
DELETE FROM users WHERE username = $1
`

func (query *Query) DeleteUserByID(ctx context.Context, username string) error {
	_, err := query.db.ExecContext(ctx, deleteUserByIDQuery, username)
	return err
}
