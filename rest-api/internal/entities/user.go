package entities

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Token    JWT    `json:"-"`
	Tasks    Tasks  `json:"-"`
}

func EmptyUser() User {
	return User{
		Timezone: time.UTC.String(),
	}
}

// HashPassword encrypts user password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword checks user password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}
