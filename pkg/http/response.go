package http

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
)

type response struct {
	errorLog *log.Logger
}

func NewResponse(errorLog *log.Logger) *response {
	return &response{
		errorLog: errorLog,
	}
}

func (r response) ClientError(w http.ResponseWriter, code int, err error) {
	r.Respond(w, code, err)
}

func (r response) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = r.errorLog.Output(2, trace)
	r.Respond(w, http.StatusInternalServerError, nil)
}

func (r response) Respond(w http.ResponseWriter, code int, err error) {
	type errorOut struct {
		Error interface{} `json:"error"`
	}
	var out errorOut

	if err == nil {
		out = errorOut{Error: "Internal server error"}
	} else {
		if _, ok := err.(json.Marshaler); !ok {
			fmt.Printf("Error of type: %T does not implements marshaller\n", err)
			err := errors2.Wrap(err, "httpResponse.Respond",
				"error does not implement marshaler so could not be marshalled for client")

			trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
			_ = r.errorLog.Output(2, trace)

			out = errorOut{Error:"an unknown error occurred"}
			return
		} else {
			fmt.Printf("Error of type: %T implements marshaller\n", err)
			out = errorOut{Error: err}
		}
	}

	errStr, err := json.Marshal(out)
	if err != nil {
		// todo:: handle
		fmt.Println("server.Respond(): cannot marshal error for client")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = fmt.Fprintln(w, string(errStr))
}
