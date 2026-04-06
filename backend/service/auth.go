package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/carnex/fittrack/backend/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	store     store.Store
	jwtSecret string
}

type LoginInput struct {
	Username string
	Password string
}

type LoginResult struct {
	Token string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	var result LoginResult
	user, err := s.store.GetUserByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}
	errPasswordMismatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errPasswordMismatch != nil {
		slog.Warn("login attempt failed: wrong password", "username", input.Username)
		return LoginResult{}, ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenstring, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return LoginResult{}, err
	}
	result.Token = tokenstring
	return result, nil
}

func NewAuthService(store store.Store, jwtSecret string) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}
