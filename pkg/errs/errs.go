package errs

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")

	ErrOrderNotFound = errors.New("order not found")
	ErrAccessDenied  = errors.New("access denied")

	ErrInvalidOrderPrice  = errors.New("order price must be greater than zero")
	ErrEmptyProductStatus = errors.New("product name cannot be empty")
	ErrCannotCancelOrder  = errors.New("cannot cancel an order that is already completed")

	ErrEmptyName   = errors.New("name cannot be empty")
	ErrInvalidRole = errors.New("invalid role provided")
	ErrInvalidStatus = errors.New("invalid order status")
)
