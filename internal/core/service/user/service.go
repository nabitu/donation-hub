package user

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/auth"
)

type Storage struct {
	storage   DataStorage
	authToken auth.Service
}

type Service interface {
	Register(context.Context, model.UserRegisterInput) (*model.UserRegisterOutput, error)
	Login(context.Context, model.UserLoginInput) (*model.UserLoginOutput, error)
	ListUser(context.Context, model.ListUserInput) (*model.ListUserOutput, error)
}

func New(storage DataStorage, authToken auth.Service) Service {
	return &Storage{
		storage:   storage,
		authToken: authToken,
	}
}

func (s *Storage) Register(ctx context.Context, input model.UserRegisterInput) (*model.UserRegisterOutput, error) {
	// todo: add validation for input here

	hasEmail, err := s.storage.HasEmail(ctx, input.Email)
	if hasEmail {
		return nil, errors.New("email already exists")
	}

	hasUsername, err := s.storage.HasUsername(ctx, input.Username)
	if (err != nil) || (hasUsername) {
		return nil, errors.New("username already exists")
	}

	u, err := s.storage.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return &model.UserRegisterOutput{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (s *Storage) Login(ctx context.Context, input model.UserLoginInput) (*model.UserLoginOutput, error) {
	// todo: add validation for input here

	user, err := s.storage.GetUserByUsername(ctx, input.Username)
	if err != nil || user.Password != input.Password {
		return nil, errors.New("invalid username or password")
	}

	tokenPayload := model.AuthPayload{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	accessToken, err := s.authToken.GenerateToken(tokenPayload)
	if err != nil {
		return nil, err
	}

	return &model.UserLoginOutput{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *Storage) ListUser(ctx context.Context, input model.ListUserInput) (output *model.ListUserOutput, err error) {
	// todo: add validation for input here

	userStorages, total, err := s.storage.GetUser(ctx, input)
	if err != nil {
		return nil, err
	}

	listUsers := make([]model.User, len(userStorages))
	for i, us := range userStorages {
		roles := strings.Split(us.Roles, ",")
		listUser := model.User{
			ID:       us.ID,
			Username: us.Username,
			Email:    us.Email,
			Roles:    roles,
		}
		listUsers[i] = listUser
	}

	// pagination
	totalPage := int64(math.Ceil(float64(total / input.Limit)))
	if total%input.Limit != 0 {
		totalPage++
	}

	output = &model.ListUserOutput{
		Users: listUsers,
		Pagination: model.ListUserMeta{
			Page:       input.Page,
			TotalPages: totalPage,
		},
	}

	return
}
