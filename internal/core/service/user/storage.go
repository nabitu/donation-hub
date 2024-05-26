package user

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

// DataStorage for driven
type DataStorage interface {
	CreateUser(ctx context.Context, user model.UserRegisterInput) (model.User, error)
	HasEmail(ctx context.Context, email string) (has bool, err error)
	HasUsername(ctx context.Context, username string) (has bool, err error)
	GetUserByUsername(ctx context.Context, username string) (user model.User, err error)
	GetUserById(ctx context.Context, id int64) (user model.User, err error)
	GetUser(ctx context.Context, input model.ListUserInput) (users []model.UserStorage, total int64, err error)
	UserHasRole(ctx context.Context, userId int64, role string) (ok bool, err error)
}
