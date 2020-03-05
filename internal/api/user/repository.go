package user

import (
	"context"

	"github.com/ldej/go-rest-example/internal/api"
)

type Repository interface {
	UserCreate(context.Context, *api.User) (*api.User, error)
	UserGetByUID(context.Context, string) (*api.User, error)
	UserGetByEmailAddress(context.Context, string) (*api.User, error)
}
