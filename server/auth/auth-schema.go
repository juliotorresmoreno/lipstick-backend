package auth

import (
	"errors"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type SignUpValidator struct {
	validator *validator.Validate
}

func NewSignUpValidator() *SignUpValidator {
	v := validator.New()
	return &SignUpValidator{validator: v}
}

func PhoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return len(phone) == 10 && strings.HasPrefix(phone, "5")
}

func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 7 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}

		if hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial {
			return true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

type SignUpValidationErrors struct {
	NameError     string `json:"name_error"`
	LastNameError string `json:"last_name_error"`
	PhoneError    string `json:"phone_error"`
	EmailError    string `json:"email_error"`
	PasswordError string `json:"password_error"`
}

func (cv *SignUpValidator) ValidateSignUp(form *SignUpPayload) (SignUpValidationErrors, error) {
	cv.validator.RegisterValidation("phone", PhoneValidation)
	cv.validator.RegisterValidation("password", PasswordValidation)

	err := cv.validator.Struct(form)
	if err != nil {
		errorsMap := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errorsMap[field] = "This field is required!"
			case "email":
				errorsMap[field] = "Invalid email format!"
			case "phone":
				errorsMap[field] = "Invalid phone number!"
			case "pattern":
				errorsMap[field] = "Password does not meet requirements!"
			default:
				errorsMap[field] = "Invalid field!"
			}
		}

		customErrors := SignUpValidationErrors{
			NameError:     errorsMap["Name"],
			LastNameError: errorsMap["LastName"],
			PhoneError:    errorsMap["Phone"],
			EmailError:    errorsMap["Email"],
			PasswordError: errorsMap["Password"],
		}

		return customErrors, errors.New("validation errors")
	}

	return SignUpValidationErrors{}, nil
}
