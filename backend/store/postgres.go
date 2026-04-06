package store

import (
	"context"

	db "github.com/carnex/fittrack/backend/db/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	queries *db.Queries
}

func Connect(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil

}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		queries: db.New(pool),
	}
}

func (s *PostgresStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
	return s.queries.CreateUser(ctx, arg)
}

func (s *PostgresStore) GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	return s.queries.GetUserByUsername(ctx, username)
}

func (s *PostgresStore) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	return s.queries.GetUserByEmail(ctx, email)
}
