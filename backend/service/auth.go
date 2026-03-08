package service

import (
	"github.com/carnex/fittrack/backend/store"
)

type AuthService struct {
	store     store.Store
	jwtSecret string
}

func NewAuthService(store store.Store, jwtSecret string) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}
