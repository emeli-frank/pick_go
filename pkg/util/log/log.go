package log

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog *log.Logger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

type errorOutputContainer struct {
	Error errorOutput
}

type errorOutput struct {
	Reason string
	Errors []string
}

func ServerError(w http.ResponseWriter, err error, errors []error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	errorLog.Output(2, trace)

	http.Error(w,
		formErrorOutput(http.StatusText(http.StatusInternalServerError), nil),
		http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int, errors []string) {
	http.Error(w,
		formErrorOutput(http.StatusText(status), errors),
		status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound, nil)
}

func formErrorOutput(reason string, errs []string) string {
	output := errorOutputContainer{}
	output.Error = errorOutput{
		Reason: reason,
		Errors: []string{},
	}

	if len(errs) > 0 {
		output.Error.Errors = []string{}
		for _, e := range errs {
			output.Error.Errors = append(output.Error.Errors, e)
		}
	}

	outputString, _ := json.Marshal(output)
	return string(outputString)
}
