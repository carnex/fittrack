package store

import (
	"context"

	db "github.com/carnex/fittrack/backend/db/gen"
)

type Store interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error)
	GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error)
}
