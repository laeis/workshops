//go:generate mockgen -source taskService.go -destination mock/taskService_mock.go -package mock
package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"net/http"
	"workshops/rest-api/internal/entities"
)

type UserRepository interface {
	Update(ctx context.Context, id string, task *entities.User) (*entities.User, error)
	Get(ctx context.Context, id string) (*entities.User, error)
	FindBy(ctx context.Context, field string, value interface{}) (*entities.User, error)
	Store(ctx context.Context, user *entities.User) (*entities.User, error)
	AddToken(ctx context.Context, user *entities.User, token string) error
	DeleteToken(ctx context.Context, user *entities.User, token string) error
	GetTimezone(ctx context.Context, userId string) (string, error)
}

type SecurityToken interface {
	ParseToken(r *http.Request) (string, error)
	GenerateToken(email string) (string, error)
	ValidateToken(signedToken string) (jwt.Claims, error)
}

type User struct {
	repo  UserRepository
	token SecurityToken
}

func NewUser(r UserRepository, t SecurityToken) *User {
	return &User{
		repo:  r,
		token: t,
	}
}

func (s *User) Update(ctx context.Context, id string, user *entities.User) (*entities.User, error) {
	return s.repo.Update(ctx, id, user)
}

func (s *User) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	if err := user.HashPassword(user.Password); err != nil {
		return nil, errors.Wrap(err, "Cant make hash for password")
	}
	return s.repo.Store(ctx, user)
}

func (s *User) Get(ctx context.Context, id string) (*entities.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *User) Login(ctx context.Context, userPayload entities.User) (*entities.JWT, error) {

	user, err := s.repo.FindBy(ctx, "email", userPayload.Email)
	if err != nil {
		return nil, err
	}
	err = user.CheckPassword(userPayload.Password)

	signedToken, err := s.token.GenerateToken(user.Email)
	if err != nil {
		return nil, err
	}
	err = s.repo.AddToken(ctx, user, signedToken)
	if err != nil {
		return nil, err
	}
	return &entities.JWT{
		Token: signedToken,
	}, nil
}

func (s *User) Logout(ctx context.Context, id string, token string) (bool, error) {
	fmt.Println(id, token)
	user, err := s.repo.FindBy(ctx, "id", id)
	if err != nil {
		return false, err
	}
	err = s.repo.DeleteToken(ctx, user, token)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *User) FindByToken(ctx context.Context, field string) (*entities.User, error) {
	user, err := s.repo.FindBy(ctx, "token", field)
	if err != nil {
		return nil, err
	}
	return user, nil
}
