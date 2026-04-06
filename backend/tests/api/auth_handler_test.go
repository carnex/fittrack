package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/carnex/fittrack/backend/db/gen"
	"github.com/carnex/fittrack/backend/handlers"
	"github.com/carnex/fittrack/backend/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
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

// availableMock — empty database, all registrations succeed
func availableMock() *MockStore {
	return &MockStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
			return db.GetUserByUsernameRow{}, pgx.ErrNoRows
		},
		getUserByEmailFunc: func(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
			return db.GetUserByEmailRow{}, pgx.ErrNoRows
		},
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, nil
		},
	}
}

// validUserMock — returns a real bcrypt hash for "password123"
func validUserMock(t *testing.T) *MockStore {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return &MockStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
			return db.GetUserByUsernameRow{
				ID:       pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true},
				Username: username,
				Password: string(hash),
			}, nil
		},
		getUserByEmailFunc: func(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
			return db.GetUserByEmailRow{}, pgx.ErrNoRows
		},
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, nil
		},
	}
}

// userNotFoundMock — simulates username not in DB
func userNotFoundMock() *MockStore {
	return &MockStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
			return db.GetUserByUsernameRow{}, pgx.ErrNoRows
		},
		getUserByEmailFunc: func(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
			return db.GetUserByEmailRow{}, pgx.ErrNoRows
		},
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, nil
		},
	}
}

func registerBody(t *testing.T, input map[string]interface{}) *bytes.Buffer {
	t.Helper()
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	return bytes.NewBuffer(body)
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		mock       *MockStore
		wantStatus int
	}{
		{
			name: "valid registration",
			body: map[string]interface{}{
				"Username":        "testuser",
				"Email":           "test@test.com",
				"Password":        "password123",
				"ConfirmPassword": "password123",
				"ResetMethod":     true,
			},
			mock:       availableMock(),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json body",
			body:       nil,
			mock:       availableMock(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "passwords do not match",
			body: map[string]interface{}{
				"Username":        "testuser",
				"Email":           "test@test.com",
				"Password":        "password123",
				"ConfirmPassword": "different",
				"ResetMethod":     true,
			},
			mock:       availableMock(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			body: map[string]interface{}{
				"Username":        "testuser",
				"Email":           "test@test.com",
				"Password":        "short",
				"ConfirmPassword": "short",
				"ResetMethod":     true,
			},
			mock:       availableMock(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing security question when reset method is false",
			body: map[string]interface{}{
				"Username":        "testuser",
				"Email":           "test@test.com",
				"Password":        "password123",
				"ConfirmPassword": "password123",
				"ResetMethod":     false,
			},
			mock:       availableMock(),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.body == nil {
				reqBody = bytes.NewBufferString("invalid json{{{")
			} else {
				reqBody = registerBody(t, tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			userService := service.NewUserService(tt.mock)
			handler := handlers.NewAuthHandler(userService, nil)
			handler.Register(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d got %d — body: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		mock       *MockStore
		wantStatus int
		wantToken  bool
	}{
		{
			name: "valid login returns token",
			body: map[string]interface{}{
				"Username": "badlinemtb",
				"Password": "password123",
			},
			mock:       nil, // set per test using t
			wantStatus: http.StatusOK,
			wantToken:  true,
		},
		{
			name: "username not found returns 401",
			body: map[string]interface{}{
				"Username": "doesnotexist",
				"Password": "password123",
			},
			mock:       userNotFoundMock(),
			wantStatus: http.StatusUnauthorized,
			wantToken:  false,
		},
		{
			name: "wrong password returns 401",
			body: map[string]interface{}{
				"Username": "badlinemtb",
				"Password": "wrongpassword",
			},
			mock:       nil, // set per test using t
			wantStatus: http.StatusUnauthorized,
			wantToken:  false,
		},
		{
			name:       "invalid json returns 400",
			body:       nil,
			mock:       userNotFoundMock(),
			wantStatus: http.StatusBadRequest,
			wantToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// build mock — some tests need the bcrypt hash
			mock := tt.mock
			if mock == nil {
				mock = validUserMock(t)
			}

			// build request body
			var reqBody *bytes.Buffer
			if tt.body == nil {
				reqBody = bytes.NewBufferString("invalid json{{{")
			} else {
				reqBody = registerBody(t, tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// wire up handler with both services
			userService := service.NewUserService(mock)
			authService := service.NewAuthService(mock, "test-secret-key")
			handler := handlers.NewAuthHandler(userService, authService)
			handler.Login(rr, req)

			// assert status code
			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d got %d — body: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			// assert token presence
			if tt.wantToken {
				var resp map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Errorf("failed to decode response body: %v", err)
					return
				}
				if resp["token"] == "" {
					t.Errorf("expected token in response but got empty string")
				}
			}
		})
	}
}
