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
	Me(context.Context, model.UserMeInput) (*model.User, error)
}

func New(storage DataStorage, authToken auth.Service) Service {
	return &Storage{
		storage:   storage,
		authToken: authToken,
	}
}

func (s *Storage) Register(ctx context.Context, input model.UserRegisterInput) (*model.UserRegisterOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

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
	// validation for input here
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.storage.GetUserByUsername(ctx, input.Username)
	if err != nil || user.Password != input.Password {
		return nil, errors.New("invalid username or password")
	}

	tokenPayload := model.AuthPayload{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Roles,
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
	err = input.Validate()
	if err != nil {
		return nil, err
	}

	userStorages, total, err := s.storage.GetUser(ctx, input)
	if err != nil {
		return nil, err
	}

	var listUsers []model.User
	for _, us := range *userStorages {
		roles := strings.Split(us.Roles, ",")

		listUsers = append(listUsers, model.User{
			ID:       us.ID,
			Username: us.Username,
			Email:    us.Email,
			Roles:    roles,
		})
	}

	// pagination
	totalPage := int64(math.Ceil(float64(*total / input.Limit)))
	if *total%input.Limit != 0 {
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

func (s *Storage) Me(ctx context.Context, i model.UserMeInput) (*model.User, error) {
	me, err := s.storage.GetUserById(ctx, i.UserID)
	if err != nil {
		return nil, err
	}

	return me, nil
}
