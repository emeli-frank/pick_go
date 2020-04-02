package validation

var (
	AddressRule = StringRule{
		NotEmpty:        true,
		MinLen:          6,
		MaxLen:          128,
	}
	EmailRule = StringRule{
		NotEmpty:        true,
		MinLen:          8,
		MaxLen:          128,
		ContentType:     FieldEmail,
	}
	PhoneRule = StringRule{
		NotEmpty:        true,
		//MinLen:          2,
		//MaxLen:          32,
		ContentType:     FieldPhone,
	}
	PasswordRule = StringRule{
		NotEmpty:        true,
		MinLen:          8,
		MaxLen:          32,
	}
)