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
			body:       nil, // we'll send malformed JSON
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
