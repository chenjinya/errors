package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCode_New(t *testing.T) {
	dbErr := DbError.Newf("there is an db error %s", "duplicate insert", errors.New("data duplicated"))
	assert.Equal(t, dbErr.code, DbError)
	assert.Equal(t, dbErr.code, dbErr.Code())
	assert.Equal(t, dbErr.statusCode, dbErr.StatusCode())
	assert.Equal(t, dbErr.message, dbErr.Message())
	assert.Equal(t, dbErr.message, "there is an db error duplicate insert")
	assert.Equal(t, dbErr.err, dbErr.Unwrap())
	assert.Equal(t, dbErr.Error(), "error(there is an db error duplicate insert), wrap(data duplicated)")

	assert.Panics(t,
		func() {
			_ = NewErrorCode(1000, http.StatusBadRequest, "wrong error code")
		})
	assert.Panics(t,
		func() {
			_ = NewErrorCode(0, 0, "")
		})

	testDefaultStatusCodeError := ErrCode(1234).New("default status code test", nil)
	assert.Equal(t, http.StatusInternalServerError, int(testDefaultStatusCodeError.StatusCode().Get()))

	testEmptyMessageError := ParamError.Neww(nil)
	assert.Equal(t, http.StatusBadRequest, int(testEmptyMessageError.StatusCode()))
	assert.Equal(t, "参数错误", testEmptyMessageError.Message())
}
