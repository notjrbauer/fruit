package fruit

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// Product errors.
const (
	ErrProductRequired   = Error("product required")
	ErrProductNotFound   = Error("product not found")
	ErrProductExists     = Error("product already exists")
	ErrProductIDRequired = Error("product id required")
)

// User errors.
const (
	ErrUserIDRequired = Error("user id required")
	ErrUserNotFound   = Error("user not found")
	ErrUserExists     = Error("user already exists")
	ErrUserRequired   = Error("user required")
)

// Error represents a fruit error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
