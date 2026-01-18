package domain

import "fmt"

type ErrorCode string

const (
	ErrNotFound   ErrorCode = "NOT_FOUND"
	ErrValidation ErrorCode = "VALIDATION_ERROR"
	ErrConflict   ErrorCode = "STATUS_CONFLICT"
	ErrInternal   ErrorCode = "INTERNAL_ERROR"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewNotFound(msg string, details map[string]any) *DomainError {
	return &DomainError{Code: ErrNotFound, Message: msg, Details: details}
}

func NewValidation(msg string, details map[string]any) *DomainError {
	return &DomainError{Code: ErrValidation, Message: msg, Details: details}
}

func NewInternal(msg string, details map[string]any) *DomainError {
	return &DomainError{Code: ErrInternal, Message: msg, Details: details}
}
