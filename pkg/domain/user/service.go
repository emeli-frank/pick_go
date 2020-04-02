package user

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
	"github.com/emeli-frank/pick_go/pkg/forms/validation"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type FormErrors []string
var ErrForm error = errors.New("form error")

type repository interface {
	SaveUser(user *User, hashedPassword string) (int, error)
	GetUser(userId int) (user *User, err error)
	GetUserToAuthenticate(email string) (user *User, hashedPassword string, err error)
	Tx() (*sql.Tx, error)
}

func New(repo repository) *service {
	return &service{
		r: repo,
	}
}

type service struct {
	r repository
}

func (s *service) Create(user *User, password string) (int, error) {
	const op = "userService.Create"

	// validation
	v := validation.New(nil)
	v.Do("first_name", op, user.Names, nameRule)
	v.Do("email", op, user.Email, validation.EmailRule)
	v.Do("password", op, password, validation.PasswordRule)
	if err := v.Errors(); err != nil {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, errors2.Wrap(err, op, "generating passwords")
	}

	id, err := s.r.SaveUser(user, string(hashedPassword))
	user.ID = id
	if err != nil {
		switch errors2.Unwrap(err).(type) {
		case *errors2.Conflict:
			return 0, errors2.Wrap(err, op, "the email is not available")
			/*v.AddError("email", "the email is not available", op, err)
			return 0, "", v.Errors()*/
		default:
			return 0, errors2.Wrap(err, op, "creating user from repo")
		}
	}

	return id, nil
}

// todo:: consider moving to a separate auth service
func (s *service) Authenticate(email string, password string) (*User, string, error){
	op := "userService.Authenticate"

	v := validation.New(nil)
	v.Do("email", op, email, validation.EmailRule)

	if err := v.Errors(); err != nil {
		return nil, "", err
	}

	// get user from repo
	user, hashedPassword, err := s.r.GetUserToAuthenticate(email)
	if err != nil {
		switch errors2.Unwrap(err).(type) {
		case *errors2.NotFound:
			return nil, "", errors2.Wrap(&errors2.Unauthorized{Err: err}, op, "getting user to authenticate")
		default:
			return nil, "", errors2.Wrap(err, op, "getting user to authenticate")
		}
	}

	// compare user provided and stored password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, "", errors2.Wrap(&errors2.Unauthorized{Err: err}, op, "hashing password")
	} else if err != nil {
		return nil, "", errors2.Wrap(err, op, "hashing password")
	}

	// generate jwt
	jwtKey := []byte("my_secrete_key") // todo:: store somewhere else

	type claims struct {
		User User `json:"user"`
		jwt.StandardClaims
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	c := &claims{
		User:           User{ID: user.ID, Email: email},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, "", errors2.Wrap(err, op,"signing token")
	}

	return user, tokenString, nil
}

func (s *service) GetUser(id int) (*User, error) {
	const op = "userService.GetUser"
	user, err := s.r.GetUser(id)
	if err != nil {
		return nil, errors2.Wrap(err, op, "getting user from user repo")
	}

	return user, nil
}

func generateToken() string {
	return strconv.Itoa(int(time.Now().Unix()))
}
