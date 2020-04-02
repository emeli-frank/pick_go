package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalid = errors.New("validation error")

type validation struct {
	e Errors
}

// New creates a new validation object with the passed validation.Errors object
// or creates a new validation.Errors if nil was passed
func New(errs Errors) *validation {
	if errs == nil {
		errs = Errors{}
	}
	return &validation{errs}
}

// AddError adds outside error to the validation object while skipping empty error strings
func (v *validation) AddError(key string, message string, op string, err error) {
	v.e.Add(key, message, op, err)
}

// SetError receives a new validation.Errors object and replaces the existing error object
func (v *validation) SetError(err Errors) {
	v.e = err
}

// Errors returns validation errors from validation object or nil if there is non
func (v *validation) Errors() error {
	if len(v.e) == 0 {
		return nil
	}

	return v.e
}

func NotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return ErrEmptyStr
	}

	return nil
}

func HasCharacter(value string) bool {
	if len(value) < 1 {
		return false
	}

	return true
}

func CheckIntGreaterThanZero(value int) error {
	if value <= 0 {
		return errors.New("not positive integer")
	}

	return nil
}

func Length(value string, min int, max int) error {
	if utf8.RuneCountInString(value) < min {
		return errors.New(fmt.Sprintf("less than minimum (%d)", min))
	}
	if utf8.RuneCountInString(value) > max {
		return errors.New(fmt.Sprintf("more than maximum (%d)", max))
	}
	return nil
}

func IsAllLetters(value string) error {
	if err := NotEmpty(value); err != nil {
		//return errors.New("is empty")
		return err
	}

	for _, r := range value {
		if !unicode.IsLetter(r) {
			//return errors.New("not letters")
			return ErrNotAllLetters
		}
	}
	return nil
}

func IsLettersWithSpaces(value string) error {
	if err := NotEmpty(value); err != nil {
		//return errors.New("is empty")
		return err
	}

	for _, r := range value {
		fmt.Println(value)
		if !unicode.IsLetter(r) || string(r) == " " {
			//return errors.New("not letters")
			return ErrNotAlphaNumericWithSpaces
		}
	}
	return nil
}

func IsNumeric(value string) error {
	if err := NotEmpty(value); err != nil {
		//return errors.New("is empty")
		return ErrEmptyStr
	}

	// todo:: fix, some values (with commas) can jump this
	if _, err := strconv.Atoi(value); err != nil {
		//return errors.New("not numeric")
		return ErrNotNumeric
	}
	return nil
}

func IsAlphaNumeric(email string) error {
	if err := NotEmpty(email); err != nil {
		//return errors.New("is empty")
		return ErrEmptyStr
	}

	pattern := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if !pattern.MatchString(email) {
		//return errors.New("not email")
		return ErrNotAlphaNumeric
	}
	return nil
}

func IsEmail(email string) error {
	if err := NotEmpty(email); err != nil {
		//return errors.New("is empty")
		return ErrEmptyStr
	}

	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !pattern.MatchString(email) {
		//return errors.New("not email")
		return ErrNotValidEmail
	}
	return nil
}

func IsPhone(phone string) error {
	if err := NotEmpty(phone); err != nil {
		//return errors.New("is empty")
		return ErrEmptyStr
	}

	pattern := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !pattern.MatchString(phone) {
		//return errors.New("not email")
		return ErrNotValidPhone
	}
	return nil
}

func PermittedValues(value string, opts ...string) error {
	if err := NotEmpty(value); err != nil {
		return errors.New("is empty")
	}

	for _, v := range opts {
		if value != v {
			return errors.New(fmt.Sprintf("value %s is not permitted", v))
		}
	}
	return nil
}

/*type Rule struct {
	Str StringRule
	Num NumberRule
}*/

type StringRule struct {
	NotEmpty        bool
	MinLen          int
	MaxLen          int
	ContentType     FieldContentType
}

type NumberRule struct {
	PermittedValues []interface{}
	//GreaterOrEqualToOne bool
	//GreaterOrEqualToZero bool
}

type FieldContentType int
const (
	 FieldAlphabets FieldContentType 	= iota + 1
	 FieldNumbers
	 FieldAlphaNumeric
	 FieldEmail
	 FieldPhone
	 FieldAlphabetsWithSpaces
	 FieldName
	 FieldNames
)


func (v *validation) Do(key string, op string, value interface{}, rule interface{}) {
	switch rule := rule.(type) {
	case StringRule:
		switch val := value.(type) {
		case string:
			if rule.NotEmpty {
				if err := NotEmpty(val); err != nil {
					v.AddError(key, "field cannot be empty", op, err)
				}
			}
			if rule.MinLen > 1 && (rule.NotEmpty || HasCharacter(val)) {
				if err := MinLen(val, rule.MinLen); err != nil {
					msg := fmt.Sprintf(
						"field is less than minimum allowed number of characters, should be more than %d",
						rule.MinLen)
					v.AddError(key, msg, op, err)
				}
			}
			if rule.MaxLen > 1 && (rule.NotEmpty || HasCharacter(val)) {
				if err := MaxLen(val, rule.MaxLen); err != nil {
					msg := fmt.Sprintf(
						"field is more than maximum allowed number of characters, should be less than %d",
						rule.MaxLen)
					v.AddError(key, msg, op, err)
				}
			}
			if rule.ContentType >= 1 && (rule.NotEmpty || HasCharacter(val)) {
				if rule.ContentType == FieldAlphabets {
					if err := IsAllLetters(val); err != nil {
						v.AddError(key, "field can contain only letters", op, err)
					}
				} else if rule.ContentType == FieldNumbers {
					if err := IsNumeric(val); err != nil {
						v.AddError(key, "string can contain only numbers", op, err)
					}
				} else if rule.ContentType == FieldAlphaNumeric {
					if err := IsAlphaNumeric(val); err != nil {
						v.AddError(key, "field can contain only alphabets and numbers", op, err)
					}
				} else if rule.ContentType == FieldEmail {
					if err := IsEmail(val); err != nil {
						v.AddError(key, "field is not a valid email", op, err)
					}
				} else if rule.ContentType == FieldPhone {
					if err := IsPhone(val); err != nil {
						v.AddError(key, "field is not a valid phone number", op, err)
					}
				} else if rule.ContentType == FieldAlphabetsWithSpaces {
					if err := IsLettersWithSpaces(val); err != nil {
						v.AddError(key, "field can contain only alphabets (words can be separated with a space)", op, err)
					}
				} else if rule.ContentType == FieldName {
					// todo:: allow other characters that appear in names
					if err := IsAllLetters(val); err != nil {
						v.AddError(key, "field can contain only alphabets", op, err)
					}
				}
			}
		}
	}

}




func MinLen(value string, min int) error {
	if utf8.RuneCountInString(value) < min {
		//return errors.New(fmt.Sprintf("less than minimum (%d)", min))
		return ErrLessThanRequired
	}
	return nil
}

func MaxLen(value string, max int) error {
	if utf8.RuneCountInString(value) > max {
		fmt.Println("max: ", max)
		fmt.Println("val: ", value)
		//return errors.New(fmt.Sprintf("more than maximum (%d)", max))
		return ErrLargerThanRequired
	}
	return nil
}

func (v *validation) CustomValidate(op string, f func() (k string, msg string, err error)) {
	k, msg, err := f()
	if err == nil {
		return
	}
	v.AddError(k, msg, op, err)
}

/*func IsAlphaNumeric(value string) error {
	if err := NotEmpty(value); err != nil {
		return errors.New("is empty")
	}

	for _, r := range value {
		if !unicode.IsLetter(r) {
			if err := IsNumeric(value); err != nil {
				return errors.New("not letters")
			}
		}
	}
	return nil
}*/

var (
	ErrEmptyStr           = errors.New("empty string")
	ErrLargerThanRequired = errors.New("larger than maximum allowed number of characters")
	ErrLessThanRequired = errors.New("less than minimum allowed number of characters")
	ErrNotAllLetters = errors.New("all characters are not letters")
	ErrNotValidEmail = errors.New("not valid email")
	ErrNotValidPhone = errors.New("not valid phone")
	ErrNotAlphaNumeric = errors.New("not alphanumeric")
	ErrNotNumeric = errors.New("not numeric")
	ErrNotAlphaNumericWithSpaces = errors.New("not alphanumeric with spaces")
)
