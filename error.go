package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
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
	stack      []uintptr
}

func (e BaseError) MarshalJSON() (b []byte, err error) {
	return json.Marshal(map[string]interface{}{
		"code":        e.Code(),
		"status_code": e.StatusCode(),
		"message":     e.Message(),
		"error":       e.Error(),
	})
}

func (e BaseError) Error() string {
	return fmt.Sprintf("error(%s), wrap(%s)", e.message, e.err)
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

func (e *BaseError) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(st, "%s", e.Message())
		for _, pc := range e.stack {
			f := Frame(pc)
			fmt.Fprintf(st, "\n%s", f.FileAndLine())
		}
	case 's':
		fmt.Fprintf(st, "%s", e.errors())
	}
}

func (e *BaseError) errors() string {
	es := []string{e.message}
	cause := e.err
	for {
		if cause == nil {
			break
		}
		be, ok := cause.(*BaseError)
		if !ok {
			es = append(es, cause.Error())
			break
		}
		es = append(es, be.message)
		cause = be.err
	}
	return strings.Join(es, "->")
}

type ErrCode int

func (e ErrCode) New(msg string, err error) *BaseError {
	return &BaseError{
		code:       e,
		statusCode: ErrorCodeWithHttpStatusCode[e].StatusCode().Get(),
		message:    e.defaultErrMessage(msg),
		err:        err,
		stack:      e.callers(),
	}
}

// NewW New with wrapped error ，只记录错误结构，使用默认错误提示
func (e ErrCode) Neww(err error) *BaseError {
	statusCode := ErrorCodeWithHttpStatusCode[e].StatusCode().Get()
	return &BaseError{
		code:       e,
		statusCode: statusCode,
		message:    e.defaultErrMessage(""),
		err:        err,
		stack:      e.callers(),
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
		stack:      e.callers(),
	}
}

func (e ErrCode) defaultErrMessage(msg string) string {
	if msg == "" {
		return ErrorCodeWithHttpStatusCode[e].Message()
	}
	return msg
}

func (e ErrCode) callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	stack := pcs[0:n]

	return stack
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")

	return name[i+1:]
}

type Frame uintptr

func (f Frame) pc() uintptr {
	return uintptr(f)
}

func (f Frame) FileAndLine() string {
	fn := runtime.FuncForPC(f.pc())
	file, line := fn.FileLine(f.pc())

	return fmt.Sprintf("%s:%d:%s", file, line, funcname(fn.Name()))
}
