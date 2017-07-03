package caps

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

// Error represents a caps error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
