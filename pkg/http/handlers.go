package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/domain/user"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
	"github.com/emeli-frank/pick_go/pkg/forms/validation"
	"log"
	"net/http"
)

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decodeError := errors.New("could not decode")
	if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type") != "application/json" {
		return decodeError
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		return decodeError
	}

	if dec.More() {
		return decodeError
	}

	return nil
}

type server struct {
	response *response
	userService user.Service
	infoLog *log.Logger
}

func NewServer(response *response, userService user.Service, infoLog *log.Logger) *server {

	return &server{
		response: response,
		userService: userService,
		infoLog: infoLog,
	}
}

func (s server) registerHandler(w http.ResponseWriter, r *http.Request) {
	op := "server.registerHandler"
	data := &struct {
		user.User
		Password string `json:"password"`
	}{}

	err := decodeJSONBody(w, r, data)
	if err != nil {
		err = errors2.WrapWithMsg(err, op, "", "error decoding request body")
		s.response.ClientError(w, http.StatusBadRequest, err)
		return
	}

	password := data.Password

	id, err := s.userService.Create(&data.User, password)
	if err != nil {
		switch errors2.Unwrap(err).(type) {
		case validation.Errors:
			s.response.ClientError(w, http.StatusBadRequest, err)
			return
		case *errors2.Conflict:
			s.response.ClientError(w, http.StatusConflict, err)
			return
		default:
			s.response.ServerError(w, err)
			return
		}
	}

	clientOutput := struct {
		ID               int `json:"id"`
	}{
		ID:                       id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/users/%d", id))
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(clientOutput)
}

func (s server) loginHandler(w http.ResponseWriter, r *http.Request) {
	op := "server.loginHandler"
	credentials := struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}{}

	err := decodeJSONBody(w, r, &credentials)
	if err != nil {
		err = errors2.WrapWithMsg(err, op, "", "error decoding request body")
		s.response.ClientError(w, http.StatusBadRequest, err)
		return
	}

	email := credentials.Email
	password := credentials.Password

	u, authToken, err := s.userService.Authenticate(email, password)
	if err != nil {
		err = errors2.Wrap(err, op, "getting user and auth token from userService")
		switch errors2.Unwrap(err).(type) {
		case validation.Errors:
			s.response.ClientError(w, http.StatusBadRequest, err)
			return
		case *errors2.Unauthorized:
			err = errors2.SetMessage(err, "email or password is not correct")
			// todo:: find a way to send empty message with unauthorized header
			s.response.ClientError(w, http.StatusUnauthorized, err)
			return
		default:
			s.response.ServerError(w, err)
			return
		}
	}

	clientOutput := struct {
		AuthorizationToken string `json:"authorization_token"`
		User               user.User `json:"user"`
	}{
		AuthorizationToken:         authToken,
		User:                       *u,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clientOutput)
}

func (s server) userHandler(w http.ResponseWriter, r *http.Request) {
	op := "server.userHandler"

	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		s.response.ServerError(w, errors2.Wrap(errors.New(""), op, "getting user object from request context"))
		return
	}

	u, err := s.userService.GetUser(u.ID)
	if err != nil {
		err = errors2.Wrap(err, op, "getting user from userService")
		switch errors2.Unwrap(err).(type) {
		case *errors2.NotFound:
			err = errors2.SetMessage(err, "user does not exist")
			s.response.ClientError(w, http.StatusNotFound, err)
			return
		default:
			s.response.ServerError(w, err)
			return
		}
	}

	clientOutput := struct {
		User               user.User `json:"user"`
	}{
		User:                       *u,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clientOutput)
}