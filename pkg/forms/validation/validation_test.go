package validation

import (
	"testing"
)

func TestValidation_Do(t *testing.T) {
	tests := []struct {
		name string
		key string
		rule StringRule
		value string
		errors string
	}{
		{
			name: "None empty string",
			key: "key",
			rule: StringRule{
				NotEmpty: true,
			},
			value: "Dictionary",
			errors: "",
		},
		{
			name: "Empty string",
			key: "key",
			rule: StringRule{
				NotEmpty: true,
			},
			value: "",
			errors: ErrEmptyStr.Error(),
		},
		{
			name: "More than min length",
			key: "key",
			rule: StringRule{
				MinLen: 5,
			},
			value: "Dictionary",
			errors: "",
		},
		{
			name: "Less than min length",
			rule: StringRule{
				MinLen: 5,
			},
			value: "Dict",
			errors: ErrLessThanRequired.Error(),
		},
		{
			name: "Less than max length",
			key: "key",
			rule: StringRule{
				MaxLen: 5,
			},
			value: "Dict",
			errors: "",
		},
		{
			name: "More than max length",
			rule: StringRule{
				MaxLen: 5,
			},
			value: "Dictionary",
			errors: ErrLargerThanRequired.Error(),
		},
		{
			name: "alphabets only",
			key: "key",
			rule: StringRule{
				ContentType: FieldAlphabets,
			},
			value: "Dictionary",
			errors: "",
		},
		{
			name: "alphabets only (invalid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldAlphabets,
			},
			value: "_Dictionary",
			errors: ErrNotAllLetters.Error(),
		},
		{
			name: "alphanumeric (valid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldAlphaNumeric,
			},
			value: "Dictionary1234",
			errors: "",
		},
		{
			name: "alphanumeric (invalid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldAlphaNumeric,
			},
			value: "Dictionary1234&",
			errors: ErrNotAlphaNumeric.Error(),
		},
		{
			name: "email (valid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldEmail,
			},
			value: "email@gmail.com",
			errors: "",
		},
		{
			name: "email (invalid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldEmail,
			},
			value: "emailgmail.com",
			errors: ErrNotValidEmail.Error(),
		},
		{
			name: "phone (valid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldPhone,
			},
			value: "08131231234",
			errors: "",
		},
		{
			name: "phone (invalid)",
			key: "key",
			rule: StringRule{
				ContentType: FieldPhone,
			},
			value: "339403d3",
			errors: ErrNotValidPhone.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Do(tt.key, tt.value, tt.rule)
			err := v.Errors()

			if len(tt.errors) == 0 && err != nil {
				t.Errorf("wanted no errors; got %v", err.ValidationErrors()[tt.key])
			} else if len(tt.errors) >= 1 && err == nil {
				t.Errorf("wanted %v; got nil",tt.errors)
			} else if err != nil { // if both wanted and obtained error are not empty
				ss := err.ValidationErrors()[tt.key]
				if ss != tt.errors {
					t.Errorf("wanted %v; got %v", tt.errors, err.ValidationErrors()[tt.key])
				}
			}
		})
	}
}


// todo:: keep, good functions
/*func foo(obtainedErr []string, wantedErrs []string) bool {
	return isSubset(obtainedErr, wantedErrs) && isSubset(wantedErrs, obtainedErr)
}

func isSubset(parent []string, sub []string) bool {
	for _, s := range sub {
		match := false
		for _, p := range parent {
			if s == p {
				match = true
			}
		}

		if !match {
			return false
		}
	}

	return true
}*/
