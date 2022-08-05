package errors

import (
	"fmt"
	"net/http"
)

var ErrorCodeWithHttpStatusCode = make(map[ErrCode]BaseError)

var (
	// base error
	ParamError      = NewErrorCode(1000, http.StatusBadRequest, "参数错误")
	UnAuthError     = NewErrorCode(1001, http.StatusUnauthorized, "资源未授权")
	ParseError      = NewErrorCode(1002, http.StatusBadRequest, "数据解析失败")
	PermissionError = NewErrorCode(1003, http.StatusForbidden, "资源未授权")
	NotFoundError   = NewErrorCode(1004, http.StatusNotFound, "资源不存在")
	ConflictError   = NewErrorCode(1005, http.StatusConflict, "资源产生冲突")
	InternalError   = NewErrorCode(1006, http.StatusInternalServerError, "系统内部错误")

	// data error
	DbError        = NewErrorCode(1100, http.StatusInternalServerError, "数据库发生意外")
	DuplicateError = NewErrorCode(1101, http.StatusBadRequest, "数据主键重复")

	// network error
	RpcError = NewErrorCode(1200, http.StatusInternalServerError, "远程调用发生意外")
)

// NewErrorCode 创建错误号，保证错误号唯一
func NewErrorCode(code int, httpStatusCode int, message string) ErrCode {
	if code == 0 {
		panic(fmt.Sprint("error code should not empty"))
	}
	errCode := ErrCode(code)
	exist := ErrorCodeWithHttpStatusCode[errCode]
	if exist.code != 0 {
		panic(fmt.Sprintf("error code is defined「%d: %s」，new 「%d: %s」", exist.code, exist.message, code, message))
	}
	ErrorCodeWithHttpStatusCode[errCode] = BaseError{
		code:       errCode,
		statusCode: HttpStatusCode(httpStatusCode),
		message:    message,
	}
	return errCode
}
