package user

import (
	"errors"
	"github.com/emeli-frank/pick_go/pkg/forms/validation"
)

type User struct {
	ID        int `json:"id"`
	Names     string `json:"names"`
	Email     string `json:"email"`
}

type Service interface {
	Create(user *User, password string) (id int, err error)
	Authenticate(email string, password string) (user *User, token string, err error)
	GetUser(id int) (user *User, err error)
}

// todo:: this might not be the best place to store this
type contextKey string
const ContextKeyUser = contextKey("user")

var nameRule = validation.StringRule{
	NotEmpty:    true,
	MinLen:      1,
	MaxLen:      64,
}

var ErrAccountAlreadyConfirmed = errors.New("account is already confirmed")
