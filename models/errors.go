package models

import "strings"

const (
	// ErrNotFound returned when resource cannot be found in the database
	ErrNotFound modelError = "models: resource not found"

	// ErrInvalidID returned when an invalid ID is provided to method like Delete
	ErrIDInvalid modelError = "models: ID must be > 0"

	// ErrPasswordIncorrect returned when incorrect password used for authentication
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired returned when is email is blank on form submission
	ErrEmailRequired modelError = "models: email address is required"

	// ErrEmailInvalid returned when email syntax is bad
	ErrEmailInvalid modelError = "models: email address is not valid"

	// ErrEmailTaken returned when an update/create is attempted with already registered address
	ErrEmailTaken modelError = "models: email address is already registered"

	// ErrPasswordTooShort return when password is too short in Create/Update
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	// ErrPasswordRequired returned when Create attempted without submitting password
	ErrPasswordRequired modelError = "models: password is required"

	// ErrRememberTooShort returned when remember token is too short
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"

	// ErrRememberRequired returned when Create/Update attempted without a
	// user remember token hash
	ErrRememberRequired modelError = "models: remember token is required"
)

// want the above variables to really be constants, but cant make type error a constant
type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
