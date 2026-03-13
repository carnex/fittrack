package service

import (
	"context"
	"errors"
	"testing"

	db "github.com/carnex/fittrack/backend/db/gen"
	"github.com/jackc/pgx/v5"
)

type MockStore struct {
	getUserByUsernameFunc func(ctx context.Context, username string) (db.GetUserByUsernameRow, error)
	getUserByEmailFunc    func(ctx context.Context, email string) (db.GetUserByEmailRow, error)
	createUserFunc        func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
}

func (m *MockStore) GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	return m.getUserByUsernameFunc(ctx, username)
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	return m.getUserByEmailFunc(ctx, email)
}

func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
	return m.createUserFunc(ctx, arg)
}

func availableMock() *MockStore {
	return &MockStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
			return db.GetUserByUsernameRow{}, pgx.ErrNoRows // not found = available
		},
		getUserByEmailFunc: func(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
			return db.GetUserByEmailRow{}, pgx.ErrNoRows // not found = available
		},
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, nil
		},
	}
}

func takenUsernameMock() *MockStore {
	mock := availableMock()
	mock.getUserByUsernameFunc = func(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
		return db.GetUserByUsernameRow{Username: username}, nil // found = taken
	}
	return mock
}

func takenEmailMock() *MockStore {
	mock := availableMock()
	mock.getUserByEmailFunc = func(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
		return db.GetUserByEmailRow{
			Email: "test@test.com",
		}, nil
	}
	return mock
}

func TestValidateRegisterInput(t *testing.T) {

	tests := []struct {
		name    string
		input   RegisterInput
		mock    *MockStore
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid input with email reset",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				ResetMethod:     true,
			},
			mock:    availableMock(),
			wantErr: false,
		},
		{
			name: "valid input with security question",
			input: RegisterInput{
				Username:         "testuser",
				Email:            "test@test.com",
				Password:         "password123",
				ConfirmPassword:  "password123",
				ResetMethod:      false,
				SecurityQuestion: "What was your first bike?",
				SecurityAnswer:   "Trek",
			},
			mock:    availableMock(),
			wantErr: false,
		},
		{
			name: "passwords do not match",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "password123",
				ConfirmPassword: "different",
				ResetMethod:     true,
			},
			mock:    availableMock(),
			wantErr: true,
			errMsg:  "passwords do not match",
		},
		{
			name: "password too short",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "short",
				ConfirmPassword: "short",
				ResetMethod:     true,
			},
			mock:    availableMock(),
			wantErr: true,
			errMsg:  "password must be at least 8 characters",
		},
		{
			name: "username already taken",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				ResetMethod:     true,
			},
			mock:    takenUsernameMock(),
			wantErr: true,
			errMsg:  "username already exists",
		},
		{
			name: "email already taken",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				ResetMethod:     true,
			},
			mock:    takenEmailMock(),
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name: "security question missing when reset method is question",
			input: RegisterInput{
				Username:        "testuser",
				Email:           "test@test.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				ResetMethod:     false,
			},
			mock:    availableMock(),
			wantErr: true,
			errMsg:  "security question and answer are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegisterInput(tt.input, tt.mock)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if !errors.Is(err, errors.New(tt.errMsg)) && err.Error() != tt.errMsg {
					t.Errorf("expected error %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %q", err.Error())
				}
			}
		})
	}
}
