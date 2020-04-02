package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SetMessage sets public message on error wrapper
func SetMessage(err error, pubMsg string) error {
	// fails silently if error is not of type *wrap
	w, ok := err.(*wrap)
	if ok {
		w.pubMsg = pubMsg
	}
	return err
}

// Message returns public error message to display to client or empty string if error is not wrapped
func Message(err error) string {
	w, ok := err.(*wrap)
	// return empty string if err is not of type *wrap
	if !ok {
		return ""
	}

	if w.pubMsg != "" {
		return w.pubMsg
	} else if w.Err != nil {
		return Message(w.Err)
	} else {
		return ""
	}
}

// WrapWithMsg wraps error and sets public message to be displayed to client
func WrapWithMsg(err error, op string, privMsg string, pubMsg string) error {
	w := Wrap(err, op, privMsg)
	w = SetMessage(w, pubMsg)
	return Wrap(w, op, pubMsg)
}

// Unwrap recursively finds non-wrap error and returns it
// e.g NotFound, Conflict etc...
func Unwrap(err error) error {
	u, ok := err.(interface{
		Unwrap() error
	})
	if !ok {
		return err
	}
	return Unwrap(u.Unwrap())
}

// Wrap adds context to already wrapped error or wraps it if it is not
func Wrap(err error, op string, message string) error {
	return &wrap{Err: err, Op:op, privMsg:message}
}

/*func As(err error, target interface{}) bool {
	u := Cause(err)
	return errors.As(err, u)
}*/

type wrap struct {
	Err     error
	Op      string
	privMsg string
	pubMsg  string
}

func (c wrap) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if c.Op != "" {
		_, _ = fmt.Fprintf(&buf, "[%s]: ", c.Op)
	} else {
		_, _ = fmt.Fprint(&buf, "_: ")
	}

	// Print the current additional context in our stack, if any.
	if c.privMsg != "" {
		_, _ = fmt.Fprintf(&buf, "[%s] >> ", c.privMsg)
	} else {
		_, _ = fmt.Fprint(&buf, "_ >> ")
	}

	// If wrapping an error, print its Error() message. Otherwise print the error code & message.
	if c.Err != nil {
		buf.WriteString(c.Err.Error())
	} else {
		_, _ = fmt.Fprintf(&buf, "<Generic error> ")
		buf.WriteString(c.privMsg)
	}

	return buf.String()
}

func(c *wrap) Unwrap() error {
	return c.Err
}

func (c *wrap) MarshalJSON() ([]byte, error) {
	e, ok := Unwrap(c).(json.Marshaler)
	if ok {
		return e.MarshalJSON()
	}

	m := Message(c)
	return json.Marshal(struct {
		Type string `json:"type"`
		Error string `json:"error"`
	}{Type:"", Error:m})
}


/*type Error struct {
	//Code    string
	message string
	context string
	op      string
	Err     error
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.op != "" {
		_, _ = fmt.Fprintf(&buf, "[%s]: ", e.op)
	} else {
		_, _ = fmt.Fprint(&buf, "_: ")
	}

	// Print the current additional context in our stack, if any.
	if e.context != "" {
		_, _ = fmt.Fprintf(&buf, "[%s] >> ", e.context)
	} else {
		_, _ = fmt.Fprint(&buf, "_ >> ")
	}

	// If wrapping an error, print its Error() message. Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		_, _ = fmt.Fprintf(&buf, "<Generic error> ")
		buf.WriteString(e.context)
	}

	return buf.String()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	v, err := json.Marshal(struct {
		Type string `json:"type"`
		Errors string `json:"errors"`
	}{Type:"validation", Errors:e.message})

	if err != nil {
		return nil, err
	}

	return v, nil
}*/

type NotFound struct {
	Err     error
}

// Error outputs stack info that should not be shown to client.
func (e *NotFound) Error() string {
	return e.Err.Error()
}

func (e *NotFound) Cause() error {
	return e.Err
}

type Conflict struct {
	Err     error
	Item string
}

func (e *Conflict) Error() string {
	return e.Err.Error()
}

func (e *Conflict) Cause() error {
	return e.Err
}

type Invalid struct {
	Err     error
}

func (e *Invalid) Error() string {
	return e.Err.Error()
}

func (e *Invalid) Cause() error {
	return e.Err
}

type Network struct {
	Err     error
}

func (e *Network) Error() string {
	return e.Err.Error()
}

func (e *Network) Cause() error {
	return e.Err
}

type Unauthorized struct {
	Err     error
}

func (e *Unauthorized) Error() string {
	return e.Err.Error()
}

func (e *Unauthorized) Cause() error {
	return e.Err
}
