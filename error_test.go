package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackTrace(t *testing.T) {
	err := ParamError.New("param error", DbError.New("db error", errors.New("base error")))
	assert.Equal(t, err.Error(), "param error")
	assert.Equal(t, fmt.Sprintf("%s", err), "param error->db error->base error")
}

func TestWithField(t *testing.T) {
	err := ParamError.New("param error", DbError.New("db error", errors.New("base error"))).WithField("dbname", "user")
	assert.Equal(t, err.Error(), "param error")
	assert.Equal(t, fmt.Sprintf("%s", err), "param error->db error->base error")
	ms, me := json.Marshal(err.Fields())
	if me != nil {
		panic(me)
	}
	assert.Equal(t, string(ms), `{"dbname":"user"}`)
}

func TestErrCode_New(t *testing.T) {
	dbErr := DbError.Newf("there is an db error %s", "duplicate insert", errors.New("data duplicated"))
	assert.Equal(t, dbErr.code, DbError)
	assert.Equal(t, dbErr.code, dbErr.Code())
	assert.Equal(t, dbErr.statusCode, dbErr.StatusCode())
	assert.Equal(t, dbErr.message, dbErr.Message())
	assert.Equal(t, dbErr.message, "there is an db error duplicate insert")
	assert.Equal(t, dbErr.err, dbErr.Unwrap())
	assert.Equal(t, dbErr.Error(), "there is an db error duplicate insert")

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
