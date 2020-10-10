package domainerror

// ErrorCode is a custom int type for DomainError
type ErrorCode int

// DomainError is an error model for handling HTTP response layer.
// Each service needs to implement this interface for its error models
type DomainError interface {
	error
	Code() ErrorCode
}
