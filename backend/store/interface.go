package store

import (
	"context"

	db "github.com/carnex/fittrack/backend/db/gen"
)

type Store interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error)
}
