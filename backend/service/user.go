package service

import (
	"context"
	"errors"
	"strings"

	"database/sql"

	db "github.com/carnex/fittrack/backend/db/gen"
	"github.com/carnex/fittrack/backend/store"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store store.Store
}

type RegisterInput struct {
	Username         string
	Email            string
	Password         string
	ConfirmPassword  string
	ResetMethod      bool
	SecurityQuestion string
	SecurityAnswer   string
}

func NewUserService(store store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) error {
	var answer []byte
	errValidtate := validateRegisterInput(input, s.store)
	if errValidtate != nil {
		return errValidtate
	}
	password, errHash := hashString(input.Password)
	if errHash != nil {
		return errHash
	}
	if !input.ResetMethod {
		securityAnswer, err := hashString(input.SecurityAnswer)
		if err != nil {
			return err
		}
		answer = securityAnswer
	}
	_, errUser := s.store.CreateUser(ctx, db.CreateUserParams{Username: input.Username, Email: input.Email, Password: string(password), ResetMethod: input.ResetMethod, SecurityQuestion: nullString(input.SecurityQuestion), SecurityQuestionAnswer: nullString(string(answer))})
	if errUser != nil {
		return errUser
	}

	return nil
}

func validateRegisterInput(input RegisterInput, store store.Store) error {
	ctx := context.Background()
	if input.Password != input.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	_, userErr := store.GetUserByUsername(ctx, input.Username)
	if userErr == nil {
		return errors.New("username already exists")
	}
	if !errors.Is(userErr, pgx.ErrNoRows) {
		return userErr
	}
	_, emailErr := store.GetUserByEmail(ctx, input.Email)
	if emailErr == nil {
		return errors.New("email already exists")
	}
	if !errors.Is(emailErr, pgx.ErrNoRows) {
		return emailErr
	}
	if !input.ResetMethod {
		if strings.TrimSpace(input.SecurityQuestion) == "" || strings.TrimSpace(input.SecurityAnswer) == "" {
			return errors.New("security question and answer are required")
		}
	}
	return nil
}

func hashString(str string) ([]byte, error) {
	value := []byte(str)
	hash, err := bcrypt.GenerateFromPassword(value, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
