package errors

import (
	"errors"
	"fmt"
)

var (
	New = errors.New
)

var (
	ErrBadRequest              = New400Response("ErrBadRequest")
	ErrInvalidParent           = New400Response("ErrInvalidParent")
	ErrNotAllowDeleteWithChild = New400Response("ErrNotAllowDeleteWithChild")
	ErrNotAllowDelete          = New400Response("ErrNotAllowDelete")
	ErrNotFoundUser            = New400Response("ErrNotFoundUser")
	ErrInvalidUsername         = New400Response("ErrInvalidUsername")
	ErrInvalidPassword         = New400Response("ErrInvalidPassword")
	ErrInvalidUser             = New400Response("ErrInvalidUser")
	ErrUserDisable             = New400Response("ErrUserDisable")
	ErrInvalidEmail            = New400Response("ErrInvalidEmail")
	ErrInvalidPhoneNumber      = New400Response("ErrInvalidPhone")

	ErrNoPerm          = NewResponse(401, 401, "ErrNoPerm")
	ErrInvalidToken    = NewResponse(9999, 401, "ErrInvalidToken")
	ErrNotFound        = NewResponse(404, 404, "ErrNotFound")
	ErrMethodNotAllow  = NewResponse(405, 405, "ErrMethodNotAllow")
	ErrTooManyRequests = NewResponse(429, 429, "ErrTooManyRequests")
	ErrInternalServer  = NewResponse(500, 500, "ErrInternalServer")
)

type ResponseError struct {
	Code       int
	Message    string
	StatusCode int
	ERR        error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

func WrapResponse(err error, code, statusCode int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(msg, args...),
		ERR:        err,
		StatusCode: statusCode,
	}
	return res
}

func Wrap400Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 400, 400, msg, args...)
}

func Wrap500Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 500, 500, msg, args...)
}

func NewResponse(code, statusCode int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(msg, args...),
		StatusCode: statusCode,
	}
	return res
}

func New400Response(msg string, args ...interface{}) error {
	return NewResponse(400, 400, msg, args...)
}

func New500Response(msg string, args ...interface{}) error {
	return NewResponse(500, 500, msg, args...)
}
