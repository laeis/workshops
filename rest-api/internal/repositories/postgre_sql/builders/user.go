package builders

import (
	"fmt"
	"github.com/pkg/errors"
)

const FindByQuery = "SELECT id, email, password, timezone FROM users WHERE 1 = 1"

var WrongFiled = errors.New("Wrong field for find")

type user struct {
	ArgStorage
	QueryPlaceholder
	queryList map[string]func(*user) string
}

func NewUser() *user {
	return &user{
		queryList: map[string]func(*user) string{
			"email": byEmail,
			"id":    byId,
			"token": byToken,
		},
	}
}

func (u *user) FindBy(field string, value interface{}) (string, error) {
	queryFunc, ok := u.queryList[field]
	if !ok {
		return "", WrongFiled
	}
	u.AddQueryArg(value)
	return queryFunc(u), nil
}

func byEmail(u *user) string {
	return fmt.Sprintf("%s AND email=$%d LIMIT 1", FindByQuery, u.NextPlaceholder())
}

func byId(u *user) string {
	return fmt.Sprintf("%s AND id=$%d LIMIT 1", FindByQuery, u.NextPlaceholder())
}

func byToken(u *user) string {
	return fmt.Sprintf("%s AND id IN (SELECT user_id FROM tokens WHERE token = $%d GROUP BY user_id)", FindByQuery, u.NextPlaceholder())
}
