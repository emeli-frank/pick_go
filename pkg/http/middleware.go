package http

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emeli-frank/pick_go/pkg/domain/user"
	"github.com/emeli-frank/pick_go/pkg/util/log"
	"net/http"
	"strings"
)

func (s server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				log.ServerError(w, fmt.Errorf("%s", err), nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s server) checkJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerTokenSlice := strings.Split(r.Header.Get("Authorization"), " ")
		if len(bearerTokenSlice) != 2 || bearerTokenSlice[0] != "Bearer" {
			//fmt.Println("token not formed properly in header")
			next.ServeHTTP(w, r)
			return
		}
		tokenStr := bearerTokenSlice[1]

		type claims struct {
			User user.User
			jwt.StandardClaims
		}

		c := &claims{}

		_, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secrete_key"), nil
		})
		if err != nil {
			// todo:: add this commented part to the part that fails when no valid jwt is provided
			/*log.ClientError(w, http.StatusForbidden, []string{"bad authorization blah blah ..."})
			return*/
			fmt.Println("invalid token")
			next.ServeHTTP(w, r)
			return
		}

		u := c.User

		ctx := context.WithValue(r.Context(), user.ContextKeyUser, &u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*func (s server) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerTokenSlice := strings.Split(r.Header.Get("Authorization"), " ")
		if len(bearerTokenSlice) != 2 || bearerTokenSlice[0] != "Bearer"{
			err := errors.New("bad authorization blah blah ...")
			err = errors2.Wrap(err, "authenticatedOnly", "bearers token") // todo:: fix
			s.response.ClientError(w, http.StatusForbidden, err)
			return
		}
		tokenStr := bearerTokenSlice[1]

		type claims struct {
			User user.User
			jwt.StandardClaims
		}

		c := &claims{}

		_, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secrete_key"), nil
		})
		if err != nil {
			log.ClientError(w, http.StatusForbidden, []string{"bad authorization blah blah ..."})
			return
		}
		next.ServeHTTP(w, r)
	})
}*/

func (s server) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(user.ContextKeyUser).(*user.User)
		if  !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(">>>>>>>>>>>", u, "<<<<<<<<<<<<<<<<")
		next.ServeHTTP(w, r)
	})
}
