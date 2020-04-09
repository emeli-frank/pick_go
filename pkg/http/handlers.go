package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/domain/product"
	"github.com/emeli-frank/pick_go/pkg/domain/user"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
	"github.com/emeli-frank/pick_go/pkg/forms/validation"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
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
	productService product.Service
}

func NewServer(response *response, userService user.Service, productService product.Service, infoLog *log.Logger) *server {

	return &server{
		response: response,
		userService: userService,
		infoLog: infoLog,
		productService: productService,
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

func (s server) productListHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.productListHandler"

	page, err := urlParamToInt(r, "page")
	if err != nil || page <= 0 {
		page = 1
	}
	numberToFetch, err := urlParamToInt(r, "number")
	if err != nil {
		numberToFetch = 30
	}

	index := (page - 1) * numberToFetch

	pp, total, err := s.productService.GetProducts(index, numberToFetch)
	if err != nil {
		s.response.ServerError(w, errors2.Wrap(err, op, "getting products from service"))
		return
	}

	o := struct {
		CurrentPage int `json:"current_page"`
		NumberLoaded int `json:"number_loaded"`
		TotalNumber int `json:"total_number"`
		Products []*product.Product `json:"products"`
	} {
		CurrentPage: page,
		NumberLoaded: len(pp),
		TotalNumber: total,
		Products: pp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(o)
}

func (s server) productDetailHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.productDetailHandler"

	vars := mux.Vars(r)
	pIdstr := vars["id"]
	productId, err := strconv.Atoi(pIdstr)
	if err != nil {
		err = errors2.Wrap(err, op, "getting product id from URL")
		s.response.ClientError(w, http.StatusNotFound, err)
		return
	}

	var userId int
	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		userId = 0
	} else {
		userId = u.ID
	}

	p, inCart, err := s.productService.GetProduct(productId, userId)
	if err != nil {
		switch errors2.Unwrap(err).(type) {
		case *errors2.NotFound:
			err = errors2.Wrap(err, op, "getting product from service")
			s.response.ClientError(w, http.StatusNotFound, err)
			return
		default:
			s.response.ServerError(w, errors2.Wrap(err, op, "getting product from service"))
			return
		}
	}

	o := struct {
		Product *product.Product `json:"product"`
		InCart bool `json:"in_cart"`
	} {
		Product: p,
		InCart: inCart,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(o)
}

func urlParamToInt(r *http.Request, key string) (int, error){
	paramMap := r.URL.Query()[key]
	if len(paramMap) != 1 {
		return 0, errors.New("param value is not 1")
	}

	value, err := strconv.Atoi(paramMap[0])
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s server) cartItemsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.cartItemsHandler"

	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		s.response.ServerError(w, errors2.Wrap(errors.New(""), op, "getting user object from request context"))
		return
	}

	pp, err := s.productService.GetCartItems(u.ID)
	if err != nil {
		s.response.ServerError(w, errors2.Wrap(err, op, "getting products from service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pp)
}

func (s server) saveProductToCartHandler(w http.ResponseWriter, r *http.Request) {
	op := "server.saveProductToCartHandler"
	data := struct {
		ProductId int `json:"product_id"`
	}{}

	err := decodeJSONBody(w, r, &data)
	if err != nil {
		err = errors2.WrapWithMsg(err, op, "", "error decoding request body")
		s.response.ClientError(w, http.StatusBadRequest, err)
		return
	}

	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		s.response.ServerError(w, errors2.Wrap(errors.New(""), op, "getting user object from request context"))
		return
	}

	err = s.productService.SaveProductToCart(u.ID, data.ProductId)
	if err != nil {
		err = errors2.Wrap(err, op, "getting user and auth token from userService")
		switch errors2.Unwrap(err).(type) {
		case *errors2.Conflict:
			err = errors2.SetMessage(err, "duplicate entry")
			s.response.ClientError(w, http.StatusConflict, err)
			return
		default:
			s.response.ServerError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s server) orderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.orderHistoryHandler"

	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		s.response.ServerError(w, errors2.Wrap(errors.New(""), op, "getting user object from request context"))
		return
	}

	pp, err := s.productService.GetOrderProducts(u.ID)
	if err != nil {
		s.response.ServerError(w, errors2.Wrap(err, op, "getting products from service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pp)
}

func (s server) saveToOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	op := "server.saveProductToCartHandler"
	data := struct {
		ProductId int `json:"product_id"`
	}{}

	err := decodeJSONBody(w, r, &data)
	if err != nil {
		err = errors2.WrapWithMsg(err, op, "", "error decoding request body")
		s.response.ClientError(w, http.StatusBadRequest, err)
		return
	}

	u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
	if  !ok {
		s.response.ServerError(w, errors2.Wrap(errors.New(""), op, "getting user object from request context"))
		return
	}

	err = s.productService.SaveToOrderHistory(u.ID, data.ProductId, time.Now())
	if err != nil {
		err = errors2.Wrap(err, op, "getting user and auth token from userService")
		switch errors2.Unwrap(err).(type) {
		case *errors2.Conflict:
			err = errors2.SetMessage(err, "duplicate entry")
			s.response.ClientError(w, http.StatusConflict, err)
			return
		default:
			s.response.ServerError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
