package userstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/isdzulqor/donation-hub/internal/core/model"
)

type Storage struct {
	container *model.Container
}

type DatabaseUser struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	Roles     string `db:"roles" json:"roles"`
}

type UsersCount struct {
	Total int64 `json:"total"`
}

func New(container *model.Container) *Storage {
	return &Storage{container: container}
}

func (s Storage) CreateUser(ctx context.Context, input model.UserRegisterInput) (*model.User, error) {
	ts := time.Now().Unix()
	tx, err := s.container.Connection.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`
	resUser, err := tx.ExecContext(ctx, query, input.Username, input.Email, input.Password, ts)
	if err != nil {
		return nil, err
	}

	userId, _ := resUser.LastInsertId()
	query = `INSERT INTO user_roles (user_id, role) VALUES (?,?)`
	_, err = tx.ExecContext(ctx, query, userId, input.Role)

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       userId,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}, nil
}

func (s Storage) HasEmail(ctx context.Context, email string) (bool, error) {
	query := "select count(*) from users where email = ?"
	var exists = false
	err := s.container.Connection.DB.GetContext(ctx, &exists, query, email)

	return exists, err
}

func (s Storage) HasUsername(ctx context.Context, username string) (bool, error) {
	query := "select count(*) from users where username = ?"
	var exists = false
	err := s.container.Connection.DB.GetContext(ctx, &exists, query, username)

	return exists, err
}

func (s Storage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var du DatabaseUser
	query := "SELECT id, username, email,  password, created_at FROM users WHERE username = ?"
	err := s.container.Connection.DB.GetContext(ctx, &du, query, username)

	if err != nil {
		return nil, err
	}

	// get role on user_role table
	var roles []string
	rolesQuery := "SELECT role FROM user_roles WHERE user_id = ?"
	err = s.container.Connection.DB.SelectContext(ctx, &roles, rolesQuery, du.ID)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       du.ID,
		Username: du.Username,
		Email:    du.Email,
		Password: du.Password,
		Roles:    roles,
	}, nil
}

func (s Storage) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	var du DatabaseUser
	query := `SELECT 
    			users.id as id,
    			users.username as username,
    			users.email as email,
    			GROUP_CONCAT(user_roles.role) AS roles
				FROM 
				    users 
				JOIN 
				    user_roles ON users.id = user_roles.user_id
                WHERE id = ?
                GROUP BY users.id`
	err := s.container.Connection.DB.GetContext(ctx, &du, query, id)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		ID:       du.ID,
		Username: du.Username,
		Email:    du.Email,
		Password: du.Password,
	}

	if du.Roles != "" {
		u.Roles = strings.Split(du.Roles, ",")
	}

	return u, nil
}

// GetUser total is a total data, not pagination
func (s Storage) GetUser(ctx context.Context, input model.ListUserInput) (*[]model.UserStorage, *int64, error) {
	offset := (input.Page - 1) * input.Limit
	var query string
	var count UsersCount
	var users []DatabaseUser
	var err error

	// misal create colum baru, name, maka structnya akan error

	if input.Role == "" {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role) AS roles
				FROM users
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role IN ("donor", "requester")
				GROUP BY users.id LIMIT ? OFFSET ? `

		s.container.Connection.DB.Unsafe() // salah satu cara untuk meminimalisir error struct
		err = s.container.Connection.DB.SelectContext(ctx, &users, query, input.Limit, offset)

		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("error when get user: %v", err))
		}

		err = s.container.Connection.DB.GetContext(ctx, &count, "SELECT COUNT(*) as total FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.role IN ('donor', 'requester')")
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("error when count user: %v", err))
		}
	} else {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role) AS roles
				FROM users
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role = ? GROUP BY users.id LIMIT ? OFFSET ? `

		err = s.container.Connection.DB.SelectContext(ctx, &users, query, input.Role, input.Limit, offset)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("error when get user: %v", err))
		}
		err = s.container.Connection.DB.GetContext(ctx, &count, "SELECT COUNT(*) as total FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.role = ? GROUP BY u.id", input.Role)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("error when count user: %v", err))
		}
	}

	var userStorages []model.UserStorage
	for _, user := range users {
		userStorages = append(userStorages, model.UserStorage{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
		})
	}

	return &userStorages, &count.Total, nil
}

func (s Storage) UserHasRole(ctx context.Context, userId int64, role string) (bool, error) {
	query := "select count(*) from user_roles where user_id = ? and role = ?"
	var exists = false
	err := s.container.Connection.DB.GetContext(ctx, &exists, query, userId, role)

	return exists, err
}
