package domainerror

import (
	"fmt"
)

const (
	JSONParseErrorCode ErrorCode = iota + 1
	NoneRequiredItemErrorCode
	ValidationErrorCode
)

// JSONParseErrorCode

type JSONParseError struct {
	requestUrl string
	err        error
}

func NewJSONParseError(requestUrl string, err error) *JSONParseError {
	return &JSONParseError{
		requestUrl: requestUrl,
		err:        err,
	}
}

func (e *JSONParseError) Error() string {
	return fmt.Sprintf("Request '%s' failed. Error message: %v",
		e.requestUrl, e.err)
}

func (e *JSONParseError) Code() ErrorCode {
	return JSONParseErrorCode
}

// NoneRequiredItemErrorCode

type NoneRequiredItemError struct {
	RequiredItemName string
}

func NewNoneRequiredItemError(itemName string) *NoneRequiredItemError {
	return &NoneRequiredItemError{
		RequiredItemName: itemName,
	}
}

func (e *NoneRequiredItemError) Error() string {
	return fmt.Sprintf(`'%s' is required`,
		e.RequiredItemName)
}

func (e *NoneRequiredItemError) Code() ErrorCode {
	return NoneRequiredItemErrorCode
}

// ValidationErrorCode

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{
		err: err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %+v", e.err)
}

func (e *ValidationError) Code() ErrorCode {
	return ValidationErrorCode
}
