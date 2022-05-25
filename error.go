package errors

import (
	"fmt"
)

type BaseErrorInterface interface {
	StatusCode() HttpStatusCode
	Code() ErrCode
	Message() string
	Error() string
	Unwrap() error
}

type BaseError struct {
	code       ErrCode
	statusCode HttpStatusCode
	message    string
	err        error
}

func (e BaseError) Error() string {
	return fmt.Sprintf("error(%s), wrap(%v)", e.message, e.err)
}

func (e BaseError) Code() ErrCode {
	return e.code
}

func (e BaseError) StatusCode() HttpStatusCode {
	return e.statusCode
}

func (e BaseError) Message() string {
	return e.message
}

func (e BaseError) Unwrap() error {
	return e.err
}

type ErrCode int

func (e ErrCode) New(msg string, err error) *BaseError {
	return &BaseError{
		code:       e,
		statusCode: ErrorCodeWithHttpStatusCode[e].StatusCode().Get(),
		message:    e.defaultErrMessage(msg),
		err:        err,
	}
}

// NewW New with wrapped error ，只记录错误结构，使用默认错误提示
func (e ErrCode) Neww(err error) *BaseError {
	statusCode := ErrorCodeWithHttpStatusCode[e].StatusCode().Get()
	return &BaseError{
		code:       e,
		statusCode: statusCode,
		message:    ErrorCodeWithHttpStatusCode[e].Message(),
		err:        err,
	}
}

// Newf 接受多个参数，按照 format 参数拼接错误信息
func (e ErrCode) Newf(format string, args ...interface{}) *BaseError {
	var fmtArgs []interface{}
	var err error
	if len(args) > 0 {
		var ok bool
		err, ok = args[len(args)-1].(error)
		if !ok {
			err = nil
			fmtArgs = args
		} else {
			fmtArgs = args[:len(args)-1]
		}
	}
	msg := fmt.Sprintf(format, fmtArgs...)
	return &BaseError{
		code:       e,
		statusCode: ErrorCodeWithHttpStatusCode[e].StatusCode().Get(),
		message:    e.defaultErrMessage(msg),
		err:        err,
	}
}

func (e ErrCode) defaultErrMessage(msg string) string {
	if msg == "" {
		return ErrorCodeWithHttpStatusCode[e].Message()
	}
	return msg
}
