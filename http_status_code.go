package errors

import (
	"net/http"
)

type HttpStatusCode int

func (c HttpStatusCode) Get() HttpStatusCode {
	if c == 0 {
		return HttpStatusCode(http.StatusInternalServerError)
	}
	return c
}
