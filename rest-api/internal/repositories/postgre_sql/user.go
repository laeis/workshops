//go:generate mockgen -source taskRepository.go -destination mock/taskRepository_mock.go -package mock
package postgre_sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/repositories/postgre_sql/builders"
)

type User struct {
	Connection *sql.DB
}

func NewUser(c *sql.DB) *User {
	return &User{Connection: c}
}

func (r *User) FindBy(ctx context.Context, field string, value interface{}) (*entities.User, error) {
	user := entities.User{}
	b := builders.NewUser()
	userSql, err := b.FindBy(field, value)
	if err != nil {
		return nil, fmt.Errorf("Serch by field %s not implemented: %w ", field, appError.NotFound)
	}
	if err := r.Connection.QueryRowContext(ctx, userSql, value).Scan(&user.Id, &user.Email, &user.Password, &user.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("User with %s %s: %w ", field, value, appError.NotFound)
		}
		log.Print("Failed to execute query: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *User) Get(ctx context.Context, id string) (*entities.User, error) {
	user := entities.User{}
	userSql := "SELECT id, email, password, timezone  FROM users WHERE id = $1"
	if err := r.Connection.QueryRow(userSql, id).Scan(&user.Id, &user.Email, &user.Password, &user.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("User with Id %d: %w ", id, appError.NotFound)
		}
		log.Print("Failed to execute query: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *User) Store(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := "INSERT INTO users(email, password, timezone) VALUES ($1, $2, $3) returning id"
	var lastInsertId string
	err := r.Connection.QueryRow(query, user.Email, user.Password, user.Timezone).Scan(&lastInsertId)
	if err != nil {
		log.Printf("Error %s when inserting row into user table", err)
		return nil, err
	}
	user.Id = lastInsertId
	return user, nil
}

func (r *User) Update(ctx context.Context, id string, user *entities.User) (*entities.User, error) {
	query := "update users set timezone=$1 where id=$2"
	_, err := r.Connection.Exec(query, user.Timezone, id)
	if err != nil {
		log.Printf("Error %s when update user %s", err, id)
		return nil, err
	}
	updUser, err := r.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("User with Id %s: %w ", id, appError.NotFound)
		}
		log.Print("Failed to execute query: ", err)
		return nil, err
	}
	return updUser, nil
}

func (r *User) AddToken(ctx context.Context, user *entities.User, token string) error {
	query := "INSERT INTO tokens(token, user_id) VALUES ($1, $2)"
	result, err := r.Connection.Exec(query, token, user.Id)
	cnt, _ := result.RowsAffected()
	if err != nil || cnt == 0 {
		log.Printf("Error %s when inserting row into user table", err)
		return err
	}
	return nil
}

func (r *User) DeleteToken(ctx context.Context, user *entities.User, token string) error {
	query := "DELETE FROM tokens WHERE token=$1 AND user_id=$2"
	result, err := r.Connection.Exec(query, token, user.Id)
	fmt.Println(err)
	cnt, _ := result.RowsAffected()
	if err != nil || cnt == 0 {
		log.Printf("Error %s when deleting row into user table", err)
		return err
	}
	return nil
}

func (r *User) GetTimezone(ctx context.Context, userId string) (string, error) {
	user, err := r.Get(ctx, userId)
	if err != nil {
		return "", errors.New("Some context auth error")
	}
	return user.Timezone, nil
}
