package user

import (
	"context"

	"github.com/isdzulqor/donation-hub/internal/core/model"
)

// DataStorage for driven
type DataStorage interface {
	CreateUser(ctx context.Context, user model.UserRegisterInput) (*model.User, error)
	HasEmail(ctx context.Context, email string) (bool, error)
	HasUsername(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	GetUser(ctx context.Context, input model.ListUserInput) (*[]model.UserStorage, *int64, error)
	UserHasRole(ctx context.Context, userId int64, role string) (bool, error)
}
