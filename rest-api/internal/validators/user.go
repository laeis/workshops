package validators

import (
	"github.com/pkg/errors"
	"net/mail"
	"time"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
)

var emailRules = func(u *entities.User) error {
	if u.Email == "" {
		return errors.Wrapf(appError.BadRequest, "User email empty")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.Wrapf(appError.BadRequest, "User email not valid: %w", err)
	}
	return nil
}

var passwordRules = func(u *entities.User) error {
	if "" == u.Password {
		return errors.Wrap(appError.BadRequest, "User password empty")
	}
	return nil
}

var timezoneRules = func(u *entities.User) error {
	if _, err := time.LoadLocation(u.Timezone); err != nil {
		return errors.Wrapf(appError.BadRequest, "User timezone not valid: %w", err)
	}
	return nil
}

type user struct {
	fields stringList
	rules  map[string]func(*entities.User) error
	entity *entities.User
}

//UserValidator create new validator instance and fill required fields
func UserValidator(u *entities.User) *user {
	return &user{
		entity: u,
		fields: stringList{"emails", "password", "timezone"},
		rules: map[string]func(*entities.User) error{
			"emails":   emailRules,
			"password": passwordRules,
			"timezone": timezoneRules,
		},
	}
}

//Validate chek User fields by rules
func (u *user) Validate(fields ...string) error {
	vld := stringList(fields)
	if len(vld) == 0 {
		vld = u.fields
	}

	for _, v := range vld {
		if rules, ok := u.rules[v]; ok {
			if err := rules(u.entity); err != nil {
				return err
			}
		}
	}

	return nil
}
