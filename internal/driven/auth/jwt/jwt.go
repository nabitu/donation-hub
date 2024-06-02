package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/auth"
	"strconv"
	"time"
)

type service struct {
	container *model.Container
}

type MyCustomClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
	jwt.RegisteredClaims
}

func (s service) GenerateToken(i model.AuthPayload) (token string, err error) {
	userIDStr := strconv.FormatInt(i.UserID, 10)
	claims := MyCustomClaims{
		UserID:   userIDStr,
		Username: i.Username,
		Email:    i.Email,
		Role:     i.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    s.container.Config.TokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "for user " + i.Username,
			Audience:  jwt.ClaimStrings{},
		},
	}
	fmt.Println(userIDStr)
	fmt.Println(claims)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = jwtToken.SignedString([]byte(s.container.Config.TokenSecretKey))

	return
}

func (s service) ValidateToken(tokenString string) (*model.AuthPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.container.Config.TokenSecretKey), nil
	})

	if err != nil {
		err = fmt.Errorf("error parsing token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		fmt.Println("id:"+claims.UserID, "username:"+claims.Username, "email:"+claims.Email, "role:"+claims.Role[0])
		id, err := strconv.Atoi(claims.UserID)
		if err != nil {
			return nil, err
		}

		return &model.AuthPayload{
			UserID:   int64(id),
			Username: claims.Username,
			Email:    claims.Email,
			Role:     claims.Role,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func New(container *model.Container) auth.Service {
	return &service{
		container: container,
	}
}
